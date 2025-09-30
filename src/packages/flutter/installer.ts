import { promises as fs } from 'fs';
import path from 'path';
import { runCommand } from '../../utils/shell.js';
import {
  FLUTTER_CLIX_SDK_VERSION,
  FLUTTER_FIREBASE_CORE_VERSION,
  FLUTTER_FIREBASE_MESSAGING_VERSION,
} from '../versions.js';

export async function handleFlutterInstall(apiKey: string, projectID: string): Promise<void> {
  const projectRoot = process.cwd();

  console.log('🚀 Installing Clix SDK for Flutter...');
  console.log('==================================================\n');

  // Check if this is a Flutter project
  if (!(await checkFlutterProject(projectRoot))) {
    console.error("This doesn't appear to be a Flutter project. Please ensure you're in the root of a Flutter project.");
    return;
  }

  // Step 1: Check and install Firebase CLI
  console.log('🔧 Checking Firebase CLI...');
  try {
    await checkAndInstallFirebaseCLI();
    console.log('✅ Firebase CLI is ready\n');
  } catch (error: any) {
    console.error('❌ Failed to setup Firebase CLI');
    console.log(`   ${error.message}`);
    return;
  }

  // Step 2: Check and install FlutterFire CLI
  console.log('🔥 Checking FlutterFire CLI...');
  try {
    await checkAndInstallFlutterFireCLI();
    console.log('✅ FlutterFire CLI is ready\n');
  } catch (error: any) {
    console.error('❌ Failed to setup FlutterFire CLI');
    console.log(`   ${error.message}`);
    return;
  }

  // Step 3: Configure Firebase project
  console.log('🚀 Configuring Firebase project...');
  try {
    await configureFirebaseProject();
    console.log('✅ Firebase project configured\n');
  } catch (error: any) {
    console.error('❌ Failed to configure Firebase project');
    console.log(`   ${error.message}`);
    console.log("⚠️  Please manually run 'flutterfire configure' to setup Firebase");
    return;
  }

  // Step 4: Add dependencies to pubspec.yaml
  console.log('📦 Adding dependencies to pubspec.yaml...');
  try {
    await addFlutterDependencies(projectRoot);
    console.log('✅ Dependencies added to pubspec.yaml\n');
  } catch (error: any) {
    console.error('❌ Failed to add dependencies to pubspec.yaml');
    console.log(`   ${error.message}`);
    console.log('Please manually add the following dependencies to your pubspec.yaml:');
    console.log('dependencies:');
    console.log(`  clix_flutter: ${FLUTTER_CLIX_SDK_VERSION}`);
    console.log(`  firebase_core: ${FLUTTER_FIREBASE_CORE_VERSION}`);
    console.log(`  firebase_messaging: ${FLUTTER_FIREBASE_MESSAGING_VERSION}`);
    return;
  }

  // Step 5: Install dependencies
  console.log('📱 Installing Flutter dependencies...');
  const pubGetResult = await runCommand('flutter', ['pub', 'get']);
  if (!pubGetResult.success) {
    console.error('❌ Failed to install Flutter dependencies');
    console.log('Run: flutter pub get');
    return;
  }
  console.log('✅ Flutter dependencies installed successfully\n');

  // Step 6: Verify Firebase configuration
  console.log('🔧 Verifying Firebase configuration...');
  try {
    await verifyFirebaseConfig(projectRoot);
    console.log('✅ Firebase configuration verified\n');
  } catch (error: any) {
    console.error('❌ Firebase configuration verification failed');
    console.log(`   ${error.message}`);
    return;
  }

  // Step 7: Update main.dart with Clix initialization
  console.log('🔗 Updating main.dart with Clix initialization...');
  try {
    await updateMainDart(projectRoot, projectID, apiKey);
    console.log('✅ main.dart updated successfully\n');
  } catch (error: any) {
    console.error('❌ Failed to update main.dart');
    console.log(`   ${error.message}`);
    console.log('⚠️  Please manually add the following to your main.dart:');
    console.log("   import 'package:firebase_core/firebase_core.dart';");
    console.log("   import 'package:clix_flutter/clix_flutter.dart';");
    console.log('   // Add Firebase.initializeApp() and Clix.initialize() in main() before runApp()');
  }

  // Step 8: iOS-specific setup instructions
  console.log('\n🍎 iOS-specific setup required:');
  console.log('==================================================');
  console.log('1. Open ios/Runner.xcworkspace in Xcode');
  console.log('2. Select Runner target > Signing & Capabilities');
  console.log("3. Add 'Push Notifications' capability");
  console.log("4. Add 'Background Modes' capability");
  console.log("5. Enable 'Remote notifications' in Background Modes");
  console.log('==================================================\n');

  // Step 9: Final instructions
  console.log('🎉 Clix SDK Flutter installation completed!');
  console.log('==================================================');
  console.log('Next steps:');
  console.log('1. Configure iOS push notifications in Xcode (as shown above)');
  console.log('2. Upload your iOS Service Account Key to Clix console');
  console.log("3. Run 'flutter run' to test your app");
  console.log("4. Run 'clix doctor --flutter' to verify your setup");
  console.log('==================================================');
}

