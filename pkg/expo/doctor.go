package expo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
)

// RunDoctor performs comprehensive checks for React Native Expo Clix SDK setup
func RunDoctor() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	fmt.Println("🔍 Running Clix Doctor for React Native Expo...")
	logx.Separatorln()

	var issues []string

	// Check 1: Expo project structure
	if !CheckExpoProject(projectRoot) {
		issues = append(issues, "Not an Expo project - missing app.json or expo dependency")
	} else {
		logx.Log().Success().Println("✅ Expo project detected")
	}

	// Check 2: Required dependencies
	missingDeps := CheckDependencies(projectRoot)
	if len(missingDeps) > 0 {
		issues = append(issues, fmt.Sprintf("Missing dependencies: %s", strings.Join(missingDeps, ", ")))
	} else {
		logx.Log().Success().Println("✅ All required dependencies installed")
	}

	// Check 3: Firebase configuration files
	hasAndroidConfig := CheckFirebaseConfig(projectRoot, "android")
	hasIOSConfig := CheckFirebaseConfig(projectRoot, "ios")
	
	if !hasAndroidConfig {
		issues = append(issues, "Missing google-services.json file")
	} else {
		logx.Log().Success().Println("✅ google-services.json found")
	}

	if !hasIOSConfig {
		issues = append(issues, "Missing GoogleService-Info.plist file")
	} else {
		logx.Log().Success().Println("✅ GoogleService-Info.plist found")
	}

	// Check 4: app.json configuration
	configIssues := CheckAppConfig(projectRoot)
	if len(configIssues) > 0 {
		issues = append(issues, configIssues...)
	} else {
		logx.Log().Success().Println("✅ app.json properly configured")
	}

	// Check 5: Clix initialization file
	if !CheckClixInitialization(projectRoot) {
		issues = append(issues, "Clix initialization file not found")
	} else {
		logx.Log().Success().Println("✅ Clix initialization file found")
	}

	// Check 6: Generated native code
	if !CheckNativeCode(projectRoot) {
		issues = append(issues, "Native code not generated - run 'npx expo prebuild --clean'")
	} else {
		logx.Log().Success().Println("✅ Native code generated")
	}

	// Report results
	logx.NewLine()
	if len(issues) == 0 {
		logx.Log().Success().Println("🎉 All checks passed! Your Expo project is ready for Clix SDK.")
		fmt.Println("You can now run 'npx expo run:android' or 'npx expo run:ios' to test push notifications.")
	} else {
		fmt.Printf("❌ Found %d issue(s):\n", len(issues))
		for i, issue := range issues {
			fmt.Printf("  %d. %s\n", i+1, issue)
		}
		logx.NewLine()
		fmt.Println("Please fix the above issues and run 'clix doctor --expo' again.")
	}

	return nil
}

// CheckDependencies checks if all required dependencies are installed
func CheckDependencies(projectRoot string) []string {
	requiredDeps := []string{
		"expo-dev-client",
		"@react-native-firebase/app",
		"@react-native-firebase/messaging",
		"expo-build-properties",
		"@clix-so/react-native-sdk",
		"@notifee/react-native",
		"react-native-device-info",
		"react-native-get-random-values",
		"react-native-mmkv",
		"uuid",
	}

	packageJSONPath := filepath.Join(projectRoot, "package.json")
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return requiredDeps // Return all as missing if can't read package.json
	}

	var packageJSON map[string]any
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return requiredDeps
	}

	dependencies := make(map[string]bool)
	if deps, ok := packageJSON["dependencies"].(map[string]any); ok {
		for dep := range deps {
			dependencies[dep] = true
		}
	}
	if devDeps, ok := packageJSON["devDependencies"].(map[string]any); ok {
		for dep := range devDeps {
			dependencies[dep] = true
		}
	}

	var missing []string
	for _, dep := range requiredDeps {
		if !dependencies[dep] {
			missing = append(missing, dep)
		}
	}

	return missing
}

