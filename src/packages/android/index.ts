// Main functions
export { handleAndroidInstall } from './installer.js';
export { runAndroidDoctor } from './doctor.js';
export { uninstallClixAndroid } from './uninstaller.js';

// Check functions
export {
  checkGradleRepository,
  checkGradleDependency,
  checkGradlePlugin,
  checkClixCoreImport,
  checkAndroidMainActivityPermissions,
  checkGoogleServicesJSON,
  contains,
  indexOf,
  stringContainsClixInitializeInOnCreate,
} from './check.js';

// Path utilities
export {
  getAppBuildGradlePath,
  getBaseDirPath,
  getSourceDirPath,
  getAndroidManifestPath,
  getVersionCatalogPath,
  hasVersionCatalog,
} from './path.js';

// Package name utilities
export { getPackageName } from './package-name.js';

// Manifest parser
export { extractApplicationClassName } from './manifest-parser.js';