async function checkFlutterProject(projectRoot: string): Promise<boolean> {
  const pubspecPath = path.join(projectRoot, 'pubspec.yaml');
  try {
    await fs.stat(pubspecPath);
  } catch {
    return false;
  }

  // Check if pubspec.yaml contains flutter dependency
  try {
    const data = await fs.readFile(pubspecPath, 'utf-8');
    return data.includes('flutter:') || data.includes('flutter_test:');
  } catch {
    return false;
  }
}

async function checkFirebaseConfig(projectRoot: string, platform: string): Promise<boolean> {
  let configPath: string;
  switch (platform) {
    case 'android':
      configPath = path.join(projectRoot, 'android', 'app', 'google-services.json');
      break;
    case 'ios':
      configPath = path.join(projectRoot, 'ios', 'Runner', 'GoogleService-Info.plist');
      break;
    default:
      return false;
  }

  try {
    await fs.stat(configPath);
    return true;
  } catch {
    return false;
  }
}

async function addFlutterDependencies(projectRoot: string): Promise<void> {
  const pubspecPath = path.join(projectRoot, 'pubspec.yaml');

  const data = await fs.readFile(pubspecPath, 'utf-8');
  const content = data;
  const lines = content.split('\n');
  const result: string[] = [];

  let dependenciesFound = false;
  let dependenciesAdded = false;

  const requiredDeps: Record<string, string> = {
    clix_flutter: FLUTTER_CLIX_SDK_VERSION,
    firebase_core: FLUTTER_FIREBASE_CORE_VERSION,
    firebase_messaging: FLUTTER_FIREBASE_MESSAGING_VERSION,
  };

  for (const line of lines) {
    result.push(line);

    // Find dependencies section
    if (line.trim() === 'dependencies:') {
      dependenciesFound = true;
    }

    // Add dependencies after finding the dependencies section
    if (dependenciesFound && !dependenciesAdded) {
      // Check if this line starts a new section (not indented under dependencies)
      if (
        line.trim() !== 'dependencies:' &&
        line.trim() !== '' &&
        !line.startsWith('  ') &&
        !line.startsWith('\t')
      ) {
        // We've moved to a new section, add dependencies before this line
        result.pop(); // Remove the current line

        // Add required dependencies
        for (const [dep, version] of Object.entries(requiredDeps)) {
          if (!content.includes(`${dep}:`)) {
            result.push(`  ${dep}: ${version}`);
          }
        }
        result.push(line); // Add back the current line
        dependenciesAdded = true;
      }
    }
  }

  // If dependencies section was found but we reached the end, add dependencies
  if (dependenciesFound && !dependenciesAdded) {
    for (const [dep, version] of Object.entries(requiredDeps)) {
      if (!content.includes(`${dep}:`)) {
        result.push(`  ${dep}: ${version}`);
      }
    }
  }

  // If no dependencies section found, throw error
  if (!dependenciesFound) {
    throw new Error('dependencies section not found in pubspec.yaml');
  }

  // Write updated pubspec.yaml
  const updatedContent = result.join('\n');
  await fs.writeFile(pubspecPath, updatedContent, 'utf-8');
}

async function updateMainDart(projectRoot: string, projectID: string, apiKey: string): Promise<void> {
  const mainDartPath = path.join(projectRoot, 'lib', 'main.dart');

  // Check if main.dart exists
  try {
    await fs.stat(mainDartPath);
  } catch {
    throw new Error(`main.dart not found at ${mainDartPath}`);
  }

  // Read existing main.dart
  const content = await fs.readFile(mainDartPath, 'utf-8');

  // Check if Clix is already initialized
  if (content.includes('Clix.initialize')) {
    return; // Already integrated
  }

  // Add necessary imports and initialization
  const modifiedContent = addClixToMainDart(content, projectID, apiKey);

  // Write the modified content back
  await fs.writeFile(mainDartPath, modifiedContent, 'utf-8');
}