// CheckAppConfig checks if app.json is properly configured
func CheckAppConfig(projectRoot string) []string {
	var issues []string

	appJSONPath := filepath.Join(projectRoot, "app.json")
	data, err := os.ReadFile(appJSONPath)
	if err != nil {
		issues = append(issues, "Could not read app.json")
		return issues
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		issues = append(issues, "Could not parse app.json")
		return issues
	}

	// Check for Firebase App plugin
	hasFirebaseAppPlugin := false
	hasFirebaseMessagingPlugin := false
	for _, plugin := range config.Expo.Plugins {
		if pluginStr, ok := plugin.(string); ok {
			if pluginStr == "@react-native-firebase/app" {
				hasFirebaseAppPlugin = true
			}
			if pluginStr == "@react-native-firebase/messaging" {
				hasFirebaseMessagingPlugin = true
			}
		}
		if pluginArray, ok := plugin.([]any); ok && len(pluginArray) > 0 {
			if pluginStr, ok := pluginArray[0].(string); ok {
				if pluginStr == "@react-native-firebase/app" {
					hasFirebaseAppPlugin = true
				}
				if pluginStr == "@react-native-firebase/messaging" {
					hasFirebaseMessagingPlugin = true
				}
			}
		}
	}

	if !hasFirebaseAppPlugin {
		issues = append(issues, "@react-native-firebase/app plugin not configured in app.json")
	}

	if !hasFirebaseMessagingPlugin {
		issues = append(issues, "@react-native-firebase/messaging plugin not configured in app.json")
	}

	// Check for expo-build-properties plugin
	hasBuildPropertiesPlugin := false
	for _, plugin := range config.Expo.Plugins {
		if pluginStr, ok := plugin.(string); ok && pluginStr == "expo-build-properties" {
			hasBuildPropertiesPlugin = true
			break
		}
		if pluginArray, ok := plugin.([]any); ok && len(pluginArray) > 0 {
			if pluginStr, ok := pluginArray[0].(string); ok && pluginStr == "expo-build-properties" {
				hasBuildPropertiesPlugin = true
				break
			}
		}
	}

	if !hasBuildPropertiesPlugin {
		issues = append(issues, "expo-build-properties plugin not configured in app.json")
	}

	// Check for Firebase configuration in android and ios sections
	if config.Expo.Android != nil {
		if _, exists := config.Expo.Android["googleServicesFile"]; !exists {
			issues = append(issues, "googleServicesFile not configured in app.json android section")
		}
		if _, exists := config.Expo.Android["package"]; !exists {
			issues = append(issues, "Android package name not configured in app.json")
		}
	} else {
		issues = append(issues, "Android configuration missing in app.json")
	}

	if config.Expo.IOS != nil {
		if _, exists := config.Expo.IOS["googleServicesFile"]; !exists {
			issues = append(issues, "googleServicesFile not configured in app.json ios section")
		}
		if _, exists := config.Expo.IOS["bundleIdentifier"]; !exists {
			issues = append(issues, "iOS bundle identifier not configured in app.json")
		}
		
		// Check for iOS push notification settings
		if entitlements, exists := config.Expo.IOS["entitlements"]; exists {
			if entMap, ok := entitlements.(map[string]any); ok {
				if _, hasAPS := entMap["aps-environment"]; !hasAPS {
					issues = append(issues, "aps-environment not configured in iOS entitlements")
				}
			}
		} else {
			issues = append(issues, "iOS entitlements not configured for push notifications")
		}

		if infoPlist, exists := config.Expo.IOS["infoPlist"]; exists {
			if plistMap, ok := infoPlist.(map[string]any); ok {
				if _, hasBackground := plistMap["UIBackgroundModes"]; !hasBackground {
					issues = append(issues, "UIBackgroundModes not configured in iOS infoPlist")
				}
			}
		} else {
			issues = append(issues, "iOS infoPlist not configured for background modes")
		}
	} else {
		issues = append(issues, "iOS configuration missing in app.json")
	}

	return issues
}

// CheckClixInitialization checks if Clix initialization file exists
func CheckClixInitialization(projectRoot string) bool {
	tsPath := filepath.Join(projectRoot, "clix-config.ts")
	jsPath := filepath.Join(projectRoot, "clix-config.js")
	
	_, tsErr := os.Stat(tsPath)
	_, jsErr := os.Stat(jsPath)
	
	return tsErr == nil || jsErr == nil
}

// CheckNativeCode checks if native code has been generated
func CheckNativeCode(projectRoot string) bool {
	androidPath := filepath.Join(projectRoot, "android")
	iosPath := filepath.Join(projectRoot, "ios")
	
	androidExists := false
	iosExists := false
	
	if info, err := os.Stat(androidPath); err == nil && info.IsDir() {
		androidExists = true
	}
	
	if info, err := os.Stat(iosPath); err == nil && info.IsDir() {
		iosExists = true
	}
	
	return androidExists && iosExists
}