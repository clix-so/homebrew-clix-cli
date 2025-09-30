import { promises as fs } from 'fs';
import path from 'path';
import { runCommand } from '../../utils/shell.js';
import { prompt } from '../../utils/prompt.js';

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

export async function handleExpoInstall(apiKey: string, projectID: string): Promise<void> {
  const projectRoot = process.cwd();

  console.log('🚀 Installing Clix SDK for React Native Expo...');
  console.log('==================================================');

  // Check if this is an Expo project
  if (!(await checkExpoProject(projectRoot))) {
    console.error("This doesn't appear to be an Expo project. Please ensure you're in the root of an Expo project.");
    return;
  }

  // Step 1: Install expo-dev-client
  console.log('📦 Installing expo-dev-client...');
  const devClientResult = await runCommand('npx', ['expo', 'install', 'expo-dev-client']);
  if (!devClientResult.success) {
    console.error('❌ Failed to install expo-dev-client');
    console.log('Run: npx expo install expo-dev-client');
    return;
  }
  console.log('✅ expo-dev-client installed successfully\n');

  // Step 2: Install Firebase modules
  console.log('🔥 Installing Firebase modules...');
  const firebaseResult = await runCommand('npx', [
    'expo',
    'install',
    '@react-native-firebase/app',
    '@react-native-firebase/messaging',
    'expo-build-properties',
  ]);
  if (!firebaseResult.success) {
    console.error('❌ Failed to install Firebase modules');
    console.log('Run: npx expo install @react-native-firebase/app @react-native-firebase/messaging expo-build-properties');
    return;
  }
  console.log('✅ Firebase modules installed successfully\n');

  // Step 3: Install Clix dependencies
  console.log('📱 Installing Clix dependencies...');

  // Get appropriate MMKV version based on React Native version
  const mmkvVersion = await getMMKVVersion(projectRoot);
  if (!mmkvVersion) {
    console.error('❌ Failed to determine MMKV version');
    return;
  }

  const dependencies = [
    '@clix-so/react-native-sdk',
    '@notifee/react-native',
    'react-native-device-info',
    'react-native-get-random-values',
    mmkvVersion,
    'uuid',
  ];

  const clixResult = await runCommand('npx', ['expo', 'install', ...dependencies]);
  if (!clixResult.success) {
    console.error('❌ Failed to install Clix dependencies');
    console.log(`Run: npx expo install ${dependencies.join(' ')}`);
    return;
  }
  console.log('✅ Clix dependencies installed successfully\n');

  // Step 4: Check Firebase configuration files
  console.log('🔧 Checking Firebase configuration files...');
  const hasAndroidConfig = await checkFirebaseConfig(projectRoot, 'android');
  const hasIOSConfig = await checkFirebaseConfig(projectRoot, 'ios');

  if (!hasAndroidConfig || !hasIOSConfig) {
    console.error('❌ Firebase configuration files missing');
    if (!hasAndroidConfig) {
      console.log('   Missing: google-services.json (place in project root)');
    }
    if (!hasIOSConfig) {
      console.log('   Missing: GoogleService-Info.plist (place in project root)');
    }
    console.log('   Download these files from Firebase Console');
    return;
  }
  console.log('✅ Firebase configuration files found\n');

  // Step 5: Update app.json with Firebase plugin
  console.log('⚙️  Updating app.json configuration...');
  try {
    await updateAppConfig(projectRoot);
    console.log('✅ app.json updated successfully\n');
  } catch (error: any) {
    console.error('❌ Failed to update app.json');
    console.log(`   ${error.message}`);
    return;
  }

  // Step 6: Create Clix initialization file
  console.log('🔨 Creating Clix initialization...');
  try {
    await createClixInitialization(projectRoot, apiKey, projectID);
    console.log('✅ Clix initialization created successfully\n');
  } catch (error: any) {
    console.error('❌ Failed to create Clix initialization');
    console.log(`   ${error.message}`);
    return;
  }

  // Step 7: Integrate Clix initialization into App component
  console.log('🔗 Integrating Clix into App component...');
  try {
    await integrateClixIntoApp(projectRoot);
    console.log('✅ Clix integration added to App component\n');
  } catch (error: any) {
    console.error('❌ Failed to integrate Clix into App component');
    console.log(`   ${error.message}`);
    console.log('⚠️  Please manually add the following to your main component:');
    console.log("   import { initializeClix } from './clix-config';");
    console.log("   // Call initializeClix() in your component's useEffect");
    console.log('   // This should be added to App.tsx, App.js, or app/_layout.tsx\n');
  }

  // Step 8: Final instructions
  console.log('🎉 Clix SDK installation completed!');
  console.log('==================================================');
  console.log('Next steps:');
  console.log("1. Run 'npx expo prebuild --clean' to generate native code");
  console.log("2. Run 'npx expo run:android' or 'npx expo run:ios' to test");
  console.log("3. Run 'clix doctor --expo' to verify your setup");
  console.log('==================================================');
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

async function updateAppConfig(projectRoot: string): Promise<void> {
  const appJSONPath = path.join(projectRoot, 'app.json');

  const data = await fs.readFile(appJSONPath, 'utf-8');
  const config: AppConfig = JSON.parse(data);

  // Check if required plugins are already present
  let hasFirebaseAppPlugin = false;
  let hasFirebaseMessagingPlugin = false;
  let buildPropertiesPluginIndex = -1;

  for (let i = 0; i < config.expo.plugins.length; i++) {
    const plugin = config.expo.plugins[i];
    if (typeof plugin === 'string') {
      if (plugin === '@react-native-firebase/app') {
        hasFirebaseAppPlugin = true;
      }
      if (plugin === '@react-native-firebase/messaging') {
        hasFirebaseMessagingPlugin = true;
      }
      if (plugin === 'expo-build-properties') {
        buildPropertiesPluginIndex = i;
      }
    }
    if (Array.isArray(plugin) && plugin.length > 0) {
      if (plugin[0] === '@react-native-firebase/app') {
        hasFirebaseAppPlugin = true;
      }
      if (plugin[0] === '@react-native-firebase/messaging') {
        hasFirebaseMessagingPlugin = true;
      }
      if (plugin[0] === 'expo-build-properties') {
        buildPropertiesPluginIndex = i;
      }
    }
  }

  // Add missing plugins
  if (!hasFirebaseAppPlugin) {
    config.expo.plugins.push('@react-native-firebase/app');
  }

  if (!hasFirebaseMessagingPlugin) {
    config.expo.plugins.push('@react-native-firebase/messaging');
  }

  // Handle expo-build-properties plugin
  const notifeeRepo = '../../node_modules/@notifee/react-native/android/libs';

  if (buildPropertiesPluginIndex === -1) {
    // Plugin doesn't exist, add complete configuration
    const buildPropertiesPlugin = [
      'expo-build-properties',
      {
        ios: {
          useFrameworks: 'static',
        },
        android: {
          extraMavenRepos: [notifeeRepo],
        },
      },
    ];
    config.expo.plugins.push(buildPropertiesPlugin);
  } else {
    // Plugin exists, ensure it has the correct configuration
    const plugin = config.expo.plugins[buildPropertiesPluginIndex];

    // Handle string plugin format - convert to array format
    if (typeof plugin === 'string' && plugin === 'expo-build-properties') {
      config.expo.plugins[buildPropertiesPluginIndex] = [
        'expo-build-properties',
        {
          ios: {
            useFrameworks: 'static',
          },
          android: {
            extraMavenRepos: [notifeeRepo],
          },
        },
      ];
    } else if (Array.isArray(plugin) && plugin.length >= 1) {
      // Handle array plugin format
      let pluginConfig: Record<string, any>;
      if (plugin.length >= 2 && typeof plugin[1] === 'object') {
        pluginConfig = plugin[1];
      } else {
        pluginConfig = {};
        plugin.push(pluginConfig);
      }

      // Ensure iOS configuration
      if (pluginConfig.ios && typeof pluginConfig.ios === 'object') {
        pluginConfig.ios.useFrameworks = 'static';
      } else {
        pluginConfig.ios = { useFrameworks: 'static' };
      }

      // Ensure Android configuration
      if (pluginConfig.android && typeof pluginConfig.android === 'object') {
        // Handle extraMavenRepos
        if (pluginConfig.android.extraMavenRepos && Array.isArray(pluginConfig.android.extraMavenRepos)) {
          const hasNotifeeRepo = pluginConfig.android.extraMavenRepos.includes(notifeeRepo);
          if (!hasNotifeeRepo) {
            pluginConfig.android.extraMavenRepos.push(notifeeRepo);
          }
        } else {
          pluginConfig.android.extraMavenRepos = [notifeeRepo];
        }
      } else {
        pluginConfig.android = { extraMavenRepos: [notifeeRepo] };
      }

      config.expo.plugins[buildPropertiesPluginIndex] = plugin;
    }
  }

  // Update Android configuration
  if (!config.expo.android) {
    config.expo.android = {};
  }
  config.expo.android.googleServicesFile = './google-services.json';

  // Check for Android package name
  if (!config.expo.android.package) {
    console.log('\n📱 Android package name is required for Firebase configuration.');
    console.log('This should match the package name in your google-services.json file.');
    console.log('Example: com.yourcompany.yourapp');
    const androidPackage = await prompt('Enter your Android package name');
    if (androidPackage) {
      config.expo.android.package = androidPackage;
    }
  }

  // Update iOS configuration
  if (!config.expo.ios) {
    config.expo.ios = {};
  }
  config.expo.ios.googleServicesFile = './GoogleService-Info.plist';

  // Check for iOS bundle identifier
  if (!config.expo.ios.bundleIdentifier) {
    console.log('\n🍎 iOS bundle identifier is required for Firebase configuration.');
    console.log('This should match the bundle ID in your GoogleService-Info.plist file.');
    console.log('Example: com.yourcompany.yourapp');
    const iosBundleId = await prompt('Enter your iOS bundle identifier');
    if (iosBundleId) {
      config.expo.ios.bundleIdentifier = iosBundleId;
    }
  }

  // Add iOS push notification entitlements
  if (config.expo.ios.entitlements && typeof config.expo.ios.entitlements === 'object') {
    config.expo.ios.entitlements['aps-environment'] = 'production';
  } else {
    config.expo.ios.entitlements = {
      'aps-environment': 'production',
    };
  }

  // Add iOS info.plist settings for background modes
  if (config.expo.ios.infoPlist && typeof config.expo.ios.infoPlist === 'object') {
    config.expo.ios.infoPlist.UIBackgroundModes = ['remote-notification'];
  } else {
    config.expo.ios.infoPlist = {
      UIBackgroundModes: ['remote-notification'],
    };
  }

  // Write updated config back to file
  const updatedData = JSON.stringify(config, null, 2);
  await fs.writeFile(appJSONPath, updatedData, 'utf-8');
}

async function createClixInitialization(projectRoot: string, apiKey: string, projectID: string): Promise<void> {
  // Check if the project uses TypeScript
  let isTypeScript = false;
  try {
    await fs.stat(path.join(projectRoot, 'tsconfig.json'));
    isTypeScript = true;
  } catch {
    isTypeScript = false;
  }

  let fileName: string;
  let content: string;

  if (isTypeScript) {
    fileName = 'clix-config.ts';
    content = `import Clix from '@clix-so/react-native-sdk';

export const initializeClix = async (): Promise<void> => {
  try {
    await Clix.initialize({
      projectId: '${projectID}',
      apiKey: '${apiKey}'
    });
    console.log('Clix SDK initialized successfully');
  } catch (error) {
    console.error('Failed to initialize Clix SDK:', error);
  }
};
`;
  } else {
    fileName = 'clix-config.js';
    content = `import Clix from '@clix-so/react-native-sdk';

export const initializeClix = async () => {
  try {
    await Clix.initialize({
      projectId: '${projectID}',
      apiKey: '${apiKey}'
    });
    console.log('Clix SDK initialized successfully');
  } catch (error) {
    console.error('Failed to initialize Clix SDK:', error);
  }
};
`;
  }

  const filePath = path.join(projectRoot, fileName);
  await fs.writeFile(filePath, content, 'utf-8');
}

async function integrateClixIntoApp(projectRoot: string): Promise<void> {
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

  let appFilePath = '';
  for (const file of appFiles) {
    const fullPath = path.join(projectRoot, file);
    try {
      await fs.stat(fullPath);
      appFilePath = fullPath;
      break;
    } catch {
      continue;
    }
  }

  if (!appFilePath) {
    throw new Error('Could not find App.tsx, App.js, or _layout.tsx file in project');
  }

  // Read the existing App component
  const content = await fs.readFile(appFilePath, 'utf-8');

  // Check if Clix is already imported
  if (content.includes('initializeClix')) {
    return; // Already integrated
  }

  // Add Clix import and useEffect
  const modifiedContent = addClixToAppComponent(content);

  // Write the modified content back
  await fs.writeFile(appFilePath, modifiedContent, 'utf-8');
}

function addClixToAppComponent(content: string): string {
  const lines = content.split('\n');
  const result: string[] = [];

  const clixImport = "import { initializeClix } from './clix-config';";
  let importAdded = false;
  let initAdded = false;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    result.push(line);

    // Add Clix import after the last import
    if (!importAdded && line.trim().startsWith('import ')) {
      // Check if this is the last import line
      let isLastImport = true;
      for (let j = i + 1; j < lines.length; j++) {
        const nextLine = lines[j].trim();
        if (nextLine === '' || nextLine.startsWith('//')) {
          continue;
        }
        if (nextLine.startsWith('import ')) {
          isLastImport = false;
          break;
        }
        break;
      }

      if (isLastImport) {
        result.push(clixImport);
        importAdded = true;
      }
    }

    // Add useEffect import to React import if needed
    if (line.includes('import') && line.includes('react') && !line.includes('useEffect')) {
      if (line.includes('{') && line.includes('}')) {
        // Modify existing React import to include useEffect
        result[result.length - 1] = line.replace('}', ', useEffect }');
      } else if (!content.includes('useEffect')) {
        // Add separate useEffect import
        result.push("import { useEffect } from 'react';");
      }
    }

    // Add Clix initialization after component opening brace
    if (!initAdded && isComponentDeclaration(line)) {
      // Find the opening brace
      let braceIndex = i;
      if (!line.includes('{')) {
        for (let j = i + 1; j < lines.length; j++) {
          if (lines[j].includes('{')) {
            braceIndex = j;
            break;
          }
        }
      }

      if (braceIndex > i) {
        // We need to add to the next line after brace
        continue;
      } else {
        // Brace is on the same line, add initialization
        const initCode = [
          '',
          '  // Initialize Clix SDK',
          '  useEffect(() => {',
          '    initializeClix();',
          '  }, []);',
          '',
        ];
        result.push(...initCode);
        initAdded = true;
      }
    }

    // Handle case where opening brace is on next line
    if (!initAdded && line.trim() === '{' && i > 0) {
      const prevLine = lines[i - 1];
      if (isComponentDeclaration(prevLine)) {
        const initCode = ['', '  // Initialize Clix SDK', '  useEffect(() => {', '    initializeClix();', '  }, []);', ''];
        result.push(...initCode);
        initAdded = true;
      }
    }
  }

  // If import wasn't added, add at the beginning
  if (!importAdded) {
    result.unshift(clixImport, '');
  }

  return result.join('\n');
}

