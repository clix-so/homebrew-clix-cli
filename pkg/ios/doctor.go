package ios

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/clix-so/clix-cli/pkg/logx"
)

// import firebase-related checks from the same package
// (functions are now in firebase_checks.go)

// RunDoctor performs a comprehensive check of the iOS project setup for Clix SDK
func RunDoctor() error {
	logx.Separatorln()
	logx.Log().Title().Println("Starting Clix SDK doctor for iOS…")
	logx.Separatorln()

	fmt.Print("[1/8] Checking Xcode project directory... ")
	time.Sleep(700 * time.Millisecond)
	projectPath, projectName, err := checkXcodeProject()
	if err != nil {
		fmt.Println("❌")
		return err
	}
	fmt.Printf("Found Xcode project: %s\n", projectName)

	fmt.Print("[2/8] Checking AppDelegate.swift... ")
	time.Sleep(700 * time.Millisecond)
	appDelegatePath, err := checkAppDelegateExists(projectPath)
	if err != nil {
		fmt.Println("❌ AppDelegate.swift not found")
		return err
	}
	fmt.Println("AppDelegate.swift found.")

	fmt.Print("[3/8] Checking required imports... ")
	time.Sleep(700 * time.Millisecond)
	importErrors := checkAppDelegateImports(appDelegatePath)
	if len(importErrors) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println()
		for _, msg := range importErrors {
			fmt.Println(msg)
		}
	}

	fmt.Print("[4/8] Checking Clix.initialize call... ")
	time.Sleep(700 * time.Millisecond)
	initErrors := checkClixInitialization(appDelegatePath)
	if len(initErrors) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println()
		for _, msg := range initErrors {
			fmt.Println(msg)
		}
	}

	// Variable to store SwiftUI integration errors
	var swiftUIErrors []string
	fmt.Print("[5/8] Checking SwiftUI integration... ")
	time.Sleep(700 * time.Millisecond)
	appSwiftPath, err := findAppSwiftFile(projectPath)
	if err == nil {
		swiftUIErrors = checkSwiftUIIntegration(appSwiftPath)
		if len(swiftUIErrors) == 0 {
			fmt.Println("OK")
		} else {
			fmt.Println()
			for _, msg := range swiftUIErrors {
				fmt.Println(msg)
			}
		}
	} else {
		fmt.Println("(skipped: not a SwiftUI app)")
	}

	// Check push notification capability
	fmt.Print("[6/8] Checking push notification capability... ")
	time.Sleep(700 * time.Millisecond)
	pushCapabilities, err := checkPushCapabilities(projectPath, projectName)
	if err != nil {
		fmt.Println("❌ Push notification capability not found or not enabled.")
	} else if !pushCapabilities {
		fmt.Println("❌ 'aps-environment' not set in entitlements file.")
	} else {
		fmt.Println("Push notification capability enabled.")
	}

	// Check Firebase integration
	fmt.Print("[7/8] Checking Firebase integration... ")
	time.Sleep(700 * time.Millisecond)
	firebaseErrors := checkFirebaseIntegration(appDelegatePath)
	if len(firebaseErrors) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println()
		for _, msg := range firebaseErrors {
			fmt.Println(msg)
		}
	}

	// Check GoogleService-Info.plist
	fmt.Print("[8/8] Checking GoogleService-Info.plist... ")
	time.Sleep(700 * time.Millisecond)
	plistError := checkGoogleServicePlist(projectPath)
	if plistError != nil {
		fmt.Println(plistError)
	} else {
		fmt.Println("GoogleService-Info.plist found.")
	}

	logx.Separatorln()
	if len(importErrors) > 0 || len(initErrors) > 0 || len(swiftUIErrors) > 0 || !pushCapabilities || len(firebaseErrors) > 0 || plistError != nil {
		fmt.Println("⚠️ Some issues were found with your Clix SDK integration.")
		fmt.Println("  └ Please fix the issues mentioned above to ensure proper push notification delivery.")
		fmt.Println("  └ Run 'clix-cli install --ios' to fix most issues automatically.")
	} else {
		fmt.Println("🎉 Your iOS project is properly configured for Clix SDK!")
		fmt.Println("  └ Push notifications should be working correctly.")
	}
	logx.Separatorln()

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

	// Find any .entitlements file in the project directory
	var entitlementsPath string
	files, err := os.ReadDir(projectPath)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".entitlements") {
				entitlementsPath = filepath.Join(projectPath, file.Name())
				break
			}
		}
	}

	// Fallback to project name if no .entitlements file found
	if entitlementsPath == "" {
		entitlementsPath = filepath.Join(projectPath, fmt.Sprintf("%s.entitlements", projectName))
	}

	_, err = os.Stat(entitlementsPath)
	if err != nil {
		return false, err
	}

	content, err := os.ReadFile(entitlementsPath)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(content), "aps-environment"), nil
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
