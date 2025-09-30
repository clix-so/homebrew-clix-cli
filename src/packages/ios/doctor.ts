import { promises as fs } from 'fs';
import path from 'path';
import { checkFirebaseIntegration, checkGoogleServicePlist } from './firebase-checks.js';

/**
 * Checks if we're in an Xcode project directory
 * @returns [projectPath, projectName] or throws error
 */
async function checkXcodeProject(): Promise<[string, string]> {
  const entries = await fs.readdir('.', { withFileTypes: true });

  for (const entry of entries) {
    if (entry.isDirectory() && entry.name.endsWith('.xcodeproj')) {
      const projectName = entry.name.replace('.xcodeproj', '');
      return [path.join('.', projectName), projectName];
    }
  }

  throw new Error('❌ No .xcodeproj found. Please run this command from the root of your Xcode project');
}

/**
 * Checks if AppDelegate.swift exists
 * @param projectPath Path to the project directory
 * @returns Path to AppDelegate.swift
 */
async function checkAppDelegateExists(projectPath: string): Promise<string> {
  const appDelegatePath = path.join(projectPath, 'AppDelegate.swift');
  await fs.access(appDelegatePath);
  return appDelegatePath;
}

/**
 * Checks required imports in AppDelegate.swift
 * @param appDelegatePath Path to AppDelegate.swift
 * @returns Array of error messages
 */
async function checkAppDelegateImports(appDelegatePath: string): Promise<string[]> {
  try {
    const content = await fs.readFile(appDelegatePath, 'utf-8');
    const errors: string[] = [];

    if (!content.includes('import Clix')) {
      errors.push('❌ Missing \'import Clix\' in AppDelegate.swift');
      errors.push('  └ Add \'import Clix\' at the top of your AppDelegate.swift file');
    }

    return errors;
  } catch (error) {
    return [`❌ Error reading AppDelegate.swift: ${error instanceof Error ? error.message : String(error)}`];
  }
}

/**
 * Checks for Clix.initialize call
 * @param appDelegatePath Path to AppDelegate.swift
 * @returns Array of error messages
 */
async function checkClixInitialization(appDelegatePath: string): Promise<string[]> {
  try {
    const content = await fs.readFile(appDelegatePath, 'utf-8');
    const errors: string[] = [];

    if (!content.includes('Clix.initialize')) {
      errors.push('❌ Missing \'Clix.initialize\' call in AppDelegate.swift');
      errors.push('  └ Add the following code in your didFinishLaunchingWithOptions method:');
      errors.push('  └ Clix.initialize(projectId: "YOUR_PROJECT_ID", username: "YOUR_USERNAME", password: "YOUR_PASSWORD")');
    }

    return errors;
  } catch (error) {
    return [`❌ Error reading AppDelegate.swift: ${error instanceof Error ? error.message : String(error)}`];
  }
}

/**
 * Checks if push notification capability is enabled
 * @param projectPath Path to the project directory
 * @param projectName Name of the project
 * @returns true if enabled, false otherwise
 */
async function checkPushCapabilities(projectPath: string, projectName: string): Promise<boolean> {
  // Find any .entitlements file in the project directory
  let entitlementsPath = '';

  try {
    const files = await fs.readdir(projectPath, { withFileTypes: true });
    for (const file of files) {
      if (!file.isDirectory() && file.name.endsWith('.entitlements')) {
        entitlementsPath = path.join(projectPath, file.name);
        break;
      }
    }
  } catch {
    // Ignore errors
  }

  // Fallback to project name if no .entitlements file found
  if (!entitlementsPath) {
    entitlementsPath = path.join(projectPath, `${projectName}.entitlements`);
  }

  try {
    await fs.access(entitlementsPath);
    const content = await fs.readFile(entitlementsPath, 'utf-8');
    return content.includes('aps-environment');
  } catch {
    return false;
  }
}

/**
 * Finds the main App.swift file for SwiftUI apps
 * @param projectPath Path to the project directory
 * @returns Path to App.swift file
 */
async function findAppSwiftFile(projectPath: string): Promise<string> {
  const walk = async (dir: string): Promise<string | null> => {
    const entries = await fs.readdir(dir, { withFileTypes: true });

    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);

      if (entry.isDirectory()) {
        const result = await walk(fullPath);
        if (result) return result;
      } else if (entry.name.endsWith('App.swift')) {
        const content = await fs.readFile(fullPath, 'utf-8');
        if (content.includes('@main')) {
          return fullPath;
        }
      }
    }

    return null;
  };

  const result = await walk(projectPath);
  if (!result) {
    throw new Error('No SwiftUI App.swift file found');
  }

  return result;
}

/**
 * Checks for UIApplicationDelegateAdaptor in SwiftUI apps
 * @param appSwiftPath Path to App.swift file
 * @returns Array of error messages
 */
