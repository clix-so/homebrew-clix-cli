import { promises as fs } from 'fs';
import path from 'path';

// AppConfig represents the app.json configuration
interface AppConfig {
  expo: ExpoConfig;
}

interface ExpoConfig {
  name: string;
  slug: string;
  version: string;
  orientation: string;
  icon: string;
  userInterfaceStyle: string;
  splash: Record<string, any>;
  updates: Record<string, any>;
  assetBundlePatterns: string[];
  ios?: Record<string, any>;
  android?: Record<string, any>;
  web: Record<string, any>;
  plugins: any[];
}

export async function runExpoDoctor(): Promise<void> {
  const projectRoot = process.cwd();

  console.log('🔍 Running Clix Doctor for React Native Expo...');
  console.log('==================================================\n');

  const issues: string[] = [];

  // Check 1: Expo project structure
  if (!(await checkExpoProject(projectRoot))) {
    issues.push('Not an Expo project - missing app.json or expo dependency');
  } else {
    console.log('✅ Expo project detected');
  }

  // Check 2: Required dependencies
  const missingDeps = await checkDependencies(projectRoot);
  if (missingDeps.length > 0) {
    issues.push(`Missing dependencies: ${missingDeps.join(', ')}`);
  } else {
    console.log('✅ All required dependencies installed');
  }

  // Check 3: Firebase configuration files
  const hasAndroidConfig = await checkFirebaseConfig(projectRoot, 'android');
  const hasIOSConfig = await checkFirebaseConfig(projectRoot, 'ios');

  if (!hasAndroidConfig) {
    issues.push('Missing google-services.json file');
  } else {
    console.log('✅ google-services.json found');
  }

  if (!hasIOSConfig) {
    issues.push('Missing GoogleService-Info.plist file');
  } else {
    console.log('✅ GoogleService-Info.plist found');
  }

  // Check 4: app.json configuration
  const configIssues = await checkAppConfig(projectRoot);
  if (configIssues.length > 0) {
    issues.push(...configIssues);
  } else {
    console.log('✅ app.json properly configured');
  }

  // Check 5: Clix initialization file
  if (!(await checkClixInitialization(projectRoot))) {
    issues.push('Clix initialization file not found');
  } else {
    console.log('✅ Clix initialization file found');
  }

  // Check 6: Clix integration in App component
  if (!(await checkClixAppIntegration(projectRoot))) {
    issues.push('Clix not integrated in App component');
  } else {
    console.log('✅ Clix integrated in App component');
  }

  // Check 7: Generated native code
  if (!(await checkNativeCode(projectRoot))) {
    issues.push("Native code not generated - run 'npx expo prebuild --clean'");
  } else {
    console.log('✅ Native code generated');
  }

  // Report results
  console.log('');
  if (issues.length === 0) {
    console.log('🎉 All checks passed! Your Expo project is ready for Clix SDK.');
    console.log("You can now run 'npx expo run:android' or 'npx expo run:ios' to test push notifications.");
  } else {
    console.log(`❌ Found ${issues.length} issue(s):`);
    issues.forEach((issue, i) => {
      console.log(`  ${i + 1}. ${issue}`);
    });
    console.log('');
    console.log("Please fix the above issues and run 'clix doctor --expo' again.");
  }
}

async function checkExpoProject(projectRoot: string): Promise<boolean> {
  const appJSONPath = path.join(projectRoot, 'app.json');
  try {
    await fs.stat(appJSONPath);
  } catch {
    return false;
  }

  // Check if package.json exists and contains expo
  const packageJSONPath = path.join(projectRoot, 'package.json');
  try {
    const data = await fs.readFile(packageJSONPath, 'utf-8');
    return data.includes('expo');
  } catch {
    return false;
  }
}

async function checkFirebaseConfig(projectRoot: string, platform: string): Promise<boolean> {
  let fileName: string;
  switch (platform) {
    case 'android':
      fileName = 'google-services.json';
      break;
    case 'ios':
      fileName = 'GoogleService-Info.plist';
      break;
    default:
      return false;
  }

  const configPath = path.join(projectRoot, fileName);
  try {
    await fs.stat(configPath);
    return true;
  } catch {
    return false;
  }
}

async function checkDependencies(projectRoot: string): Promise<string[]> {
  const requiredDeps = [
    'expo-dev-client',
    '@react-native-firebase/app',
    '@react-native-firebase/messaging',
    'expo-build-properties',
    '@clix-so/react-native-sdk',
    '@notifee/react-native',
    'react-native-device-info',
    'react-native-get-random-values',
    'uuid',
  ];

  const packageJSONPath = path.join(projectRoot, 'package.json');
  try {
    const data = await fs.readFile(packageJSONPath, 'utf-8');
    const packageJSON = JSON.parse(data);

    const dependencies: Record<string, boolean> = {};
    if (packageJSON.dependencies) {
      Object.keys(packageJSON.dependencies).forEach((dep) => {
        dependencies[dep] = true;
      });
    }
    if (packageJSON.devDependencies) {
      Object.keys(packageJSON.devDependencies).forEach((dep) => {
        dependencies[dep] = true;
      });
    }

    const missing: string[] = [];
    for (const dep of requiredDeps) {
      if (!dependencies[dep]) {
        missing.push(dep);
      }
    }

    // Check MMKV separately as it needs version-specific validation
    if (!(await checkMMKVVersion(projectRoot, dependencies))) {
      missing.push('react-native-mmkv (incorrect version)');
    }

    return missing;
  } catch {
    return requiredDeps; // Return all as missing if can't read package.json
  }
}

