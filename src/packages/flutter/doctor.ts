import { promises as fs } from 'fs';
import path from 'path';
import { runCommand } from '../../utils/shell.js';
import {
  FLUTTER_CLIX_SDK_VERSION,
  FLUTTER_FIREBASE_CORE_VERSION,
  FLUTTER_FIREBASE_MESSAGING_VERSION,
} from '../versions.js';

export async function runFlutterDoctor(): Promise<void> {
  const projectRoot = process.cwd();

  console.log('🏥 Clix Flutter Doctor');
  console.log('==================================================\n');

  const allChecks: Array<{
    name: string;
    fn: (projectRoot: string) => Promise<{ passed: boolean; message: string }>;
  }> = [
    { name: 'Flutter Project Detection', fn: checkFlutterProject },
    { name: 'Firebase CLI Installation', fn: checkFirebaseCLI },
    { name: 'FlutterFire CLI Installation', fn: checkFlutterFireCLI },
    { name: 'Firebase Options Configuration', fn: checkFirebaseOptions },
    { name: 'Clix Flutter SDK Dependency', fn: checkClixDependency },
    { name: 'Firebase Core Dependency', fn: checkFirebaseCoreDependency },
    { name: 'Firebase Messaging Dependency', fn: checkFirebaseMessagingDependency },
    { name: 'Main.dart Configuration', fn: checkMainDartConfiguration },
  ];

  let allPassed = true;

  for (const check of allChecks) {
    const { passed, message } = await check.fn(projectRoot);
    if (passed) {
      console.log(`✅ ${check.name}`);
    } else {
      console.log(`❌ ${check.name}`);
      if (message) {
        console.log(`   ${message}`);
      }
      allPassed = false;
    }
  }

  console.log('');

  if (allPassed) {
    console.log('🎉 All checks passed! Your Flutter project is properly configured for Clix SDK.');
  } else {
    console.log("❗ Some checks failed. Please fix the issues above and run 'clix doctor --flutter' again.");
  }
}

async function checkFlutterProject(projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const pubspecPath = path.join(projectRoot, 'pubspec.yaml');
  try {
    await fs.stat(pubspecPath);
  } catch {
    return { passed: false, message: 'pubspec.yaml not found - not a Flutter project' };
  }

  try {
    const data = await fs.readFile(pubspecPath, 'utf-8');
    if (!data.includes('flutter:')) {
      return { passed: false, message: 'flutter dependency not found in pubspec.yaml' };
    }
  } catch {
    return { passed: false, message: 'failed to read pubspec.yaml' };
  }

  return { passed: true, message: '' };
}

async function checkClixDependency(projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const pubspecPath = path.join(projectRoot, 'pubspec.yaml');
  try {
    const data = await fs.readFile(pubspecPath, 'utf-8');
    if (!data.includes('clix_flutter:')) {
      return {
        passed: false,
        message: `Add 'clix_flutter: ${FLUTTER_CLIX_SDK_VERSION}' to dependencies in pubspec.yaml`,
      };
    }
  } catch {
    return { passed: false, message: 'failed to read pubspec.yaml' };
  }

  return { passed: true, message: '' };
}

async function checkFirebaseCoreDependency(projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const pubspecPath = path.join(projectRoot, 'pubspec.yaml');
  try {
    const data = await fs.readFile(pubspecPath, 'utf-8');
    if (!data.includes('firebase_core:')) {
      return {
        passed: false,
        message: `Add 'firebase_core: ${FLUTTER_FIREBASE_CORE_VERSION}' to dependencies in pubspec.yaml`,
      };
    }
  } catch {
    return { passed: false, message: 'failed to read pubspec.yaml' };
  }

  return { passed: true, message: '' };
}

async function checkFirebaseMessagingDependency(projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const pubspecPath = path.join(projectRoot, 'pubspec.yaml');
  try {
    const data = await fs.readFile(pubspecPath, 'utf-8');
    if (!data.includes('firebase_messaging:')) {
      return {
        passed: false,
        message: `Add 'firebase_messaging: ${FLUTTER_FIREBASE_MESSAGING_VERSION}' to dependencies in pubspec.yaml`,
      };
    }
  } catch {
    return { passed: false, message: 'failed to read pubspec.yaml' };
  }

  return { passed: true, message: '' };
}

async function checkFirebaseCLI(_projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const result = await runCommand('firebase', ['--version']);
  if (!result.success) {
    return { passed: false, message: 'Firebase CLI not installed. Run: npm install -g firebase-tools' };
  }
  return { passed: true, message: '' };
}

async function checkFlutterFireCLI(_projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const result = await runCommand('flutterfire', ['--version']);
  if (!result.success) {
    return { passed: false, message: 'FlutterFire CLI not installed. Run: dart pub global activate flutterfire_cli' };
  }
  return { passed: true, message: '' };
}

async function checkFirebaseOptions(projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const configPath = path.join(projectRoot, 'lib', 'firebase_options.dart');
  try {
    await fs.stat(configPath);
  } catch {
    return { passed: false, message: 'firebase_options.dart not found. Run: flutterfire configure' };
  }

  // Check if Firebase config files exist (they should be created by flutterfire configure)
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
    return { passed: false, message: `missing Firebase config files: ${missingFiles.join(', ')}. Run: flutterfire configure` };
  }

  return { passed: true, message: '' };
}

async function checkMainDartConfiguration(projectRoot: string): Promise<{ passed: boolean; message: string }> {
  const mainPath = path.join(projectRoot, 'lib', 'main.dart');
  try {
    await fs.stat(mainPath);
  } catch {
    return { passed: false, message: 'main.dart not found at lib/main.dart' };
  }

  try {
    const data = await fs.readFile(mainPath, 'utf-8');
    const content = data;
    const issues: string[] = [];

    if (!content.includes('firebase_core')) {
      issues.push('Missing firebase_core import');
    }

    if (!content.includes('firebase_options.dart')) {
      issues.push('Missing firebase_options.dart import');
    }

    if (!content.includes('clix_flutter')) {
      issues.push('Missing clix_flutter import');
    }

    if (!content.includes('Firebase.initializeApp')) {
      issues.push('Missing Firebase.initializeApp call');
    }

    if (!content.includes('DefaultFirebaseOptions.currentPlatform')) {
      issues.push('Missing DefaultFirebaseOptions.currentPlatform in Firebase.initializeApp');
    }

    if (!content.includes('Clix.initialize')) {
      issues.push('Missing Clix.initialize call');
    }

    if (!content.includes('WidgetsFlutterBinding.ensureInitialized()')) {
      issues.push('Missing WidgetsFlutterBinding.ensureInitialized() call');
    }

    if (issues.length > 0) {
      return { passed: false, message: issues.join('; ') };
    }

    return { passed: true, message: '' };
  } catch {
    return { passed: false, message: 'failed to read main.dart' };
  }
}
