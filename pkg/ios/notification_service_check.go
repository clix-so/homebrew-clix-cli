package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// checkNotificationServiceExtension checks if NotificationServiceExtension target and files exist
// and verifies it has the correct setup with app groups and proper implementation
func checkNotificationServiceExtension(projectPath string) []string {
	var errors []string

	// Check if the extension directory exists
	extensionDir := filepath.Join(projectPath, "NotificationServiceExtension")
	infoPlist := filepath.Join(extensionDir, "Info.plist")
	serviceSwift := filepath.Join(extensionDir, "NotificationService.swift")

	if _, err := os.Stat(extensionDir); err != nil {
		errors = append(errors, "❌ NotificationServiceExtension directory not found.")
		errors = append(errors, "  └ Please add a Notification Service Extension target in Xcode.")
		// Return early since we can't perform other checks without the extension
		return errors
	}

	// Check if Info.plist exists
	infoPlistExists := true
	if _, err := os.Stat(infoPlist); err != nil {
		errors = append(errors, "❌ NotificationServiceExtension/Info.plist not found.")
		errors = append(errors, "  └ Please ensure Info.plist exists in NotificationServiceExtension.")
		infoPlistExists = false
	}

	// Check if NotificationService.swift exists
	serviceSwiftExists := true
	if _, err := os.Stat(serviceSwift); err != nil {
		errors = append(errors, "❌ NotificationService.swift not found in NotificationServiceExtension.")
		errors = append(errors, "  └ Please ensure NotificationService.swift exists in NotificationServiceExtension.")
		serviceSwiftExists = false
	}

	// Get project name for entitlement file paths
	projectName := filepath.Base(projectPath)

	// Get project ID from AppDelegate.swift to check app group format
	projectID := extractProjectIDFromAppDelegate(filepath.Join(projectPath, "AppDelegate.swift"))

	// Check App Groups in both main app and extension targets
	appEntitlements := filepath.Join(projectPath, fmt.Sprintf("%s.entitlements", projectName))
	extensionEntitlements := filepath.Join(extensionDir, "NotificationServiceExtension.entitlements")

	// Check if both app and extension have entitlements files
	appEntitlementsExists := false
	extensionEntitlementsExists := false

	if _, err := os.Stat(appEntitlements); err == nil {
		appEntitlementsExists = true
	}

	if _, err := os.Stat(extensionEntitlements); err == nil {
		extensionEntitlementsExists = true
	}

	// Check App Group Configuration if both entitlements files exist
	if appEntitlementsExists && extensionEntitlementsExists {
		appEntitlementsContent, err := os.ReadFile(appEntitlements)
		if err == nil {
			extensionEntitlementsContent, err := os.ReadFile(extensionEntitlements)
			if err == nil {
				// Check if both entitlements have app groups
				appGroupPattern := `<key>com\.apple\.security\.application-groups</key>.*?<array>(.*?)</array>`
				appGroupRegex := regexp.MustCompile(appGroupPattern)

				appGroups := appGroupRegex.FindStringSubmatch(string(appEntitlementsContent))
				extensionGroups := appGroupRegex.FindStringSubmatch(string(extensionEntitlementsContent))

				if len(appGroups) < 2 || len(extensionGroups) < 2 {
					errors = append(errors, "❌ App Groups not properly configured in entitlements files.")
					errors = append(errors, "  └ Please ensure both main app and extension have identical app groups.")
				} else {
					// Check if they share the same app group
					if appGroups[1] != extensionGroups[1] {
						errors = append(errors, "❌ App and Extension have different App Groups.")
						errors = append(errors, "  └ The app and extension must share identical App Groups.")
					}

					// Check app group format
					expectedAppGroup := fmt.Sprintf("group.clix.%s", projectID)
					if !strings.Contains(appGroups[1], expectedAppGroup) {
						errors = append(errors, fmt.Sprintf("❌ App Group doesn't follow the required format: %s", expectedAppGroup))
						errors = append(errors, "  └ App Group should be in the format 'group.clix.{project_id}'.")
					}
				}
			}
		}
	} else {
		errors = append(errors, "❌ Missing entitlements files for app group configuration.")
		errors = append(errors, "  └ Both main app and extension need entitlements files with app groups.")
	}

	// Check Info.plist for NSAppTransportSecurity setting
	if infoPlistExists {
		infoPlistContent, err := os.ReadFile(infoPlist)
		if err == nil {
			requiredPlistContent := `<key>NSAppTransportSecurity</key>
	<dict>
		<key>NSAllowsArbitraryLoads</key>
		<true/>
	</dict>`

			if !strings.Contains(string(infoPlistContent), "NSAppTransportSecurity") || 
               !strings.Contains(string(infoPlistContent), "NSAllowsArbitraryLoads") {
				errors = append(errors, "❌ NotificationServiceExtension Info.plist missing NSAppTransportSecurity configuration.")
				errors = append(errors, "  └ Please add the following to your Info.plist:")
				errors = append(errors, fmt.Sprintf("  └ %s", requiredPlistContent))
			}
		}
	}

	// Check NotificationService.swift implementation
	if serviceSwiftExists {
		serviceContent, err := os.ReadFile(serviceSwift)
		if err == nil {
			content := string(serviceContent)

			// Check for proper imports
			if !strings.Contains(content, "import Clix") {
				errors = append(errors, "❌ Missing 'import Clix' in NotificationService.swift.")
				errors = append(errors, "  └ Please add 'import Clix' to your NotificationService.swift.")
			}

			if !strings.Contains(content, "import UserNotifications") {
				errors = append(errors, "❌ Missing 'import UserNotifications' in NotificationService.swift.")
				errors = append(errors, "  └ Please add 'import UserNotifications' to your NotificationService.swift.")
			}

			// Check for proper class inheritance
			if !strings.Contains(content, "class NotificationService: ClixNotificationServiceExtension") {
				errors = append(errors, "❌ NotificationService doesn't inherit from ClixNotificationServiceExtension.")
				errors = append(errors, "  └ Please update your NotificationService class to inherit from ClixNotificationServiceExtension.")
			}

			// Check for project ID registration
			if !strings.Contains(content, "register(projectId:") {
				errors = append(errors, "❌ Missing project ID registration in NotificationService.swift.")
				errors = append(errors, "  └ Please add 'register(projectId: \"your-project-id\")' in your init() method.")
			}

			// Check for required overrides
			if !strings.Contains(content, "override func didReceive") {
				errors = append(errors, "❌ Missing 'didReceive' method override in NotificationService.swift.")
				errors = append(errors, "  └ Please implement the 'didReceive' method that calls super.didReceive().")
			}

			if !strings.Contains(content, "override func serviceExtensionTimeWillExpire") {
				errors = append(errors, "❌ Missing 'serviceExtensionTimeWillExpire' method in NotificationService.swift.")
				errors = append(errors, "  └ Please implement the 'serviceExtensionTimeWillExpire' method that calls super.serviceExtensionTimeWillExpire().")
			}
		}
	}

	return errors
}

// extractProjectIDFromAppDelegate extracts the project ID from AppDelegate.swift
func extractProjectIDFromAppDelegate(appDelegatePath string) string {
	content, err := os.ReadFile(appDelegatePath)
	if err != nil {
		return ""
	}

	// Look for projectId in Clix.initialize
	projectIDRegex := regexp.MustCompile(`projectId:\s*"([^"]*)"`) 
	matches := projectIDRegex.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
