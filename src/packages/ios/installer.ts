import { promises as fs } from 'fs';
import path from 'path';
import { findAppPath } from './locator.js';
import { configureXcodeProject } from './xcode-project.js';

/**
 * Detects whether the iOS project is using CocoaPods or Swift Package Manager (SPM)
 * @returns [usingSPM, usingCocoaPods]
 */
async function detectPackageManager(): Promise<[boolean, boolean]> {
  let usingCocoaPods = false;
  let usingSPM = false;

  // Check for Podfile which indicates CocoaPods
  try {
    await fs.access('Podfile');
    usingCocoaPods = true;
  } catch {
    // Podfile doesn't exist
  }

  // Check for Package.swift which indicates SPM
  try {
    await fs.access('Package.swift');
    usingSPM = true;
  } catch {
    // Package.swift doesn't exist
  }

  // Check for .xcodeproj files with SPM dependencies
  try {
    const files = await fs.readdir('.', { withFileTypes: true });

    for (const f of files) {
      if (f.name.endsWith('.xcodeproj')) {
        // Check if project.pbxproj contains SPM references
        const pbxprojPath = path.join(f.name, 'project.pbxproj');
        try {
          const data = await fs.readFile(pbxprojPath, 'utf-8');
          if (data.includes('XCRemoteSwiftPackageReference')) {
            usingSPM = true;
          }
        } catch {
          // Ignore errors reading pbxproj
        }
      }

      // Check for .xcworkspace which typically indicates CocoaPods
      if (f.name.endsWith('.xcworkspace') && !f.name.endsWith('xcodeproj.xcworkspace')) {
        usingCocoaPods = true;
      }
    }
  } catch {
    // Ignore errors
  }

  // If both are detected, prioritize the one that seems more actively used
  if (usingSPM && usingCocoaPods) {
    // Check if Podfile.lock exists, which indicates active use of CocoaPods
    try {
      await fs.access('Podfile.lock');
      // Podfile.lock exists, prioritize CocoaPods
      usingSPM = false;
      usingCocoaPods = true;
    } catch {
      // No Podfile.lock, prioritize SPM
      usingSPM = true;
      usingCocoaPods = false;
    }
  }

  return [usingSPM, usingCocoaPods];
}

/**
 * Prompts user for SPM usage
 * @returns User response
 */
function promptForSPM(): Promise<string> {
  return new Promise((resolve) => {
    process.stdout.write('Could not automatically detect package manager. Are you using Swift Package Manager (SPM)? (Y/n) ');

    process.stdin.once('data', (data) => {
      resolve(data.toString().trim());
    });
  });
}

/**
 * Waits for user to press Enter
 */
function waitForEnter(message: string = 'Press Enter to continue...'): Promise<void> {
  return new Promise((resolve) => {
    console.log(message);
    process.stdin.once('data', () => {
      resolve();
    });
  });
}

/**
 * Displays iOS installation instructions
 * @param projectID Clix project ID
 * @param apiKey Clix API key
 * @param verbose Enable verbose output
 * @param dryRun Perform dry run
 */