function isComponentDeclaration(line: string): boolean {
  const trimmed = line.trim();

  // Check for various component patterns
  const patterns = [
    'export default function',
    'const App',
    'function App',
    'const RootLayout',
    'function RootLayout',
    'const Layout',
    'function Layout',
  ];

  for (const pattern of patterns) {
    if (trimmed.includes(pattern)) {
      return true;
    }
  }

  // Check for export default function with any name
  if (trimmed.startsWith('export default function ')) {
    return true;
  }

  return false;
}

async function getMMKVVersion(projectRoot: string): Promise<string | null> {
  const packageJSONPath = path.join(projectRoot, 'package.json');
  try {
    const data = await fs.readFile(packageJSONPath, 'utf-8');
    const packageJSON = JSON.parse(data);

    // Get React Native version from dependencies
    let reactNativeVersion = '';
    if (packageJSON.dependencies && packageJSON.dependencies['react-native']) {
      reactNativeVersion = packageJSON.dependencies['react-native'];
    }

    if (!reactNativeVersion) {
      // Default to MMKV 2.x for safety if version cannot be determined
      return 'react-native-mmkv@^2.12.2';
    }

    // Parse React Native version
    const version = parseReactNativeVersion(reactNativeVersion);
    if (version === null) {
      // Default to MMKV 2.x for safety if version cannot be parsed
      return 'react-native-mmkv@^2.12.2';
    }

    // Determine MMKV version based on React Native version
    if (version >= 75) {
      // React Native 0.75+ - use MMKV 3.0.2+
      return 'react-native-mmkv@^3.0.2';
    } else if (version >= 74) {
      // React Native 0.74 - use MMKV 3.0.1
      return 'react-native-mmkv@^3.0.1';
    } else {
      // React Native < 0.74 - use MMKV 2.x
      return 'react-native-mmkv@^2.12.2';
    }
  } catch (error) {
    return 'react-native-mmkv@^2.12.2';
  }
}

function parseReactNativeVersion(versionStr: string): number | null {
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
