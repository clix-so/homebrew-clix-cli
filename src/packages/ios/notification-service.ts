import { promises as fs } from 'fs';
import path from 'path';

/**
 * Extracts project ID from AppDelegate.swift by finding Clix.initialize call
 * @param appDelegatePath Path to AppDelegate.swift file
 * @returns Project ID or empty string if not found
 */
export function extractProjectIDFromAppDelegate(content: string): string {
  // Look for projectId in Clix.initialize
  const projectIDRegex = /projectId:\s*"([^"]*)"/;
  const matches = content.match(projectIDRegex);

  if (matches && matches[1]) {
    return matches[1];
  }

  return '';
}

/**
 * Checks if NotificationServiceExtension target and files exist
 * and verifies it has the correct setup with app groups and proper implementation
 * @param projectPath Path to the project directory
 * @returns Array of error messages (empty if no errors)
 */
export async function checkNotificationServiceExtension(projectPath: string): Promise<string[]> {
  const errors: string[] = [];

  // Check if the extension directory exists one level above the project root
  const projectRoot = path.dirname(projectPath);
  const parentDir = path.dirname(projectRoot);
  const extensionDir = path.join(parentDir, 'NotificationServiceExtension');
  const infoPlist = path.join(extensionDir, 'Info.plist');
  const serviceSwift = path.join(extensionDir, 'NotificationService.swift');

  try {
    await fs.access(extensionDir);
  } catch {
    errors.push('❌ NotificationServiceExtension directory not found.');
    errors.push('  └ Please add a Notification Service Extension target in Xcode.');
    return errors;
  }

  // Check if Info.plist exists
  let infoPlistExists = true;
  try {
    await fs.access(infoPlist);
  } catch {
    errors.push('❌ NotificationServiceExtension/Info.plist not found.');
    errors.push('  └ Please ensure Info.plist exists in NotificationServiceExtension.');
    infoPlistExists = false;
  }

  // Check if NotificationService.swift exists
  let serviceSwiftExists = true;
  try {
    await fs.access(serviceSwift);
  } catch {
    errors.push('❌ NotificationService.swift not found in NotificationServiceExtension.');
    errors.push('  └ Please ensure NotificationService.swift exists in NotificationServiceExtension.');
    serviceSwiftExists = false;
  }

  // Get project name for entitlement file paths
  const projectName = path.basename(projectPath);

  // Get project ID from AppDelegate.swift to check app group format
  let projectID = '';
  try {
    const appDelegateContent = await fs.readFile(path.join(projectPath, 'AppDelegate.swift'), 'utf-8');
    projectID = extractProjectIDFromAppDelegate(appDelegateContent);
  } catch {
    // Ignore if AppDelegate.swift cannot be read
  }

  // Check App Groups in both main app and extension targets
  const files = await fs.readdir(projectPath, { withFileTypes: true });
  let appEntitlements = '';

  for (const file of files) {
    if (!file.isDirectory() && file.name.endsWith('.entitlements')) {
      appEntitlements = path.join(projectPath, file.name);
      break;
    }
  }

  if (!appEntitlements) {
    // Fallback to project name if no .entitlements file found
    appEntitlements = path.join(projectPath, `${projectName}.entitlements`);
  }

  const extensionEntitlements = path.join(extensionDir, 'NotificationServiceExtension.entitlements');

  // Check if both app and extension have entitlements files
  let appEntitlementsExists = false;
  let extensionEntitlementsExists = false;

  try {
    await fs.access(appEntitlements);
    appEntitlementsExists = true;
  } catch {
    // File doesn't exist
  }

  try {
    await fs.access(extensionEntitlements);
    extensionEntitlementsExists = true;
  } catch {
    // File doesn't exist
  }

  // Check App Group Configuration if both entitlements files exist
  if (appEntitlementsExists && extensionEntitlementsExists) {
    try {
      const appEntitlementsContent = await fs.readFile(appEntitlements, 'utf-8');
      const extensionEntitlementsContent = await fs.readFile(extensionEntitlements, 'utf-8');

      // Extract app groups from XML format
      const extractGroups = (content: string): string[] => {
        const groups: string[] = [];

        // Try XML format first: <string>group.name</string>
        const appGroupPattern = /<key>\s*com\.apple\.security\.application-groups\s*<\/key>\s*<array>(.*?)<\/array>/s;
        const match = content.match(appGroupPattern);

        if (match && match[1]) {
          const stringPattern = /<string>(.*?)<\/string>/g;
          let stringMatch;
          while ((stringMatch = stringPattern.exec(match[1])) !== null) {
            if (stringMatch[1] && stringMatch[1].trim()) {
              groups.push(stringMatch[1].trim());
            }
          }
        }

        // Try alternative format: "group.name"
        if (groups.length === 0) {
          const altPattern = /"(group\.[^"]+)"/g;
          let altMatch;
          while ((altMatch = altPattern.exec(content)) !== null) {
            if (altMatch[1] && altMatch[1].trim()) {
              groups.push(altMatch[1].trim());
            }
          }
        }

        return groups;
      };

      const appGroupsFlat = extractGroups(appEntitlementsContent);
      const extensionGroupsFlat = extractGroups(extensionEntitlementsContent);

      if (appGroupsFlat.length === 0 || extensionGroupsFlat.length === 0) {
        errors.push('❌ App Groups not properly configured in entitlements files.');
        errors.push('  └ Please ensure both main app and extension have identical app groups.');
      } else {
        // Check if they share the same app group (intersection)
        const shared = appGroupsFlat.some(ag => extensionGroupsFlat.includes(ag));

        if (!shared) {
          errors.push('❌ App and Extension have different App Groups.');
          errors.push('  └ The app and extension must share at least one identical App Group.');
        }

        // Check app group format in main app
        if (projectID) {
          const expectedAppGroup = `group.clix.${projectID}`;
          const foundFormat = appGroupsFlat.includes(expectedAppGroup);

          if (!foundFormat) {
            errors.push(`❌ App Group doesn't follow the required format: ${expectedAppGroup}`);
            errors.push('  └ App Group should be in the format \'group.clix.{project_id}\'.');
          }
        }
      }
    } catch (error) {
      // Error reading entitlements files
    }
  } else {
    errors.push('❌ Missing entitlements files for app group configuration.');
    errors.push('  └ Both main app and extension need entitlements files with app groups.');
  }

  // Check Info.plist for NSAppTransportSecurity setting
  if (infoPlistExists) {
    try {
      const infoPlistContent = await fs.readFile(infoPlist, 'utf-8');

      if (!infoPlistContent.includes('NSAppTransportSecurity') || !infoPlistContent.includes('NSAllowsArbitraryLoads')) {
        errors.push('❌ NotificationServiceExtension Info.plist missing NSAppTransportSecurity configuration.');
        errors.push('  └ Please add the following to your Info.plist:');
        errors.push('  └ <key>NSAppTransportSecurity</key>');
        errors.push('  └   <dict>');
        errors.push('  └     <key>NSAllowsArbitraryLoads</key>');
        errors.push('  └     <true/>');
        errors.push('  └   </dict>');
      }
    } catch {
      // Ignore read errors
    }
  }

  // Check NotificationService.swift implementation
  if (serviceSwiftExists) {
    try {
      const serviceContent = await fs.readFile(serviceSwift, 'utf-8');

      // Check for proper imports
      if (!serviceContent.includes('import Clix')) {
        errors.push('❌ Missing \'import Clix\' in NotificationService.swift.');
        errors.push('  └ Please add \'import Clix\' to your NotificationService.swift.');
      }

      if (!serviceContent.includes('import UserNotifications')) {
        errors.push('❌ Missing \'import UserNotifications\' in NotificationService.swift.');
        errors.push('  └ Please add \'import UserNotifications\' to your NotificationService.swift.');
      }

      // Check for proper class inheritance
      if (!serviceContent.includes('class NotificationService: ClixNotificationServiceExtension')) {
        errors.push('❌ NotificationService doesn\'t inherit from ClixNotificationServiceExtension.');
        errors.push('  └ Please update your NotificationService class to inherit from ClixNotificationServiceExtension.');
      }

      // Check for project ID registration
      if (!serviceContent.includes('register(projectId:')) {
        errors.push('❌ Missing project ID registration in NotificationService.swift.');
        errors.push('  └ Please add \'register(projectId: "your-project-id")\' in your init() method.');
      }

      // Check for required overrides
      if (!serviceContent.includes('override func didReceive')) {
        errors.push('❌ Missing \'didReceive\' method override in NotificationService.swift.');
        errors.push('  └ Please implement the \'didReceive\' method that calls super.didReceive().');
      }

      if (!serviceContent.includes('override func serviceExtensionTimeWillExpire')) {
        errors.push('❌ Missing \'serviceExtensionTimeWillExpire\' method in NotificationService.swift.');
        errors.push('  └ Please implement the \'serviceExtensionTimeWillExpire\' method that calls super.serviceExtensionTimeWillExpire().');
      }
    } catch {
      // Ignore read errors
    }
  }

  return errors;
}
