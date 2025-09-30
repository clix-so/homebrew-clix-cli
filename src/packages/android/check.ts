import { readFile, stat } from 'fs/promises';
import { join, sep } from 'path';
import { readdir } from 'fs/promises';
import { getAppBuildGradlePath, getVersionCatalogPath, hasVersionCatalog, getSourceDirPath, getAndroidManifestPath, getBaseDirPath } from './path.js';
import { extractApplicationClassName } from './manifest-parser.js';

/**
 * Helper function to check if a string contains a substring
 */
export function contains(str: string, substr: string): boolean {
  return str.includes(substr);
}

/**
 * Helper function to find the index of a substring
 */
export function indexOf(str: string, substr: string): number {
  return str.indexOf(substr);
}

/**
 * Removes all whitespace from a string for pattern matching
 */
function removeAllWhitespace(s: string): string {
  return s.replace(/[\s\t\n\r]/g, '');
}

/**
 * Checks if Clix.initialize(this, ...) is called inside onCreate
 */
export function stringContainsClixInitializeInOnCreate(content: string): boolean {
  // Find signature variants
  let onCreateIdx = -1;
  const voidIdx = content.indexOf('void onCreate'); // Java
  const funIdx = content.indexOf('fun onCreate'); // Kotlin

  if (voidIdx !== -1) {
    onCreateIdx = voidIdx;
  } else if (funIdx !== -1) {
    onCreateIdx = funIdx;
  }

  if (onCreateIdx === -1) {
    return false;
  }

  // Find opening brace after signature
  let openIdx = -1;
  for (let i = onCreateIdx; i < content.length; i++) {
    const c = content[i];
    if (c === '{') {
      openIdx = i;
      break;
    }
    if (c === ';') {
      // Not a method definition
      return false;
    }
  }

  if (openIdx === -1) {
    return false;
  }

  // Extract full block by brace depth tracking
  let depth = 0;
  let endIdx = -1;
  for (let i = openIdx; i < content.length; i++) {
    const ch = content[i];
    if (ch === '{') {
      depth++;
    } else if (ch === '}') {
      depth--;
      if (depth === 0) {
        endIdx = i + 1;
        break;
      }
    }
  }

  if (endIdx === -1) {
    // Malformed braces fallback (limit 500 chars)
    endIdx = Math.min(openIdx + 500, content.length);
  }

  const block = content.substring(openIdx, endIdx);
  const norm = removeAllWhitespace(block);

  // Basic pattern (no whitespace) e.g., Clix.initialize(this
  if (contains(norm, 'Clix.initialize(this')) {
    return true;
  }

  // Allow generics: Clix.initialize<...>(this
  if (contains(norm, 'Clix.initialize<') && contains(norm, '(this')) {
    return true;
  }

  return false;
}

/**
 * Checks if mavenCentral() is present in settings.gradle(.kts) or build.gradle(.kts)
 */
export async function checkGradleRepository(projectRoot: string): Promise<boolean> {
  const gradleFiles = [
    join(projectRoot, 'settings.gradle'),
    join(projectRoot, 'settings.gradle.kts'),
    join(projectRoot, 'build.gradle'),
    join(projectRoot, 'build.gradle.kts'),
  ];

  for (const file of gradleFiles) {
    try {
      const data = await readFile(file, 'utf-8');
      const content = data.toString();
      if (contains(content, 'repositories') && contains(content, 'mavenCentral()')) {
        return true;
      }
    } catch {
      // File doesn't exist, continue
    }
  }

  return false;
}

/**
 * Checks if so.clix:clix-android-sdk is present in app/build.gradle(.kts)
 */
export async function checkGradleDependency(projectRoot: string): Promise<boolean> {
  const appBuildGradleFilePath = await getAppBuildGradlePath(projectRoot);

  if (!appBuildGradleFilePath) {
    return false;
  }

  try {
    const data = await readFile(appBuildGradleFilePath, 'utf-8');
    const content = data.toString();

    if (contains(content, 'implementation("so.clix:clix-android-sdk:')) {
      return true;
    } else if (contains(content, 'implementation(libs.clix.android.sdk)')) {
      // If using version catalog, also confirm alias exists in libs.versions.toml when possible
      if (await hasVersionCatalog(projectRoot)) {
        const catalog = await getVersionCatalogPath(projectRoot);
        try {
          const catData = await readFile(catalog, 'utf-8');
          const cat = catData.toString();
          if (contains(cat, '[libraries]') && contains(cat, 'clix-android-sdk = ')) {
            return true;
          }
        } catch {
          // Catalog read failed
        }
      } else {
        return true;
      }
    }

    return false;
  } catch {
    return false;
  }
}