function addClixToMainDart(content: string, projectID: string, apiKey: string): string {
  const lines = content.split('\n');
  const result: string[] = [];

  // Imports to add
  const firebaseImport = "import 'package:firebase_core/firebase_core.dart';";
  const firebaseOptionsImport = "import 'firebase_options.dart';";
  const clixImport = "import 'package:clix_flutter/clix_flutter.dart';";

  let importsAdded = false;
  let mainFunctionModified = false;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    // Add imports after existing import statements
    if (!importsAdded && line.trim().startsWith('import ')) {
      result.push(line);

      // Check if this is the last import
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
        if (!content.includes('firebase_core')) {
          result.push(firebaseImport);
        }
        if (!content.includes('firebase_options.dart')) {
          result.push(firebaseOptionsImport);
        }
        if (!content.includes('clix_flutter')) {
          result.push(clixImport);
        }
        importsAdded = true;
      }
    } else if (line.includes('void main()') && !mainFunctionModified) {
      // Modify main function to be async and add initialization
      if (line.includes('async')) {
        // Already async, just add the line
        result.push(line);
      } else {
        // Make it async
        const modifiedLine = line.replace('void main()', 'void main() async');
        result.push(modifiedLine);
      }

      // Add initialization code after the opening brace
      result.push('  WidgetsFlutterBinding.ensureInitialized();');
      result.push('');
      result.push('  await Firebase.initializeApp(');
      result.push('    options: DefaultFirebaseOptions.currentPlatform,');
      result.push('  );');
      result.push('');
      result.push('  await Clix.initialize(const ClixConfig(');
      result.push(`    projectId: '${projectID}',`);
      result.push(`    apiKey: '${apiKey}',`);
      result.push('  ));');
      result.push('');

      mainFunctionModified = true;
    } else {
      result.push(line);
    }
  }

  // If imports weren't added at the beginning, add them
  if (!importsAdded) {
    const imports = [firebaseImport, firebaseOptionsImport, clixImport, ''];
    result.unshift(...imports);
  }

  return result.join('\n');
}

async function checkAndInstallFirebaseCLI(): Promise<void> {
  // Check if Firebase CLI is already installed
  const checkResult = await runCommand('firebase', ['--version']);
  if (checkResult.success) {
    return;
  }

  console.log('Firebase CLI not found. Installing Firebase CLI...');

  // Try to install Firebase CLI via npm
  const installResult = await runCommand('npm', ['install', '-g', 'firebase-tools']);
  if (!installResult.success) {
    throw new Error('failed to install Firebase CLI via npm. Please install manually: npm install -g firebase-tools');
  }

  // Verify installation
  const verifyResult = await runCommand('firebase', ['--version']);
  if (!verifyResult.success) {
    throw new Error('Firebase CLI installation failed. Please install manually: npm install -g firebase-tools');
  }
}

async function checkAndInstallFlutterFireCLI(): Promise<void> {
  // Check if FlutterFire CLI is already installed
  const checkResult = await runCommand('flutterfire', ['--version']);
  if (checkResult.success) {
    return;
  }

  console.log('FlutterFire CLI not found. Installing FlutterFire CLI...');

  // Install FlutterFire CLI
  const installResult = await runCommand('dart', ['pub', 'global', 'activate', 'flutterfire_cli']);
  if (!installResult.success) {
    throw new Error('failed to install FlutterFire CLI. Please install manually: dart pub global activate flutterfire_cli');
  }

  // Verify installation
  const verifyResult = await runCommand('flutterfire', ['--version']);
  if (!verifyResult.success) {
    throw new Error(
      'FlutterFire CLI installation failed. Please ensure dart is in PATH and run: dart pub global activate flutterfire_cli'
    );
  }
}

async function configureFirebaseProject(): Promise<void> {
  // Check if firebase_options.dart already exists
  try {
    await fs.stat('lib/firebase_options.dart');
    console.log('Firebase project appears to be already configured (firebase_options.dart found)');
    return;
  } catch {
    // File doesn't exist, continue with configuration
  }

  console.log("Running 'flutterfire configure' to setup Firebase project...");
  console.log('Please follow the interactive prompts to:');
  console.log('1. Select your Firebase project');
  console.log('2. Choose platforms (iOS and Android)');
  console.log('3. Configure bundle IDs');

  // Run flutterfire configure interactively
  const configResult = await runCommand('flutterfire', ['configure']);
  if (!configResult.success) {
    throw new Error('flutterfire configure failed. Please run manually: flutterfire configure');
  }

  // Verify firebase_options.dart was created
  try {
    await fs.stat('lib/firebase_options.dart');
  } catch {
    throw new Error("firebase_options.dart was not created. Please run 'flutterfire configure' manually");
  }
}

async function verifyFirebaseConfig(projectRoot: string): Promise<void> {
  // Check if firebase_options.dart exists
  const firebaseOptionsPath = path.join(projectRoot, 'lib', 'firebase_options.dart');
  try {
    await fs.stat(firebaseOptionsPath);
  } catch {
    throw new Error("firebase_options.dart not found. Please run 'flutterfire configure' first");
  }

  // Check if Firebase configuration files exist in their proper locations
  const androidConfigPath = path.join(projectRoot, 'android', 'app', 'google-services.json');
  const iosConfigPath = path.join(projectRoot, 'ios', 'Runner', 'GoogleService-Info.plist');

  const missingFiles: string[] = [];

  try {
    await fs.stat(androidConfigPath);
  } catch {
    missingFiles.push('android/app/google-services.json');
  }

  try {
    await fs.stat(iosConfigPath);
  } catch {
    missingFiles.push('ios/Runner/GoogleService-Info.plist');
  }

  if (missingFiles.length > 0) {
    throw new Error(
      `missing Firebase config files: ${missingFiles.join(', ')}. These should be automatically created by 'flutterfire configure'`
    );
  }
}
