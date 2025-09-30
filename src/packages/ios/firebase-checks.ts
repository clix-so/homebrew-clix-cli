import { promises as fs } from 'fs';
import path from 'path';

/**
 * Checks for 'import Firebase' and 'FirebaseApp.configure()' in AppDelegate.swift
 * @param appDelegatePath Path to AppDelegate.swift file
 * @returns Array of error messages (empty if no errors)
 */
export async function checkFirebaseIntegration(appDelegatePath: string): Promise<string[]> {
  try {
    const content = await fs.readFile(appDelegatePath, 'utf-8');
    const errors: string[] = [];

    if (!content.includes('import Firebase')) {
      errors.push('❌ Missing \'import Firebase\' in AppDelegate.swift');
      errors.push('  └ Add \'import Firebase\' at the top of your AppDelegate.swift file');
    }

    if (!content.includes('FirebaseApp.configure')) {
      errors.push('❌ Missing \'FirebaseApp.configure\' call in AppDelegate.swift');
      errors.push('  └ Add \'FirebaseApp.configure()\' or \'FirebaseApp.configure(options:)\' in your didFinishLaunchingWithOptions method');
    }

    return errors;
  } catch (error) {
    return [`❌ Error reading AppDelegate.swift: ${error instanceof Error ? error.message : String(error)}`];
  }
}

/**
 * Checks if GoogleService-Info.plist exists in the project directory
 * @param projectPath Path to the project directory
 * @returns Error or null if file exists
 */
export async function checkGoogleServicePlist(projectPath: string): Promise<Error | null> {
  const plistPath = path.join(projectPath, 'GoogleService-Info.plist');

  try {
    await fs.access(plistPath);
    return null;
  } catch {
    return new Error('❌ GoogleService-Info.plist not found in project directory.\n  └ Download it from Firebase Console and add it to your Xcode project root.');
  }
}