export async function displayIOSInstructions(
  projectID: string,
  apiKey: string,
  verbose: boolean = false,
  dryRun: boolean = false
): Promise<void> {
  // Automatically detect whether the project is using CocoaPods or SPM
  const [usingSPM, usingCocoaPods] = await detectPackageManager();

  if (usingSPM) {
    console.log('');
    console.log('📦 Swift Package Manager (SPM) detected!');
    console.log('📦 Please add the Clix SDK via SPM in Xcode:');
    console.log('  1. Open your Xcode project.');
    console.log('  2. Go to File > Add Package Dependencies');
    console.log('  3. Enter the URL below to the input on the right side');
    console.log('     https://github.com/clix-so/clix-ios-sdk.git');
    console.log('  4. Select \'Up to Next Major\' for the version rule');
    console.log('');
    await waitForEnter();
  } else if (usingCocoaPods) {
    console.log('');
    console.log('📦 CocoaPods detected!');
    console.log('🤖 Installing Clix SDK for iOS via CocoaPods');
    console.log('');
  } else {
    // If neither is detected, ask the user
    const useSPM = await promptForSPM();
    if (useSPM === '' || useSPM.toLowerCase() === 'y') {
      console.log('');
      console.log('📦 Please add the Clix SDK via SPM in Xcode:');
      console.log('  1. Open your Xcode project.');
      console.log('  2. Go to File > Add Package Dependencies');
      console.log('  3. Enter the URL below to the input on the right side');
      console.log('     https://github.com/clix-so/clix-ios-sdk.git');
      console.log('  4. Select \'Up to Next Major\' for the version rule');
      console.log('  5. Click \'Add Package\' to add the Clix SDK');
      console.log('  6. Add your main app to the target list');
      console.log('');
      await waitForEnter();
    } else {
      console.log('');
      console.log('🤖 Installing Clix SDK for iOS via CocoaPods');
      console.log('');
    }
  }

  console.log('');
  console.log('📱 Integrating Clix SDK for iOS...');
  console.log('');

  console.log('1️⃣  Notification Service Extension & App Group Setup');
  console.log('  1. Open your Xcode project.');
  console.log('  2. Go to File > New > Target');
  console.log('  3. Select \'Notification Service Extension\' and click Next.');
  console.log('  4. Name it \'NotificationServiceExtension\' and click Finish.');
  console.log('  5. When prompted to activate the scheme, click \'Don\'t Activate\'.');
  console.log('  6. Add Clix framework to NotificationServiceExtension target:');
  console.log('     a. Select \'NotificationServiceExtension\' target in the project navigator');
  console.log('     b. Go to \'General\' tab');
  console.log('     c. Under \'Frameworks, Libraries, and Embedded Content\', click \'+\'');
  console.log('     d. Search for and add \'Clix\' framework');
  console.log('     e. Ensure \'Embed & Sign\' is selected for the Clix framework');
  console.log('');
  await waitForEnter('Press Enter after you have created the NotificationServiceExtension...');

  console.log('');
  console.log('2️⃣  Configuring App Groups and NotificationServiceExtension');
  console.log('');
  console.log('🤖 Automating Xcode project configuration...');

  // Try to configure the Xcode project automatically
  try {
    await configureXcodeProject(projectID, verbose, dryRun);
    console.log('✅ Xcode project configured successfully!');
    console.log('  - App Groups capability added to main app target');
    console.log('  - Background Modes (\'Background fetch\', \'Remote notifications\') enabled on main app');
    console.log('  - App Groups capability added to NotificationServiceExtension target (if present)');
    console.log('  - Clix framework added to NotificationServiceExtension target with \'Embed & Sign\' (if present)');
    console.log('  - NotificationServiceExtension is now ready to handle Clix push notifications');
    console.log('');
    await waitForEnter();
  } catch (error) {
    console.log(`⚠️ Automatic configuration failed: ${error instanceof Error ? error.message : String(error)}`);
    console.log('⚙️ Switching to manual configuration...');
    // Fall back to manual configuration
    console.log('');
    console.log('2️⃣  App Group Configuration (Manual)');
    console.log('  1. Select your main app target.');
    console.log('  2. Go to the \'Signing & Capabilities\' tab.');
    console.log('  3. Click \'+\' to add a capability.');
    console.log('  4. Search for and add \'App Groups\'.');
    console.log('  5. Click \'+\' under App Groups to add a new group.');
    console.log(`  6. Enter 'group.clix.${projectID}' as the group name.`);
    console.log('');
    await waitForEnter('Press Enter after you have configured App Groups for the main app...');

    console.log('');
    console.log('3️⃣  NotificationServiceExtension Setup (Manual)');
    console.log('  1. Select the NotificationServiceExtension target.');
    console.log('  2. Go to the \'Signing & Capabilities\' tab.');
    console.log('  3. Add the App Groups capability.');
    console.log(`  4. Select the same group: 'group.clix.${projectID}'.`);
    console.log('');
    await waitForEnter('Press Enter after you have configured App Groups for the extension target...');

    console.log('');
    console.log('4️⃣  Update NotificationServiceExtension Dependencies (Manual)');
    console.log('  1. Select the NotificationServiceExtension target.');
    console.log('  2. Go to the \'General\' tab.');
    console.log('  3. Click \'+\' under \'Frameworks, Libraries, and Embedded Content\'.');
    console.log('  4. Search for and add \'Clix\' framework.');
    console.log('  5. Ensure \'Embed & Sign\' is selected for the Clix framework.');
    console.log('  6. Verify that Clix appears in the frameworks list for NotificationServiceExtension.');
    console.log('');
    await waitForEnter('Press Enter after you have configured everything for the extension target...');

    // Add manual Background Modes steps
    console.log('');
    console.log('5️⃣  Enable Background Modes on Main App (Manual)');
    console.log('  1. Select your MAIN app target.');
    console.log('  2. Open the \'Signing & Capabilities\' tab.');
    console.log('  3. Click \'+\' and add \'Background Modes\'.');
    console.log('  4. Check the boxes for:');
    console.log('     - Background fetch');
    console.log('     - Remote notifications');
    console.log('');
    await waitForEnter('Press Enter after you have enabled Background Modes...');
  }

  console.log('');
  console.log('🚀 Clix SDK iOS setup instructions complete!');
  console.log('');
  console.log('Run \'clix-cli doctor --ios\' to verify your setup.');
}