async function checkAppConfig(projectRoot: string): Promise<string[]> {
  const issues: string[] = [];

  const appJSONPath = path.join(projectRoot, 'app.json');
  try {
    const data = await fs.readFile(appJSONPath, 'utf-8');
    const config: AppConfig = JSON.parse(data);

    // Check for Firebase App plugin
    let hasFirebaseAppPlugin = false;
    let hasFirebaseMessagingPlugin = false;
    for (const plugin of config.expo.plugins) {
      if (typeof plugin === 'string') {
        if (plugin === '@react-native-firebase/app') {
          hasFirebaseAppPlugin = true;
        }
        if (plugin === '@react-native-firebase/messaging') {
          hasFirebaseMessagingPlugin = true;
        }
      }
      if (Array.isArray(plugin) && plugin.length > 0) {
        if (plugin[0] === '@react-native-firebase/app') {
          hasFirebaseAppPlugin = true;
        }
        if (plugin[0] === '@react-native-firebase/messaging') {
          hasFirebaseMessagingPlugin = true;
        }
      }
    }

    if (!hasFirebaseAppPlugin) {
      issues.push('@react-native-firebase/app plugin not configured in app.json');
    }

    if (!hasFirebaseMessagingPlugin) {
      issues.push('@react-native-firebase/messaging plugin not configured in app.json');
    }

    // Check for expo-build-properties plugin and its configuration
    let hasBuildPropertiesPlugin = false;
    let hasIOSUseFrameworks = false;
    let hasAndroidExtraMavenRepos = false;

    for (const plugin of config.expo.plugins) {
      if (typeof plugin === 'string' && plugin === 'expo-build-properties') {
        hasBuildPropertiesPlugin = true;
        // String plugin format doesn't have configuration, so these are missing
      }
      if (Array.isArray(plugin) && plugin.length >= 2) {
        if (plugin[0] === 'expo-build-properties') {
          hasBuildPropertiesPlugin = true;

          // Check the plugin configuration
          const pluginConfig = plugin[1];
          if (pluginConfig && typeof pluginConfig === 'object') {
            // Check iOS useFrameworks
            if (pluginConfig.ios && typeof pluginConfig.ios === 'object') {
              if (pluginConfig.ios.useFrameworks === 'static') {
                hasIOSUseFrameworks = true;
              }
            }

            // Check Android extraMavenRepos
            if (pluginConfig.android && typeof pluginConfig.android === 'object') {
              const repos = pluginConfig.android.extraMavenRepos;
              const notifeeRepo = '../../node_modules/@notifee/react-native/android/libs';
              if (Array.isArray(repos)) {
                if (repos.includes(notifeeRepo)) {
                  hasAndroidExtraMavenRepos = true;
                }
              }
            }
          }
          break;
        }
      }
    }

    if (!hasBuildPropertiesPlugin) {
      issues.push('expo-build-properties plugin not configured in app.json');
    } else {
      if (!hasIOSUseFrameworks) {
        issues.push("iOS useFrameworks not set to 'static' in expo-build-properties");
      }
      if (!hasAndroidExtraMavenRepos) {
        issues.push('Android extraMavenRepos missing Notifee path in expo-build-properties');
      }
    }

    // Check for Firebase configuration in android and ios sections
    if (config.expo.android) {
      if (!config.expo.android.googleServicesFile) {
        issues.push('googleServicesFile not configured in app.json android section');
      }
      if (!config.expo.android.package) {
        issues.push('Android package name not configured in app.json');
      }
    } else {
      issues.push('Android configuration missing in app.json');
    }

    if (config.expo.ios) {
      if (!config.expo.ios.googleServicesFile) {
        issues.push('googleServicesFile not configured in app.json ios section');
      }
      if (!config.expo.ios.bundleIdentifier) {
        issues.push('iOS bundle identifier not configured in app.json');
      }

      // Check for iOS push notification settings
      if (config.expo.ios.entitlements && typeof config.expo.ios.entitlements === 'object') {
        if (!config.expo.ios.entitlements['aps-environment']) {
          issues.push('aps-environment not configured in iOS entitlements');
        }
      } else {
        issues.push('iOS entitlements not configured for push notifications');
      }

      if (config.expo.ios.infoPlist && typeof config.expo.ios.infoPlist === 'object') {
        if (!config.expo.ios.infoPlist.UIBackgroundModes) {
          issues.push('UIBackgroundModes not configured in iOS infoPlist');
        }
      } else {
        issues.push('iOS infoPlist not configured for background modes');
      }
    } else {
      issues.push('iOS configuration missing in app.json');
    }

    return issues;
  } catch {
    issues.push('Could not read or parse app.json');
    return issues;
  }
}

