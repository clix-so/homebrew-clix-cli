package ios

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

)

// import firebase-related checks from the same package
// (functions are now in firebase_checks.go)


// RunDoctor performs a comprehensive check of the iOS project setup for Clix SDK
func RunDoctor() error {
	fmt.Println("🔍 Starting Clix SDK doctor for iOS...")

	// 1. Check if we're in an Xcode project directory
	projectPath, projectName, err := checkXcodeProject()
	if err != nil {
		return err
	}
	fmt.Printf("✅ Found Xcode project: %s\n", projectName)

	// 2. Check AppDelegate.swift existence
	appDelegatePath, err := checkAppDelegateExists(projectPath)
	if err != nil {
		fmt.Println("❌ AppDelegate.swift not found")
		fmt.Println("  └ You need to create AppDelegate.swift and set it up for Clix SDK integration.")
		fmt.Println("  └ Run 'clix-cli install --ios' to automatically create the required files.")
		return nil
	}
	fmt.Println("✅ Found AppDelegate.swift")

	// 3. Check imports in AppDelegate.swift
	importErrors := checkAppDelegateImports(appDelegatePath)
	if len(importErrors) > 0 {
		for _, errMsg := range importErrors {
			fmt.Println(errMsg)
		}
	} else {
		fmt.Println("✅ All required imports are present in AppDelegate.swift")
	}

	// 4. Check Clix.initialize call
	initErrors := checkClixInitialization(appDelegatePath)
	if len(initErrors) > 0 {
		for _, errMsg := range initErrors {
			fmt.Println(errMsg)
		}
	} else {
		fmt.Println("✅ Clix.initialize is properly implemented")
	}

	// 5. Check Firebase integration in AppDelegate.swift
	firebaseErrors := checkFirebaseIntegration(appDelegatePath)
	if len(firebaseErrors) > 0 {
		for _, errMsg := range firebaseErrors {
			fmt.Println(errMsg)
		}
	} else {
		fmt.Println("✅ Firebase is properly imported and configured in AppDelegate.swift")
	}

	// 6. Check GoogleService-Info.plist existence
	plistError := checkGoogleServicePlist(projectPath)
	if plistError != nil {
		fmt.Println(plistError)
	} else {
		fmt.Println("✅ GoogleService-Info.plist exists in the project directory")
	}

	// 7. Check NotificationServiceExtension
	nseErrors := checkNotificationServiceExtension(projectPath)
	if len(nseErrors) > 0 {
		for _, errMsg := range nseErrors {
			fmt.Println(errMsg)
		}
	} else {
		fmt.Println("✅ NotificationServiceExtension is present and correctly structured")
	}

	// 8. Check push notification capabilities
	pushCapabilities, err := checkPushCapabilities(projectPath, projectName)
	if err != nil {
		fmt.Printf("❌ Error checking push notification capabilities: %s\n", err)
		fmt.Println("  └ Please open your project in Xcode and enable push notifications:")
		fmt.Println("  └ 1. Select your target under 'Targets'")
		fmt.Println("  └ 2. Go to 'Signing & Capabilities'")
		fmt.Println("  └ 3. Click '+' and add 'Push Notifications'")
	} else if !pushCapabilities {
		fmt.Println("❌ Push Notifications capability is not enabled")
		fmt.Println("  └ Please open your project in Xcode and enable push notifications:")
		fmt.Println("  └ 1. Select your target under 'Targets'")
		fmt.Println("  └ 2. Go to 'Signing & Capabilities'")
		fmt.Println("  └ 3. Click '+' and add 'Push Notifications'")
	} else {
		fmt.Println("✅ Push Notifications capability is enabled")
	}

	// 6. Check UNUserNotificationCenterDelegate implementation
	delegateErrors := checkNotificationDelegates(appDelegatePath)
	if len(delegateErrors) > 0 {
		for _, errMsg := range delegateErrors {
			fmt.Println(errMsg)
		}
	} else {
		fmt.Println("✅ Notification delegates are properly implemented")
	}

	// 7. Check UIApplicationDelegateAdaptor in SwiftUI app (if applicable)
	appSwiftPath, err := findAppSwiftFile(projectPath)
	var swiftUIErrors []string
	if err == nil {
		swiftUIErrors = checkSwiftUIIntegration(appSwiftPath)
		if len(swiftUIErrors) > 0 {
			for _, errMsg := range swiftUIErrors {
				fmt.Println(errMsg)
			}
		} else {
			fmt.Println("✅ SwiftUI app delegate adaptor is properly set up")
		}
	}

	fmt.Println("\n📋 Summary:")
	if len(importErrors) > 0 || len(initErrors) > 0 || len(delegateErrors) > 0 || (err == nil && len(swiftUIErrors) > 0) || !pushCapabilities || len(firebaseErrors) > 0 || plistError != nil {
		fmt.Println("⚠️ Some issues were found with your Clix SDK integration.")
		fmt.Println("  └ Please fix the issues mentioned above to ensure proper push notification delivery.")
		fmt.Println("  └ Run 'clix-cli install --ios' to fix most issues automatically.")
	} else {
		fmt.Println("🎉 Your iOS project is properly configured for Clix SDK!")
		fmt.Println("  └ Push notifications should be working correctly.")
	}

	return nil
}

