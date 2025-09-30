import { readFile, writeFile } from 'fs/promises';
import { join } from 'path';
import {
  checkGoogleServicesJSON,
  checkGradleRepository,
  checkGradleDependency,
  checkGradlePlugin,
  checkClixCoreImport,
  checkAndroidMainActivityPermissions,
  contains,
  indexOf,
} from './check.js';
import {
  getAppBuildGradlePath,
  getVersionCatalogPath,
  hasVersionCatalog,
  getSourceDirPath,
  getAndroidManifestPath,
} from './path.js';
import { getPackageName } from './package-name.js';
import { ANDROID_CLIX_SDK_VERSION, ANDROID_GMS_PLUGIN_VERSION } from '../versions.js';

/**
 * Guides the user through the Android installation checklist.
 * Returns an object with the status of each check and any auto-fix results.
 */
export async function handleAndroidInstall(
  apiKey: string,
  projectId: string
): Promise<{
  googleServicesJson: { ok: boolean; autoFixed: boolean };
  gradleRepository: { ok: boolean; autoFixed: boolean };
  clixDependency: { ok: boolean; autoFixed: boolean; usedVersionCatalog: boolean };
  gmsPlugin: { ok: boolean; autoFixed: boolean };
  clixInitialization: { ok: boolean; autoFixed: boolean; code: string };
  permissions: { ok: boolean };
}> {
  const projectRoot = process.cwd();

  // Check google-services.json
  const googleServicesJsonOk = await checkGoogleServicesJSON(projectRoot);

  // Check Gradle repository
  let gradleRepoOk = await checkGradleRepository(projectRoot);
  let gradleRepoAutoFixed = false;
  if (!gradleRepoOk) {
    gradleRepoAutoFixed = await addGradleRepository(projectRoot);
    if (gradleRepoAutoFixed) {
      gradleRepoOk = true;
    }
  }

  // Check Clix dependency
  let clixDepOk = await checkGradleDependency(projectRoot);
  let clixDepAutoFixed = false;
  let usedVersionCatalog = false;

  if (!clixDepOk) {
    // Prefer Version Catalog if available
    if (await hasVersionCatalog(projectRoot)) {
      if (
        (await ensureClixInVersionCatalog(projectRoot)) &&
        (await wireClixDependencyAlias(projectRoot))
      ) {
        clixDepAutoFixed = true;
        clixDepOk = true;
        usedVersionCatalog = true;
      }
    }

    // Fallback to direct dependency if still not ok
    if (!clixDepOk) {
      clixDepAutoFixed = await addGradleDependency(projectRoot);
      if (clixDepAutoFixed) {
        clixDepOk = true;
      }
    }
  }

  // Check GMS plugin
  let gmsPluginOk = await checkGradlePlugin(projectRoot);
  let gmsPluginAutoFixed = false;
  if (!gmsPluginOk) {
    gmsPluginAutoFixed = await addGradlePlugin(projectRoot);
    if (gmsPluginAutoFixed) {
      gmsPluginOk = true;
    }
  }

  // Check Clix initialization
  let [appOk, code] = await checkClixCoreImport(projectRoot);
  const originalOk = appOk;
  let autoCreated = false;

  if (!appOk && code === 'missing-application') {
    const [ok] = await addApplication(projectRoot, apiKey, projectId);
    if (ok) {
      appOk = true;
      autoCreated = true;
    }
  }

  // Check permissions
  const mainActivityOk = await checkAndroidMainActivityPermissions(projectRoot);

  return {
    googleServicesJson: { ok: googleServicesJsonOk, autoFixed: false },
    gradleRepository: { ok: gradleRepoOk, autoFixed: gradleRepoAutoFixed },
    clixDependency: { ok: clixDepOk, autoFixed: clixDepAutoFixed, usedVersionCatalog },
    gmsPlugin: { ok: gmsPluginOk, autoFixed: gmsPluginAutoFixed },
    clixInitialization: { ok: appOk, autoFixed: autoCreated && !originalOk, code },
    permissions: { ok: mainActivityOk },
  };
}

/**
 * Tries to insert mavenCentral() into settings.gradle(.kts) or build.gradle(.kts)
 */
