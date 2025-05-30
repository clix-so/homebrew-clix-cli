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

	// Check if the extension directory exists one level above the project root
	// Get the project root directory (parent of the project path)
	projectRoot := filepath.Dir(projectPath)                                 // First get project root
	parentDir := filepath.Dir(projectRoot)                                   // Then get one level above
	extensionDir := filepath.Join(parentDir, "NotificationServiceExtension") // One level above project root
	infoPlist := filepath.Join(extensionDir, "Info.plist")
	serviceSwift := filepath.Join(extensionDir, "NotificationService.swift")

	if _, err := os.Stat(extensionDir); err != nil {
		errors = append(errors, "❌ NotificationServiceExtension directory not found.")
		errors = append(errors, "  └ Please add a Notification Service Extension target in Xcode.")
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
	// Find any .entitlements file in the project directory
	var appEntitlements string
	files, err := os.ReadDir(projectPath)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".entitlements") {
				appEntitlements = filepath.Join(projectPath, file.Name())
				break
			}
		}
	}
	if appEntitlements == "" {
		// Fallback to project name if no .entitlements file found
		appEntitlements = filepath.Join(projectPath, fmt.Sprintf("%s.entitlements", projectName))
	}
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
				// Support multiline <array> blocks and extract <string> values
				// More flexible regex pattern that can handle whitespace and newlines
				appGroupPattern := `(?s)<key>\s*com\.apple\.security\.application-groups\s*</key>\s*<array>(.*?)</array>`
				appGroupRegex := regexp.MustCompile(appGroupPattern)

				// Alternative pattern if the first one doesn't match
				altAppGroupPattern := `(?s)com\.apple\.security\.application-groups\s*=\s*\(([^\)]*)\)`
				altAppGroupRegex := regexp.MustCompile(altAppGroupPattern)

				// Try primary pattern first
				appGroupsArray := appGroupRegex.FindStringSubmatch(string(appEntitlementsContent))
				extensionGroupsArray := appGroupRegex.FindStringSubmatch(string(extensionEntitlementsContent))

				// If primary pattern didn't match for app, try alternative pattern
				if len(appGroupsArray) < 2 {
					appGroupsArray = altAppGroupRegex.FindStringSubmatch(string(appEntitlementsContent))
				}

				// If primary pattern didn't match for extension, try alternative pattern
				if len(extensionGroupsArray) < 2 {
					extensionGroupsArray = altAppGroupRegex.FindStringSubmatch(string(extensionEntitlementsContent))
				}

				var appGroupsFlat, extensionGroupsFlat []string

				// Function to extract group identifiers from content
				extractGroups := func(content string) []string {
					var groups []string

					// Try XML format first: <string>group.name</string>
					stringPattern := `<string>(.*?)</string>`
					stringRegex := regexp.MustCompile(stringPattern)
					matches := stringRegex.FindAllStringSubmatch(content, -1)
					for _, m := range matches {
						if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
							groups = append(groups, strings.TrimSpace(m[1]))
						}
					}

					// If no XML format found, try alternative format: "group.name"
					if len(groups) == 0 {
						altStringPattern := `"(group\.[^"]+)"`
						altStringRegex := regexp.MustCompile(altStringPattern)
						altMatches := altStringRegex.FindAllStringSubmatch(content, -1)
						for _, m := range altMatches {
							if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
								groups = append(groups, strings.TrimSpace(m[1]))
							}
						}
					}

					return groups
				}

				// Extract groups from both app and extension
				if len(appGroupsArray) >= 2 {
					appGroupsFlat = extractGroups(appGroupsArray[1])
				}
				if len(extensionGroupsArray) >= 2 {
					extensionGroupsFlat = extractGroups(extensionGroupsArray[1])
				}

				if len(appGroupsFlat) == 0 || len(extensionGroupsFlat) == 0 {
					errors = append(errors, "❌ App Groups not properly configured in entitlements files.")
					errors = append(errors, "  └ Please ensure both main app and extension have identical app groups.")
				} else {
					// Check if they share the same app group (intersection)
					shared := false
					for _, ag := range appGroupsFlat {
						for _, eg := range extensionGroupsFlat {
							if ag == eg {
								shared = true
								break
							}
						}
						if shared {
							break
						}
					}
					if !shared {
						errors = append(errors, "❌ App and Extension have different App Groups.")
						errors = append(errors, "  └ The app and extension must share at least one identical App Group.")
					}

					// Check app group format in main app
					expectedAppGroup := fmt.Sprintf("group.clix.%s", projectID)
					foundFormat := false
					for _, ag := range appGroupsFlat {
						if ag == expectedAppGroup {
							foundFormat = true
							break
						}
					}
					if !foundFormat {
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