/**
 * Creates a new AppDelegate.swift file
 * @param projectId Clix project ID
 * @param apiKey Clix API key
 */
async function createAppDelegate(projectId: string, apiKey: string): Promise<void> {
  const template = `import UIKit
import Clix
import Firebase

class AppDelegate: ClixAppDelegate {
    override func application(_ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {

        FirebaseApp.configure()

        Task {
            await Clix.initialize(
                config: ClixConfig(
                    projectId: "${projectId}",
                    apiKey: "${apiKey}"
                )
            )
        }

        return super.application(application, didFinishLaunchingWithOptions: launchOptions)
    }
}
`;

  const appPath = await findAppPath();
  const appDelegatePath = path.join(appPath, 'AppDelegate.swift');

  await fs.writeFile(appDelegatePath, template, 'utf-8');

  // Locate and modify <YourProjectName>App.swift
  const projectDir = path.dirname(appDelegatePath);
  let appSwiftPath = '';

  const walk = async (dir: string): Promise<void> => {
    const entries = await fs.readdir(dir, { withFileTypes: true });

    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);

      if (entry.isDirectory()) {
        await walk(fullPath);
      } else if (entry.name.endsWith('App.swift')) {
        const content = await fs.readFile(fullPath, 'utf-8');
        if (content.includes('@main')) {
          appSwiftPath = fullPath;
          return;
        }
      }
    }
  };

  await walk(projectDir);

  if (!appSwiftPath) {
    // Could not find App.swift with @main
    return;
  }

  let content = await fs.readFile(appSwiftPath, 'utf-8');

  if (content.includes('@UIApplicationDelegateAdaptor(AppDelegate.self)')) {
    // Already contains the adaptor, no change needed
    return;
  }

  // Find struct declaration line with 'struct ...App: App'
  const lines = content.split('\n');
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const trimmed = line.trim();

    if (trimmed.startsWith('struct ') && trimmed.includes(': App') && content.includes('@main')) {
      // Insert '@UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate' just before first '{' in this line
      const idx = line.indexOf('{');
      if (idx !== -1) {
        let indent = '';
        for (const ch of line.substring(0, idx)) {
          if (ch === ' ' || ch === '\t') {
            indent += ch;
          } else {
            indent = '';
          }
        }
        const insertLine = indent + '    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate';
        // Insert after this line
        lines.splice(i + 1, 0, insertLine);
        content = lines.join('\n');
        break;
      }
    }
  }

  await fs.writeFile(appSwiftPath, content, 'utf-8');
}