async function checkClixInitialization(projectRoot: string): Promise<boolean> {
  const tsPath = path.join(projectRoot, 'clix-config.ts');
  const jsPath = path.join(projectRoot, 'clix-config.js');

  try {
    await fs.stat(tsPath);
    return true;
  } catch {
    try {
      await fs.stat(jsPath);
      return true;
    } catch {
      return false;
    }
  }
}

async function checkNativeCode(projectRoot: string): Promise<boolean> {
  const androidPath = path.join(projectRoot, 'android');
  const iosPath = path.join(projectRoot, 'ios');

  let androidExists = false;
  let iosExists = false;

  try {
    const info = await fs.stat(androidPath);
    if (info.isDirectory()) {
      androidExists = true;
    }
  } catch {
    androidExists = false;
  }

  try {
    const info = await fs.stat(iosPath);
    if (info.isDirectory()) {
      iosExists = true;
    }
  } catch {
    iosExists = false;
  }

  return androidExists && iosExists;
}

async function checkClixAppIntegration(projectRoot: string): Promise<boolean> {
  // Common App component file paths in Expo projects
  const appFiles = [
    'App.tsx',
    'App.js',
    'src/App.tsx',
    'src/App.js',
    'app/_layout.tsx', // Expo Router
    'app/_layout.js', // Expo Router
    'src/app/_layout.tsx',
    'src/app/_layout.js',
  ];

  for (const file of appFiles) {
    const fullPath = path.join(projectRoot, file);
    try {
      const content = await fs.readFile(fullPath, 'utf-8');
      // Check if Clix is imported and initialized
      const hasImport = content.includes('initializeClix');
      const hasCall = content.includes('initializeClix()');
      return hasImport && hasCall;
    } catch {
      continue;
    }
  }

  return false;
}

async function checkMMKVVersion(projectRoot: string, dependencies: Record<string, boolean>): Promise<boolean> {
  // Check if MMKV is installed at all
  if (!dependencies['react-native-mmkv']) {
    return false;
  }

  // Get package.json to check versions
  const packageJSONPath = path.join(projectRoot, 'package.json');
  try {
    const data = await fs.readFile(packageJSONPath, 'utf-8');
    const packageJSON = JSON.parse(data);

    // Get installed package versions
    let reactNativeVersion = '';
    let mmkvVersion = '';
    if (packageJSON.dependencies) {
      if (packageJSON.dependencies['react-native']) {
        reactNativeVersion = packageJSON.dependencies['react-native'];
      }
      if (packageJSON.dependencies['react-native-mmkv']) {
        mmkvVersion = packageJSON.dependencies['react-native-mmkv'];
      }
    }

    if (!reactNativeVersion || !mmkvVersion) {
      return false;
    }

    // Parse React Native version
    const rnVersion = parseReactNativeVersionForDoctor(reactNativeVersion);
    if (rnVersion === null) {
      return false;
    }

    // Parse MMKV version to get major version
    const mmkvMajor = parseMMKVMajorVersion(mmkvVersion);
    if (mmkvMajor === null) {
      return false;
    }

    // Check version compatibility
    if (rnVersion >= 74 && mmkvMajor >= 3) {
      return true; // RN 0.74+ should use MMKV 3.x
    } else if (rnVersion < 74 && mmkvMajor === 2) {
      return true; // RN < 0.74 should use MMKV 2.x
    } else {
      return false; // Version mismatch
    }
  } catch {
    return false;
  }
}

function parseReactNativeVersionForDoctor(versionStr: string): number | null {
  // Remove common prefixes and suffixes
  let version = versionStr
    .replace(/^\^/, '')
    .replace(/^~/, '')
    .replace(/^>=/, '')
    .replace(/^<=/, '')
    .replace(/^>/, '')
    .replace(/^</, '');

  // Split by dots to get major.minor
  const parts = version.split('.');
  if (parts.length < 2) {
    return null;
  }

  const major = parts[0].trim();
  const minor = parts[1].trim();

  // Parse major version
  const majorInt = parseInt(major, 10);
  if (isNaN(majorInt)) {
    return null;
  }

  // Parse minor version
  const minorInt = parseInt(minor, 10);
  if (isNaN(minorInt)) {
    return null;
  }

  // Return as single integer (e.g., 0.74 -> 74, 0.75 -> 75)
  return majorInt * 100 + minorInt;
}

function parseMMKVMajorVersion(versionStr: string): number | null {
  // Remove common prefixes
  let version = versionStr
    .replace(/^\^/, '')
    .replace(/^~/, '')
    .replace(/^>=/, '')
    .replace(/^<=/, '')
    .replace(/^>/, '')
    .replace(/^</, '');

  // Split by dots to get major version
  const parts = version.split('.');
  if (parts.length < 1) {
    return null;
  }

  const major = parts[0].trim();

  // Parse major version
  const majorInt = parseInt(major, 10);
  if (isNaN(majorInt)) {
    return null;
  }

  return majorInt;
}