async function checkSwiftUIIntegration(appSwiftPath: string): Promise<string[]> {
  try {
    const content = await fs.readFile(appSwiftPath, 'utf-8');
    const errors: string[] = [];

    if (!content.includes('@UIApplicationDelegateAdaptor')) {
      errors.push('❌ Missing @UIApplicationDelegateAdaptor in SwiftUI App');
      errors.push('  └ Add \'@UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate\' to your App struct');
    }

    return errors;
  } catch (error) {
    return [`❌ Error reading App.swift: ${error instanceof Error ? error.message : String(error)}`];
  }
}

/**
 * Helper function to sleep for a specified duration
 * @param ms Milliseconds to sleep
 */
const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * Performs a comprehensive check of the iOS project setup for Clix SDK
 */
export async function runIOSDoctor(): Promise<void> {
  console.log('====================================================');
  console.log('🔍 Starting Clix SDK doctor for iOS...');
  console.log('===================================================');

  process.stdout.write('[1/8] Checking Xcode project directory... ');
  await sleep(700);
  let projectPath: string;
  let projectName: string;
  try {
    [projectPath, projectName] = await checkXcodeProject();
    console.log(`✅ Found Xcode project: ${projectName}`);
  } catch (error) {
    console.log('❌');
    throw error;
  }

  process.stdout.write('[2/8] Checking AppDelegate.swift... ');
  await sleep(700);
  let appDelegatePath: string;
  try {
    appDelegatePath = await checkAppDelegateExists(projectPath);
    console.log('✅ AppDelegate.swift found.');
  } catch {
    console.log('❌ AppDelegate.swift not found');
    throw new Error('AppDelegate.swift not found');
  }

  process.stdout.write('[3/8] Checking required imports... ');
  await sleep(700);
  const importErrors = await checkAppDelegateImports(appDelegatePath);
  if (importErrors.length === 0) {
    console.log('✅ OK');
  } else {
    console.log('');
    importErrors.forEach(msg => console.log(msg));
  }

  process.stdout.write('[4/8] Checking Clix.initialize call... ');
  await sleep(700);
  const initErrors = await checkClixInitialization(appDelegatePath);
  if (initErrors.length === 0) {
    console.log('✅ OK');
  } else {
    console.log('');
    initErrors.forEach(msg => console.log(msg));
  }

  // Variable to store SwiftUI integration errors
  let swiftUIErrors: string[] = [];
  process.stdout.write('[5/8] Checking SwiftUI integration... ');
  await sleep(700);
  try {
    const appSwiftPath = await findAppSwiftFile(projectPath);
    swiftUIErrors = await checkSwiftUIIntegration(appSwiftPath);
    if (swiftUIErrors.length === 0) {
      console.log('✅ OK');
    } else {
      console.log('');
      swiftUIErrors.forEach(msg => console.log(msg));
    }
  } catch {
    console.log('(skipped: not a SwiftUI app)');
  }

  // Check push notification capability
  process.stdout.write('[6/8] Checking push notification capability... ');
  await sleep(700);
  let pushCapabilities: boolean;
  try {
    pushCapabilities = await checkPushCapabilities(projectPath, projectName);
    if (!pushCapabilities) {
      console.log('❌ \'aps-environment\' not set in entitlements file.');
    } else {
      console.log('✅ Push notification capability enabled.');
    }
  } catch {
    console.log('❌ Push notification capability not found or not enabled.');
    pushCapabilities = false;
  }

  // Check Firebase integration
  process.stdout.write('[7/8] Checking Firebase integration... ');
  await sleep(700);
  const firebaseErrors = await checkFirebaseIntegration(appDelegatePath);
  if (firebaseErrors.length === 0) {
    console.log('✅ OK');
  } else {
    console.log('');
    firebaseErrors.forEach(msg => console.log(msg));
  }

  // Check GoogleService-Info.plist
  process.stdout.write('[8/8] Checking GoogleService-Info.plist... ');
  await sleep(700);
  const plistError = await checkGoogleServicePlist(projectPath);
  if (plistError) {
    console.log(plistError.message);
  } else {
    console.log('✅ GoogleService-Info.plist found.');
  }

  console.log('====================================================');
  if (
    importErrors.length > 0 ||
    initErrors.length > 0 ||
    swiftUIErrors.length > 0 ||
    !pushCapabilities ||
    firebaseErrors.length > 0 ||
    plistError
  ) {
    console.log('⚠️ Some issues were found with your Clix SDK integration.');
    console.log('  └ Please fix the issues mentioned above to ensure proper push notification delivery.');
    console.log('  └ Run \'clix-cli install --ios\' to fix most issues automatically.');
  } else {
    console.log('🎉 Your iOS project is properly configured for Clix SDK!');
    console.log('  └ Push notifications should be working correctly.');
  }
  console.log('===================================================');
}