/**
 * Updates the NotificationService.swift file
 * @param projectID Clix project ID
 * @returns Array of errors (empty if successful)
 */
export async function updateNotificationServiceExtension(projectID: string): Promise<Error[]> {
  const errors: Error[] = [];

  // Find the project path
  let projectPath: string;
  try {
    projectPath = await findAppPath();
  } catch (error) {
    errors.push(new Error(`Failed to find Xcode project: ${error instanceof Error ? error.message : String(error)}`));
    return errors;
  }

  // Assume NotificationServiceExtension is already added in Xcode
  // Get the directory one level above the project root
  const projectRoot = path.dirname(projectPath);
  const parentDir = path.dirname(projectRoot);
  const extensionDir = path.join(parentDir, 'NotificationServiceExtension');
  const serviceSwift = path.join(extensionDir, 'NotificationService.swift');
  const infoPlist = path.join(extensionDir, 'Info.plist');

  // Debug info for extension directory path
  console.log(`Looking for extension at: ${extensionDir}`);

  // Patch code if NotificationService.swift file exists
  try {
    await fs.access(serviceSwift);

    const serviceSwiftContent = `import Clix
import UserNotifications

/// NotificationService inherits all logic from ClixNotificationServiceExtension
/// No additional logic is needed unless you want to customize notification handling.
class NotificationService: ClixNotificationServiceExtension {

\t// Initialize with your Clix project ID
\toverride init() {
\t\tsuper.init()

\t\t// Register your Clix project ID
\t\tregister(projectId: "${projectID}")
\t}

\toverride func didReceive(
\t\t_ request: UNNotificationRequest,
\t\twithContentHandler contentHandler: @escaping (UNNotificationContent) -> Void
\t) {
\t\t// Call super to handle image downloading and send push received event
\t\tsuper.didReceive(request, withContentHandler: contentHandler)
\t}

\toverride func serviceExtensionTimeWillExpire() {
\t\tsuper.serviceExtensionTimeWillExpire()
\t}
}
`;

    await fs.writeFile(serviceSwift, serviceSwiftContent, 'utf-8');
    console.log('Created or updated NotificationService.swift');
  } catch (error) {
    errors.push(new Error(`Failed to write NotificationService.swift: ${error instanceof Error ? error.message : String(error)}`));
  }

  // Create or update Info.plist with NSAppTransportSecurity
  const infoPlistContent = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
\t<key>CFBundleDevelopmentRegion</key>
\t<string>$(DEVELOPMENT_LANGUAGE)</string>
\t<key>CFBundleDisplayName</key>
\t<string>NotificationServiceExtension</string>
\t<key>CFBundleExecutable</key>
\t<string>$(EXECUTABLE_NAME)</string>
\t<key>CFBundleIdentifier</key>
\t<string>$(PRODUCT_BUNDLE_IDENTIFIER)</string>
\t<key>CFBundleInfoDictionaryVersion</key>
\t<string>6.0</string>
\t<key>CFBundleName</key>
\t<string>$(PRODUCT_NAME)</string>
\t<key>CFBundlePackageType</key>
\t<string>$(PRODUCT_BUNDLE_PACKAGE_TYPE)</string>
\t<key>CFBundleShortVersionString</key>
\t<string>1.0</string>
\t<key>CFBundleVersion</key>
\t<string>1</string>
\t<key>NSExtension</key>
\t<dict>
\t\t<key>NSExtensionPointIdentifier</key>
\t\t<string>com.apple.usernotifications.service</string>
\t\t<key>NSExtensionPrincipalClass</key>
\t\t<string>$(PRODUCT_MODULE_NAME).NotificationService</string>
\t</dict>
\t<key>NSAppTransportSecurity</key>
\t<dict>
\t\t<key>NSAllowsArbitraryLoads</key>
\t\t<true/>
\t</dict>
</dict>
</plist>
`;

  // Check if Info.plist exists and update NSAppTransportSecurity if needed
  try {
    await fs.access(infoPlist);
    const content = await fs.readFile(infoPlist, 'utf-8');

    if (!content.includes('NSAppTransportSecurity')) {
      const insertKey = '<key>NSAppTransportSecurity</key><dict><key>NSAllowsArbitraryLoads</key><true/></dict>';
      const updated = content.replace('<dict>', '<dict>\n\t' + insertKey);
      await fs.writeFile(infoPlist, updated, 'utf-8');
      console.log('Inserted NSAppTransportSecurity into Info.plist');
    }
  } catch {
    // File doesn't exist, create it
    try {
      await fs.writeFile(infoPlist, infoPlistContent, 'utf-8');
      console.log('Created Info.plist with NSAppTransportSecurity');
    } catch (error) {
      errors.push(new Error(`Failed to write Info.plist: ${error instanceof Error ? error.message : String(error)}`));
    }
  }

  return errors;
}

/**
 * Installs Clix iOS SDK by modifying AppDelegate.swift
 * @param projectID Clix project ID
 * @param apiKey Clix API key
 */
export async function handleIOSInstall(projectID: string, apiKey: string): Promise<void> {
  // Store errors to display at the end
  const installErrors: string[] = [];

  const appPath = await findAppPath();
  const appDelegatePath = path.join(appPath, 'AppDelegate.swift');

  try {
    await fs.access(appDelegatePath);
  } catch {
    // If AppDelegate.swift not found, create one and return
    console.log('AppDelegate.swift not found, creating one...');
    await createAppDelegate(projectID, apiKey);
    return;
  }

  let content = await fs.readFile(appDelegatePath, 'utf-8');
  let updated = content;

  // 1. Add required imports
  if (!updated.includes('import Clix')) {
    updated = updated.replace('import UIKit', 'import UIKit\nimport Clix');
  }

  if (!updated.includes('import Firebase')) {
    // Add Firebase import after last import statement
    const lines = updated.split('\n');
    let insertIdx = 0;

    for (let i = 0; i < lines.length; i++) {
      if (lines[i].startsWith('import ')) {
        insertIdx = i + 1;
      }
    }

    lines.splice(insertIdx, 0, 'import Firebase');
    updated = lines.join('\n');
  }

  // 2. Update class declaration to inherit from ClixAppDelegate
  if (!updated.includes('ClixAppDelegate')) {
    const lines = updated.split('\n');

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];

      if (line.includes('class AppDelegate')) {
        // Replace the class declaration line
        let indent = '';
        for (const ch of line) {
          if (ch === ' ' || ch === '\t') {
            indent += ch;
          } else {
            break;
          }
        }
        lines[i] = indent + 'class AppDelegate: ClixAppDelegate {';
        break;
      }
    }

    updated = lines.join('\n');
  }

  // 3. Update didFinishLaunchingWithOptions method to include override keyword
  if (updated.includes('didFinishLaunchingWithOptions') && !updated.includes('override func application')) {
    const lines = updated.split('\n');

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      const trimmedLine = line.trim();

      if (trimmedLine.startsWith('func application')) {
        // Check if this line or subsequent lines contain didFinishLaunchingWithOptions
        let foundMethod = line.includes('didFinishLaunchingWithOptions');

        if (!foundMethod) {
          // Check next few lines for didFinishLaunchingWithOptions (multiline case)
          for (let j = i + 1; j < Math.min(lines.length, i + 5); j++) {
            if (lines[j].includes('didFinishLaunchingWithOptions')) {
              foundMethod = true;
              break;
            }
            // Stop searching if we hit another function or closing brace
            const nextTrimmed = lines[j].trim();
            if (nextTrimmed.startsWith('func ') || nextTrimmed === '}') {
              break;
            }
          }
        }

        if (foundMethod) {
          // Add override keyword
          let indent = '';
          for (const ch of line) {
            if (ch === ' ' || ch === '\t') {
              indent += ch;
            } else {
              break;
            }
          }
          lines[i] = indent + 'override ' + trimmedLine;
          break;
        }
      }
    }

    updated = lines.join('\n');
  }

  // 3.1. Add override keyword to other AppDelegate lifecycle methods
  const appDelegateMethods = [
    'applicationDidBecomeActive',
    'applicationWillResignActive',
    'applicationDidEnterBackground',
    'applicationWillEnterForeground',
    'applicationWillTerminate',
    'applicationDidReceiveMemoryWarning',
  ];

  for (const methodName of appDelegateMethods) {
    if (updated.includes(methodName) && !updated.includes(`override func ${methodName}`)) {
      const lines = updated.split('\n');

      for (let i = 0; i < lines.length; i++) {
        const line = lines[i];
        const trimmedLine = line.trim();

        if (trimmedLine.startsWith(`func ${methodName}`)) {
          // Check if this is the method we're looking for
          const foundMethod = line.includes(methodName);

          if (foundMethod) {
            // Add override keyword
            let indent = '';
            for (const ch of line) {
              if (ch === ' ' || ch === '\t') {
                indent += ch;
              } else {
                break;
              }
            }
            lines[i] = indent + 'override ' + trimmedLine;
            break;
          }
        }
      }

      updated = lines.join('\n');
    }
  }

  // 4. Add FirebaseApp.configure and Clix.initialize before return statement
  if (updated.includes('didFinishLaunchingWithOptions')) {
    // Check if Firebase is already configured
    const hasFirebaseConfig = updated.includes('FirebaseApp.configure()');
    const hasClixInit = updated.includes('Clix.initialize');

    // If we need to add either Firebase config or Clix init
    if (!hasFirebaseConfig || !hasClixInit) {
      // Find the return statement
      const lines = updated.split('\n');
      let returnLineIndex = -1;
      let returnStatement = '';

      for (let i = 0; i < lines.length; i++) {
        const trimmed = lines[i].trim();
        if (trimmed.startsWith('return ')) {
          returnLineIndex = i;
          returnStatement = trimmed;
          break;
        }
      }

      if (returnLineIndex !== -1) {
        // Get indentation from the return line
        let indent = '';
        for (const ch of lines[returnLineIndex]) {
          if (ch === ' ' || ch === '\t') {
            indent += ch;
          } else {
            break;
          }
        }

        // Build the insertion content
        let insertContent = '';

        // Add Firebase configuration if needed
        if (!hasFirebaseConfig) {
          insertContent += indent + 'FirebaseApp.configure()\n\n';
        }

        // Add Clix initialization if needed
        if (!hasClixInit) {
          insertContent += indent + 'Task {\n' +
            indent + '    await Clix.initialize(\n' +
            indent + '        config: ClixConfig(\n' +
            indent + `            projectId: "${projectID}",\n` +
            indent + `            apiKey: "${apiKey}"\n` +
            indent + '        )\n' +
            indent + '    )\n' +
            indent + '}\n\n';
        }

        // Replace the return line with our insertions + the original return statement
        lines[returnLineIndex] = insertContent + indent + returnStatement;
        updated = lines.join('\n');
      }
    }

    // Ensure super.application is called after Firebase and Clix initialization
    if (updated.includes('return true') && !updated.includes('return super.application')) {
      updated = updated.replace(
        'return true',
        'return super.application(application, didFinishLaunchingWithOptions: launchOptions)'
      );
    }
  }

  await fs.writeFile(appDelegatePath, updated, 'utf-8');
  console.log('✅ Clix SDK successfully integrated into AppDelegate.swift');

  // Report any errors that occurred during installation
  if (installErrors.length > 0) {
    console.log('\n⚠️ Some issues occurred during installation:');
    for (const err of installErrors) {
      console.log(' -', err);
    }
    console.log('\nPlease address these issues manually or contact support.');
    throw new Error('installation completed with some issues');
  }
}