/**
 * Checks if com.google.gms:google-services is present in app/build.gradle(.kts)
 */
export async function checkGradlePlugin(projectRoot: string): Promise<boolean> {
  const gradleFiles = [
    join(projectRoot, 'app', 'build.gradle'),
    join(projectRoot, 'app', 'build.gradle.kts'),
  ];

  for (const file of gradleFiles) {
    try {
      const data = await readFile(file, 'utf-8');
      const content = data.toString();

      if (contains(content, 'alias(libs.plugins.gms')) {
        return true;
      }
      if (contains(content, 'id("com.google.gms.google-services")')) {
        return true;
      }
    } catch {
      // File doesn't exist, continue
    }
  }

  return false;
}

/**
 * Checks if any Application class imports so.clix.core.Clix (Java or Kotlin).
 * Returns [success, errorCode]
 */
export async function checkClixCoreImport(projectRoot: string): Promise<[boolean, string]> {
  const manifestPath = join(projectRoot, 'app', 'src', 'main', 'AndroidManifest.xml');
  let appName: string;

  try {
    appName = await extractApplicationClassName(manifestPath);
  } catch {
    return [false, 'unknown'];
  }

  if (!appName) {
    return [false, 'missing-application'];
  }

  let appPath = appName.startsWith('.') ? appName.substring(1) : appName;
  appPath = appPath.replace(/\./g, sep);

  const sourceDir = await getSourceDirPath(projectRoot);
  if (!sourceDir) {
    return [false, 'unknown'];
  }

  const ktPath = join(sourceDir, appPath + '.kt');
  const javaPath = join(sourceDir, appPath + '.java');

  let actualPath = '';
  try {
    await stat(javaPath);
    actualPath = javaPath;
  } catch {
    try {
      await stat(ktPath);
      actualPath = ktPath;
    } catch {
      return [false, 'unknown'];
    }
  }

  try {
    const data = await readFile(actualPath, 'utf-8');
    const content = data.toString();

    if (stringContainsClixInitializeInOnCreate(content)) {
      return [true, ''];
    } else {
      return [false, 'missing-content'];
    }
  } catch {
    return [false, 'unknown'];
  }
}

/**
 * Checks MainActivity for permission request code
 */
export async function checkAndroidMainActivityPermissions(projectRoot: string): Promise<boolean> {
  const javaDir = join(projectRoot, 'app', 'src', 'main', 'java');
  const kotlinDir = join(projectRoot, 'app', 'src', 'main', 'kotlin');
  const mainActivityFiles: string[] = [];

  const findMainActivity = async (root: string) => {
    try {
      const walk = async (dir: string) => {
        const files = await readdir(dir, { withFileTypes: true });
        for (const file of files) {
          const fullPath = join(dir, file.name);
          if (file.isDirectory()) {
            await walk(fullPath);
          } else if (file.name === 'MainActivity.java' || file.name === 'MainActivity.kt') {
            mainActivityFiles.push(fullPath);
          }
        }
      };
      await walk(root);
    } catch {
      // Directory doesn't exist
    }
  };

  await findMainActivity(javaDir);
  await findMainActivity(kotlinDir);

  if (mainActivityFiles.length === 0) {
    return false;
  }

  const permissionPattern = [
    'requestPermissions(',
    'ActivityCompat.requestPermissions(',
    'ContextCompat.checkSelfPermission(',
    'Manifest.permission.',
  ];

  for (const file of mainActivityFiles) {
    try {
      const data = await readFile(file, 'utf-8');
      const content = data.toString();

      for (const pat of permissionPattern) {
        if (contains(content, pat)) {
          return true;
        }
      }
    } catch {
      // File read failed, continue
    }
  }

  return false;
}

/**
 * Checks if google-services.json exists in the correct location
 */
export async function checkGoogleServicesJSON(projectRoot: string): Promise<boolean> {
  const gsPath = join(projectRoot, 'app', 'google-services.json');
  try {
    await stat(gsPath);
    return true;
  } catch {
    return false;
  }
}
