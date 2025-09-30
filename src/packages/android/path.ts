import { stat } from 'fs/promises';
import { join } from 'path';
import { getPackageName } from './package-name.js';

/**
 * Returns the path to app/build.gradle or app/build.gradle.kts
 */
export async function getAppBuildGradlePath(projectRoot: string): Promise<string> {
  const gradlePaths = [
    join(projectRoot, 'app', 'build.gradle'),
    join(projectRoot, 'app', 'build.gradle.kts'),
  ];

  for (const path of gradlePaths) {
    try {
      await stat(path);
      return path;
    } catch {
      // File doesn't exist, continue
    }
  }

  return '';
}

/**
 * Returns the base directory path (app/src/main/java or app/src/main/kotlin)
 */
export async function getBaseDirPath(projectRoot: string): Promise<string> {
  const javaDir = join(projectRoot, 'app', 'src', 'main', 'java');
  const kotlinDir = join(projectRoot, 'app', 'src', 'main', 'kotlin');

  try {
    await stat(javaDir);
    return javaDir;
  } catch {
    // Java dir doesn't exist, try Kotlin
  }

  try {
    await stat(kotlinDir);
    return kotlinDir;
  } catch {
    // Neither exists
  }

  return '';
}

/**
 * Returns the source directory path (e.g., app/src/main/java/com/example/app)
 */
export async function getSourceDirPath(projectRoot: string): Promise<string> {
  const baseDir = await getBaseDirPath(projectRoot);
  if (!baseDir) {
    return '';
  }

  const packageName = await getPackageName(projectRoot);
  if (!packageName) {
    return '';
  }

  const packagePath = packageName.replace(/\./g, '/');
  const sourceDir = join(baseDir, packagePath);

  return sourceDir;
}

/**
 * Returns the path to app/src/main/AndroidManifest.xml
 */
export function getAndroidManifestPath(projectRoot: string): string {
  return join(projectRoot, 'app', 'src', 'main', 'AndroidManifest.xml');
}

/**
 * Returns the path to gradle/libs.versions.toml
 */
export async function getVersionCatalogPath(projectRoot: string): Promise<string> {
  const path = join(projectRoot, 'gradle', 'libs.versions.toml');
  try {
    await stat(path);
    return path;
  } catch {
    return '';
  }
}

/**
 * Checks if the project has a version catalog (gradle/libs.versions.toml)
 */
export async function hasVersionCatalog(projectRoot: string): Promise<boolean> {
  const path = await getVersionCatalogPath(projectRoot);
  return path !== '';
}