// checkXcodeProject checks if we're in an Xcode project directory
func checkXcodeProject() (string, string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return "", "", err
	}

	var projectName string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcodeproj") {
			projectName = strings.TrimSuffix(entry.Name(), ".xcodeproj")
			return filepath.Join(".", projectName), projectName, nil
		}
	}

	return "", "", errors.New("❌ No .xcodeproj found. Please run this command from the root of your Xcode project")
}

// checkAppDelegateExists checks if AppDelegate.swift exists
func checkAppDelegateExists(projectPath string) (string, error) {
	appDelegatePath := filepath.Join(projectPath, "AppDelegate.swift")
	_, err := os.Stat(appDelegatePath)
	if err != nil {
		return "", err
	}
	return appDelegatePath, nil
}

// checkAppDelegateImports checks required imports in AppDelegate.swift
func checkAppDelegateImports(appDelegatePath string) []string {
	content, err := os.ReadFile(appDelegatePath)
	if err != nil {
		return []string{fmt.Sprintf("❌ Error reading AppDelegate.swift: %s", err)}
	}

	var errors []string
	contentStr := string(content)

	if !strings.Contains(contentStr, "import Clix") {
		errors = append(errors, "❌ Missing 'import Clix' in AppDelegate.swift")
		errors = append(errors, "  └ Add 'import Clix' at the top of your AppDelegate.swift file")
	}

	return errors
}

// checkClixInitialization checks for Clix.initialize call
func checkClixInitialization(appDelegatePath string) []string {
	content, err := os.ReadFile(appDelegatePath)
	if err != nil {
		return []string{fmt.Sprintf("❌ Error reading AppDelegate.swift: %s", err)}
	}

	var errors []string
	contentStr := string(content)

	if !strings.Contains(contentStr, "Clix.initialize") {
		errors = append(errors, "❌ Missing 'Clix.initialize' call in AppDelegate.swift")
		errors = append(errors, "  └ Add the following code in your didFinishLaunchingWithOptions method:")
		errors = append(errors, "  └ Clix.initialize(projectId: \"YOUR_PROJECT_ID\", username: \"YOUR_USERNAME\", password: \"YOUR_PASSWORD\")")
	}

	return errors
}

// checkPushCapabilities checks if push notification capability is enabled
func checkPushCapabilities(projectPath, projectName string) (bool, error) {
	// This is a simplified check. In a real implementation, you would parse the Xcode project file
	// to check if push notifications are enabled in the capabilities.
	entitlementsPath := filepath.Join(projectPath, fmt.Sprintf("%s.entitlements", projectName))
	_, err := os.Stat(entitlementsPath)
	if err != nil {
		return false, err
	}

	content, err := os.ReadFile(entitlementsPath)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(content), "aps-environment"), nil
}

// checkNotificationDelegates checks for proper notification delegate implementation
func checkNotificationDelegates(appDelegatePath string) []string {
	content, err := os.ReadFile(appDelegatePath)
	if err != nil {
		return []string{fmt.Sprintf("❌ Error reading AppDelegate.swift: %s", err)}
	}

	var errors []string
	contentStr := string(content)

	if !strings.Contains(contentStr, "UNUserNotificationCenterDelegate") {
		errors = append(errors, "❌ AppDelegate doesn't conform to UNUserNotificationCenterDelegate")
		errors = append(errors, "  └ Add 'UNUserNotificationCenterDelegate' to your AppDelegate class declaration:")
		errors = append(errors, "  └ class AppDelegate: UIResponder, UIApplicationDelegate, UNUserNotificationCenterDelegate {")
	}

	if !strings.Contains(contentStr, "userNotificationCenter") &&
		!strings.Contains(contentStr, "didReceiveNotificationResponse") {
		errors = append(errors, "❌ Missing notification handling methods")
		errors = append(errors, "  └ Add notification handling methods to your AppDelegate:")
		errors = append(errors, "  └ func userNotificationCenter(_ center: UNUserNotificationCenter, willPresent notification: UNNotification, withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void) {")
		errors = append(errors, "  └     completionHandler([.alert, .badge, .sound])")
		errors = append(errors, "  └ }")
	}

	return errors
}

// findAppSwiftFile finds the main App.swift file for SwiftUI apps
func findAppSwiftFile(projectPath string) (string, error) {
	var appSwiftPath string
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "App.swift") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(content), "@main") {
				appSwiftPath = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if appSwiftPath == "" {
		return "", errors.New("No SwiftUI App.swift file found")
	}
	return appSwiftPath, nil
}

// checkSwiftUIIntegration checks for UIApplicationDelegateAdaptor in SwiftUI apps
func checkSwiftUIIntegration(appSwiftPath string) []string {
	content, err := os.ReadFile(appSwiftPath)
	if err != nil {
		return []string{fmt.Sprintf("❌ Error reading App.swift: %s", err)}
	}

	var errors []string
	contentStr := string(content)

	if !strings.Contains(contentStr, "@UIApplicationDelegateAdaptor") {
		errors = append(errors, "❌ Missing @UIApplicationDelegateAdaptor in SwiftUI App")
		errors = append(errors, "  └ Add '@UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate' to your App struct")
	}

	return errors
}
