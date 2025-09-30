// iOS Package - Main entry point
// Export all public functions for iOS SDK integration

export { findAppPath } from './locator.js';
export { checkFirebaseIntegration, checkGoogleServicePlist } from './firebase-checks.js';
export { checkNotificationServiceExtension, extractProjectIDFromAppDelegate } from './notification-service.js';
export { ensureRuby, ensureXcodeproj, findXcodeProject, configureXcodeProject } from './xcode-project.js';
export { runIOSDoctor } from './doctor.js';
export { handleIOSInstall, displayIOSInstructions, updateNotificationServiceExtension } from './installer.js';
export { uninstallClixIOS } from './uninstaller.js';
