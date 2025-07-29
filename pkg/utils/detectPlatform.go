package utils

import (
	"fmt"
	"os"
	"strings"
)

// DetectPlatform detects the platform based on the files in the current directory
func DetectPlatform() (isIOS, isAndroid, isExpo bool) {
	files, err := os.ReadDir(".")
	if err != nil {
		return false, false, false
	}

	iosSignals := 0
	androidSignals := 0
	expoSignals := 0

	// Check for app.json (Expo indicator)
	appJSONFound := false
	packageJSONFound := false
	hasExpo := false

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

	// Determine if it's an Expo project
	if appJSONFound && hasExpo {
		expoSignals = 1
	}

	// Simple threshold-based judgment
	isIOS = iosSignals >= 1
	isAndroid = androidSignals >= 1
	isExpo = expoSignals >= 1

	// Prioritize Expo detection over native iOS/Android if both are present
	if isExpo {
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

	return
}
