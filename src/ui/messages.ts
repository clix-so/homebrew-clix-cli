import { ANDROID_CLIX_SDK_VERSION, ANDROID_GMS_PLUGIN_VERSION } from '../packages/versions.js';

export const MESSAGES = {
  // Gradle Repository
  TITLE_GRADLE_REPO_CHECK: 'Checking Gradle repository configuration...',
  MSG_GRADLE_REPO_SUCCESS: 'Gradle repositories are properly configured.',
  MSG_GRADLE_REPO_FAILURE: 'Gradle repository configuration is missing.',
  FIX_GRADLE_REPO: 'To resolve this, add the following to settings.gradle(.kts) or build.gradle(.kts):',
  MSG_GRADLE_REPO_FIX_FAILURE: 'Automatic fix failed. Please add the following to settings.gradle(.kts) or build.gradle(.kts):',
  CODE_GRADLE_REPO: `repositories {
\tmavenCentral()
}`,

  // Clix SDK Dependency
  TITLE_CLIX_DEPENDENCY_CHECK: 'Checking Clix SDK dependency...',
  MSG_CLIX_DEPENDENCY_SUCCESS: 'Clix SDK dependency detected.',
  MSG_CLIX_DEPENDENCY_FAILURE: 'Clix SDK dependency is missing.',
  FIX_CLIX_DEPENDENCY: 'To resolve this, add the following to app/build.gradle(.kts) or use Version Catalog entries:',
  MSG_CLIX_DEPENDENCY_FIX_FAILURE: 'Automatic fix failed. Please add the following to app/build.gradle(.kts):',
  CODE_CLIX_DEPENDENCY: `dependencies {
\timplementation("so.clix:clix-android-sdk:${ANDROID_CLIX_SDK_VERSION}")
}`,

  // Version Catalog
  TITLE_VERSION_CATALOG_CHECK: 'Checking Version Catalog (libs.versions.toml)...',
  MSG_VERSION_CATALOG_FOUND: 'Version Catalog detected (gradle/libs.versions.toml).',
  MSG_VERSION_CATALOG_MISSING: 'Version Catalog not found. Falling back to direct dependency.',
  MSG_VERSION_CATALOG_CLIX_ALIAS_ADDED: 'Added Clix alias to Version Catalog.',
  MSG_VERSION_CATALOG_CLIX_ALIAS_EXISTS: 'Clix alias already present in Version Catalog.',
  MSG_VERSION_CATALOG_CLIX_ALIAS_FAILED: 'Failed to update Version Catalog with Clix alias.',

  // Google Services Plugin
  TITLE_GMS_PLUGIN_CHECK: 'Checking for Google Services plugin...',
  MSG_GMS_PLUGIN_FOUND: 'Google Services plugin detected.',
  MSG_GMS_PLUGIN_NOT_FOUND: 'Google Services plugin is missing.',
  FIX_GMS_PLUGIN: 'To resolve this, add the following to build.gradle(.kts):',
  MSG_GMS_PLUGIN_FIX_FAILURE: 'Automatic fix failed. Please add the following to app/build.gradle(.kts):',
  CODE_GMS_PLUGIN: `plugins {
\tid("com.google.gms.google-services") version "${ANDROID_GMS_PLUGIN_VERSION}"
}`,

  // Clix Application Import & Initialization
  TITLE_CLIX_INITIALIZATION_CHECK: 'Checking Clix SDK initialization...',
  MSG_MANIFEST_READ_FAIL: 'Unable to read AndroidManifest.xml.',
  MSG_APPLICATION_CLASS_NOT_DEFINED: 'Application class is not defined in AndroidManifest.xml.',
  MSG_APPLICATION_CLASS_MISSING: 'Application class is missing from expected locations.',
  MSG_APPLICATION_FILE_READ_FAIL: 'Unable to read Application class file.',
  MSG_CLIX_INIT_SUCCESS: 'Clix SDK initialization detected.',
  MSG_CLIX_INIT_MISSING: 'Clix SDK initialization is missing from the Application class.',
  MSG_APP_CREATE_SUCCESS: 'Fixed: Application class created.',
  FIX_CLIX_INITIALIZATION: 'To resolve this, follow the guide below:',
  MSG_APP_FIX_FAILURE: 'Automatic fix failed. Please follow the guide below:',
  CLIX_INITIALIZATION_LINK: 'https://docs.clix.so/sdk-quickstart-android#2-initialize-clix-with-config',

  // MainActivity Permissions
  TITLE_PERMISSION_CHECK: 'Checking permission request implementation...',
  MSG_MAIN_ACTIVITY_NOT_FOUND: 'MainActivity.java or MainActivity.kt is missing.',
  MSG_PERMISSION_FOUND: 'MainActivity contains permission request code.',
  MSG_PERMISSION_MISSING: 'MainActivity does not contain code to request permissions.',
  FIX_PERMISSION_REQUEST: 'To resolve this, follow the guide below:',
  MSG_PERMISSION_FIX_FAILURE: 'Automatic fix failed. Please follow the guide below:',
  PERMISSION_REQUEST_LINK: 'https://docs.clix.so/sdk-quickstart-android#2-initialize-clix-with-config',

  // google-services.json
  TITLE_GOOGLE_SERVICES_JSON_CHECK: 'Checking for google-services.json...',
  MSG_GOOGLE_JSON_MISSING: 'google-services.json is missing from app/google-services.json.',
  MSG_GOOGLE_JSON_FOUND: 'google-services.json is present.',
  FIX_GOOGLE_SERVICES_JSON: 'To resolve this, follow the guide below:',
  MSG_GOOGLE_JSON_FIX_FAILURE: 'Automatic fix failed. Please follow the guide below:',
  GOOGLE_SERVICES_JSON_LINK: 'https://docs.clix.so/firebase-setting',

  // Auto-fix messages
  MSG_AUTO_FIX_SUCCESS: 'Fixed: Changes applied automatically.',

  // General
  MSG_SOURCE_DIR_NOT_FOUND: 'Source directory is missing.',
  MSG_WORKING_DIRECTORY_NOT_FOUND: 'Working directory is missing.',
  MSG_APP_BUILD_GRADLE_NOT_FOUND: 'app/build.gradle(.kts) is missing.',
  MSG_APP_BUILD_GRADLE_READ_FAIL: 'Unable to read app/build.gradle(.kts).',
};
