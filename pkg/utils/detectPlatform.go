package utils

import (
	"fmt"
	"os"
	"strings"
)

// DetectPlatform detects the platform based on the files in the current directory
func DetectPlatform() (isIOS, isAndroid bool) {
	files, err := os.ReadDir(".")
	if err != nil {
		return false, false
	}

	iosSignals := 0
	androidSignals := 0

	for _, f := range files {
		name := f.Name()

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

	// Simple threshold-based judgment
	isIOS = iosSignals >= 1
	isAndroid = androidSignals >= 1

	if isIOS {
		fmt.Println("iOS project detected")
	}
	if isAndroid {
		fmt.Println("Android project detected")
	}

	return
}