async function addGradleRepository(projectRoot: string): Promise<boolean> {
  const gradleFiles = [
    join(projectRoot, 'settings.gradle'),
    join(projectRoot, 'settings.gradle.kts'),
    join(projectRoot, 'build.gradle'),
    join(projectRoot, 'build.gradle.kts'),
  ];

  for (const file of gradleFiles) {
    try {
      const data = await readFile(file, 'utf-8');
      let content = data.toString();

      if (contains(content, 'repositories') && contains(content, 'mavenCentral()')) {
        return true; // already present
      }

      // Try to insert after 'repositories {' or at end
      const idx = indexOf(content, 'repositories {');
      if (idx !== -1) {
        const insertAt = idx + 'repositories {'.length;
        const newContent = content.slice(0, insertAt) + '\n    mavenCentral()' + content.slice(insertAt);
        await writeFile(file, newContent, 'utf-8');
        return true;
      }
    } catch {
      // File doesn't exist, continue
    }
  }

  return false;
}

/**
 * Tries to insert the Clix SDK dependency into app/build.gradle(.kts)
 */
async function addGradleDependency(projectRoot: string): Promise<boolean> {
  const gradleFiles = [
    join(projectRoot, 'app', 'build.gradle'),
    join(projectRoot, 'app', 'build.gradle.kts'),
  ];

  for (const file of gradleFiles) {
    try {
      const data = await readFile(file, 'utf-8');
      let content = data.toString();

      if (contains(content, 'implementation("so.clix:clix-android-sdk')) {
        return true; // already present
      }

      // Try to insert after 'dependencies {' or at end
      const idx = indexOf(content, 'dependencies {');
      if (idx !== -1) {
        const insertAt = idx + 'dependencies {'.length;
        const newContent =
          content.slice(0, insertAt) +
          `\n    implementation("so.clix:clix-android-sdk:${ANDROID_CLIX_SDK_VERSION}")` +
          content.slice(insertAt);
        await writeFile(file, newContent, 'utf-8');
        return true;
      }
    } catch {
      // File doesn't exist, continue
    }
  }

  return false;
}

/**
 * Adds clix coordinates and alias into libs.versions.toml if missing.
 * Adds under [versions], [libraries], and optionally [bundles] if needed.
 */
async function ensureClixInVersionCatalog(projectRoot: string): Promise<boolean> {
  const catalog = await getVersionCatalogPath(projectRoot);
  if (!catalog) {
    return false;
  }

  try {
    const data = await readFile(catalog, 'utf-8');
    let content = data.toString();

    const desiredVersionKey = 'clix';
    const desiredVersion = ANDROID_CLIX_SDK_VERSION;
    const libAlias = 'clix-android-sdk';

    let changed = false;

    if (!contains(content, '[versions]')) {
      content += '\n[versions]\n';
      changed = true;
    }
    if (!contains(content, `${desiredVersionKey} = "`)) {
      content = content.replace(
        '[versions]',
        `[versions]\n${desiredVersionKey} = "${desiredVersion}"`
      );
      changed = true;
    }

    if (!contains(content, '[libraries]')) {
      content += '\n[libraries]\n';
      changed = true;
    }
    const libLine = `${libAlias} = { module = "so.clix:clix-android-sdk", version.ref = "${desiredVersionKey}" }`;
    if (!contains(content, `${libAlias} = `)) {
      content = content.replace('[libraries]', `[libraries]\n${libLine}`);
      changed = true;
    }

    if (!changed) {
      // Already present
      return true;
    }

    await writeFile(catalog, content, 'utf-8');
    return true;
  } catch {
    return false;
  }
}

/**
 * Ensures app/build.gradle(.kts) uses implementation(libs.clix.android.sdk)
 */
