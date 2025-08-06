package utils

import (
	"fmt"
	"os"
	"strings"
)

// DetectPlatform detects the platform based on the files in the current directory
func DetectPlatform() (isIOS, isAndroid, isExpo bool) {
	isIOS, isAndroid, isExpo, _ = DetectAllPlatforms()
	return
}

// DetectAllPlatforms detects all supported platforms
func DetectAllPlatforms() (isIOS, isAndroid, isExpo, isFlutter bool) {
	files, err := os.ReadDir(".")
	if err != nil {
		return false, false, false, false
	}

	iosSignals := 0
	androidSignals := 0
	expoSignals := 0
	flutterSignals := 0

	// Check for app.json (Expo indicator)
	appJSONFound := false
	packageJSONFound := false
	hasExpo := false
	pubspecFound := false

	for _, f := range files {
		name := f.Name()

		// Check for app.json
		if name == "app.json" {
			appJSONFound = true
		}

		// Check for package.json
		if name == "package.json" {
			packageJSONFound = true
		}

		// Check for pubspec.yaml (Flutter indicator)
		if name == "pubspec.yaml" {
			pubspecFound = true
		}

		// iOS
		if strings.HasSuffix(name, ".xcodeproj") ||
			strings.HasSuffix(name, ".xcworkspace") ||
			name == "Podfile" ||
			name == "Package.swift" ||
			name == "Info.plist" {
			iosSignals++
		}

		// Android
		if name == "build.gradle" ||
			name == "settings.gradle" ||
			name == "AndroidManifest.xml" ||
			name == "gradlew" {
			androidSignals++
		}

		// Flutter
		if name == "pubspec.yaml" ||
			name == "pubspec.lock" ||
			f.IsDir() && (name == "lib" || name == "test" || name == "android" || name == "ios") {
			flutterSignals++
		}
	}

	// Check if package.json contains expo
	if packageJSONFound {
		if data, err := os.ReadFile("package.json"); err == nil {
			packageContent := string(data)
			if strings.Contains(packageContent, "expo") {
				hasExpo = true
			}
		}
	}

	// Check if pubspec.yaml contains flutter
	if pubspecFound {
		if data, err := os.ReadFile("pubspec.yaml"); err == nil {
			pubspecContent := string(data)
			if strings.Contains(pubspecContent, "flutter:") || strings.Contains(pubspecContent, "flutter_test:") {
				flutterSignals++
			}
		}
	}

	// Determine if it's an Expo project
	if appJSONFound && hasExpo {
		expoSignals = 1
	}

	// Simple threshold-based judgment
	isIOS = iosSignals >= 1
	isAndroid = androidSignals >= 1
	isExpo = expoSignals >= 1
	isFlutter = flutterSignals >= 2 // Need at least pubspec.yaml and flutter dependency

	// Prioritize Flutter, then Expo detection over native iOS/Android if multiple are present
	if isFlutter {
		isIOS = false
		isAndroid = false
		isExpo = false
	} else if isExpo {
		isIOS = false
		isAndroid = false
	}

	if isIOS {
		fmt.Print("📦 iOS project detected\n\n")
	}
	if isAndroid {
		fmt.Print("📦 Android project detected\n\n")
	}
	if isExpo {
		fmt.Print("📦 React Native Expo project detected\n\n")
	}
	if isFlutter {
		fmt.Print("📦 Flutter project detected\n\n")
	}

	return
}
