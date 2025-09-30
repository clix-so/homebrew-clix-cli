import {
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
import { hasVersionCatalog } from './path.js';

/**
 * Runs all Android doctor checks.
 * This function performs read-only checks and reports the status.
 * Note: In the TypeScript version, we don't have access to the logx module,
 * so the caller should handle logging based on the return values.
 */
export async function runAndroidDoctor(projectRoot: string): Promise<{
  gradleRepository: boolean;
  hasVersionCatalog: boolean;
  clixDependency: boolean;
  gmsPlugin: boolean;
  clixInitialization: { success: boolean; code: string };
  permissions: boolean;
  googleServicesJson: boolean;
}> {
  const gradleRepository = await checkGradleRepository(projectRoot);
  const versionCatalog = await hasVersionCatalog(projectRoot);
  const clixDependency = await checkGradleDependency(projectRoot);
  const gmsPlugin = await checkGradlePlugin(projectRoot);
  const [initSuccess, initCode] = await checkClixCoreImport(projectRoot);
  const permissions = await checkAndroidMainActivityPermissions(projectRoot);
  const googleServicesJson = await checkGoogleServicesJSON(projectRoot);

  return {
    gradleRepository,
    hasVersionCatalog: versionCatalog,
    clixDependency,
    gmsPlugin,
    clixInitialization: { success: initSuccess, code: initCode },
    permissions,
    googleServicesJson,
  };
}

// Export helper functions for use in other modules
export { contains, indexOf, stringContainsClixInitializeInOnCreate };