async function wireClixDependencyAlias(projectRoot: string): Promise<boolean> {
  const gradleFiles = [
    join(projectRoot, 'app', 'build.gradle'),
    join(projectRoot, 'app', 'build.gradle.kts'),
  ];

  for (const file of gradleFiles) {
    try {
      const data = await readFile(file, 'utf-8');
      let content = data.toString();

      if (
        contains(content, 'implementation(libs.clix.android.sdk)') ||
        contains(content, 'implementation("so.clix:clix-android-sdk:')
      ) {
        return true;
      }

      const idx = indexOf(content, 'dependencies {');
      if (idx !== -1) {
        const insertAt = idx + 'dependencies {'.length;
        const newContent =
          content.slice(0, insertAt) + '\n    implementation(libs.clix.android.sdk)' + content.slice(insertAt);
        await writeFile(file, newContent, 'utf-8');
        return true;
      }
    } catch {
      // File doesn't exist, continue
    }
  }

  return false;
}

/**
 * Tries to insert the Google services plugin into app/build.gradle(.kts)
 */
async function addGradlePlugin(projectRoot: string): Promise<boolean> {
  const gradleFiles = [
    join(projectRoot, 'app', 'build.gradle'),
    join(projectRoot, 'app', 'build.gradle.kts'),
  ];

  for (const file of gradleFiles) {
    try {
      const data = await readFile(file, 'utf-8');
      let content = data.toString();

      if (contains(content, 'id("com.google.gms.google-services")')) {
        return true; // already present
      }

      // Try to insert after 'plugins {' or at end
      const idx = indexOf(content, 'plugins {');
      if (idx !== -1) {
        const insertAt = idx + 'plugins {'.length;
        const newContent =
          content.slice(0, insertAt) +
          `\n    id("com.google.gms.google-services") version "${ANDROID_GMS_PLUGIN_VERSION}"` +
          content.slice(insertAt);
        await writeFile(file, newContent, 'utf-8');
        return true;
      }
    } catch {
      // File doesn't exist, continue
    }
  }

  return false;
}

/**
 * Creates the BasicApplication.kt file with Clix initialization
 */
async function addApplicationFile(
  projectRoot: string,
  apiKey: string,
  projectId: string
): Promise<[boolean, string]> {
  const sourceDir = await getSourceDirPath(projectRoot);
  if (!sourceDir) {
    return [false, 'Could not find source directory for Android project.'];
  }

  const appBuildGradlePath = await getAppBuildGradlePath(projectRoot);
  if (!appBuildGradlePath) {
    return [false, 'Could not find app/build.gradle(.kts) file.'];
  }

  const packageName = await getPackageName(projectRoot);
  if (!packageName) {
    return [false, 'Could not extract package name from app/build.gradle(.kts).'];
  }

  const filePath = join(sourceDir, 'BasicApplication.kt');
  const code = `package ${packageName}

import android.app.Application
import so.clix.core.Clix
import so.clix.core.ClixConfig

class BasicApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        Clix.initialize(this, ClixConfig(
\t\t\tprojectId = "${projectId}",
            apiKey = "${apiKey}",
        ))
    }
}`;

  try {
    await writeFile(filePath, code, 'utf-8');
    return [true, 'Application class created successfully'];
  } catch (err: any) {
    return [false, err.message || 'Failed to create Application file'];
  }
}

/**
 * Adds the application name to AndroidManifest.xml
 */
async function addApplicationNameToManifest(projectRoot: string): Promise<[boolean, string]> {
  const manifestPath = getAndroidManifestPath(projectRoot);

  try {
    const data = await readFile(manifestPath, 'utf-8');
    let content = data.toString();

    const newContent = content.replace('<application', '<application android:name=".BasicApplication"');
    await writeFile(manifestPath, newContent, 'utf-8');
    return [true, 'Application name added to AndroidManifest.xml'];
  } catch {
    return [false, 'Failed to write AndroidManifest.xml'];
  }
}

/**
 * Creates the Application class and adds it to the manifest
 */
async function addApplication(
  projectRoot: string,
  apiKey: string,
  projectId: string
): Promise<[boolean, string]> {
  // Step 1: Create BasicApplication.kt file
  const [ok1, message1] = await addApplicationFile(projectRoot, apiKey, projectId);
  if (!ok1) {
    return [false, message1];
  }

  // Step 2: Add application name to AndroidManifest.xml
  const [ok2, message2] = await addApplicationNameToManifest(projectRoot);
  if (!ok2) {
    return [false, message2];
  }

  return [true, 'Application class setup complete'];
}
