import { readdir, readFile } from 'fs/promises';
import { existsSync } from 'fs';

export interface DetectedPlatforms {
  isIOS: boolean;
  isAndroid: boolean;
  isExpo: boolean;
  isFlutter: boolean;
}

export async function detectPlatform(): Promise<{ isIOS: boolean; isAndroid: boolean; isExpo: boolean }> {
  const all = await detectAllPlatforms();
  return {
    isIOS: all.isIOS,
    isAndroid: all.isAndroid,
    isExpo: all.isExpo,
  };
}

export async function detectAllPlatforms(): Promise<DetectedPlatforms> {
  try {
    const files = await readdir('.');

    let iosSignals = 0;
    let androidSignals = 0;
    let expoSignals = 0;
    let flutterSignals = 0;

    let appJSONFound = false;
    let packageJSONFound = false;
    let hasExpo = false;
    let pubspecFound = false;

    for (const name of files) {
      // Check for app.json
      if (name === 'app.json') {
        appJSONFound = true;
      }

      // Check for package.json
      if (name === 'package.json') {
        packageJSONFound = true;
      }

      // Check for pubspec.yaml (Flutter indicator)
      if (name === 'pubspec.yaml') {
        pubspecFound = true;
      }

      // iOS
      if (
        name.endsWith('.xcodeproj') ||
        name.endsWith('.xcworkspace') ||
        name === 'Podfile' ||
        name === 'Package.swift' ||
        name === 'Info.plist'
      ) {
        iosSignals++;
      }

      // Android
      if (
        name === 'build.gradle' ||
        name === 'settings.gradle' ||
        name === 'AndroidManifest.xml' ||
        name === 'gradlew'
      ) {
        androidSignals++;
      }

      // Flutter
      if (
        name === 'pubspec.yaml' ||
        name === 'pubspec.lock'
      ) {
        flutterSignals++;
      }
    }

    // Check if package.json contains expo
    if (packageJSONFound && existsSync('package.json')) {
      try {
        const packageContent = await readFile('package.json', 'utf-8');
        if (packageContent.includes('expo')) {
          hasExpo = true;
        }
      } catch {}
    }

    // Check if pubspec.yaml contains flutter
    if (pubspecFound && existsSync('pubspec.yaml')) {
      try {
        const pubspecContent = await readFile('pubspec.yaml', 'utf-8');
        if (pubspecContent.includes('flutter:') || pubspecContent.includes('flutter_test:')) {
          flutterSignals++;
        }
      } catch {}
    }

    // Determine if it's an Expo project
    if (appJSONFound && hasExpo) {
      expoSignals = 1;
    }

    // Simple threshold-based judgment
    let isIOS = iosSignals >= 1;
    let isAndroid = androidSignals >= 1;
    let isExpo = expoSignals >= 1;
    let isFlutter = flutterSignals >= 2; // Need at least pubspec.yaml and flutter dependency

    // Prioritize Flutter, then Expo detection over native iOS/Android if multiple are present
    if (isFlutter) {
      isIOS = false;
      isAndroid = false;
      isExpo = false;
    } else if (isExpo) {
      isIOS = false;
      isAndroid = false;
    }

    if (isIOS) {
      console.log('📦 iOS project detected\n');
    }
    if (isAndroid) {
      console.log('📦 Android project detected\n');
    }
    if (isExpo) {
      console.log('📦 React Native Expo project detected\n');
    }
    if (isFlutter) {
      console.log('📦 Flutter project detected\n');
    }

    return { isIOS, isAndroid, isExpo, isFlutter };
  } catch (error) {
    return { isIOS: false, isAndroid: false, isExpo: false, isFlutter: false };
  }
}
