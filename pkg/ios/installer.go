package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InstallClixIOS(projectID, apiKey string) error {
	// Store errors to display at the end
	var installErrors []string
	appPath, err := FindAppPath()
	if err != nil {
		return err
	}
	appPath = filepath.Join(appPath, "AppDelegate.swift")
	if _, err := os.Stat(appPath); err != nil {
		// If AppDelegate.swift not found, create one and return its result
		fmt.Println("AppDelegate.swift not found, creating one...")
		return createAppDelegate(projectID, apiKey)
	}

	content, err := os.ReadFile(appPath)
	if err != nil {
		return err
	}

	updated := string(content)

	// 1. Import
	if !strings.Contains(updated, "import Clix") {
		updated = strings.Replace(updated, "import UIKit", "import UIKit\nimport Clix", 1)
	}
	if !strings.Contains(updated, "import Firebase") {
		// Add Firebase import after last import statement
		lines := strings.Split(updated, "\n")
		insertIdx := 0
		for i, line := range lines {
			if strings.HasPrefix(line, "import ") {
				insertIdx = i + 1
			}
		}
		lines = append(lines[:insertIdx], append([]string{"import Firebase"}, lines[insertIdx:]...)...)
		updated = strings.Join(lines, "\n")
	}

	// 2. Clix.initialize
	if strings.Contains(updated, "didFinishLaunchingWithOptions") && !strings.Contains(updated, "Clix.initialize") {
		updated = strings.Replace(updated,
			"return true",
			fmt.Sprintf(`
        Clix.initialize(
			apiKey: "%s",
			projectId: "%s"
		)

        return true`, apiKey, projectID),
			1)
	}

	// 3. FirebaseApp.configure
	if strings.Contains(updated, "didFinishLaunchingWithOptions") && !strings.Contains(updated, "FirebaseApp.configure") {
		updated = strings.Replace(updated,
			"return true",
			"FirebaseApp.configure()\n\n        return true",
			1)
	}

	err = os.WriteFile(appPath, []byte(updated), 0644)
	if err != nil {
		return fmt.Errorf("failed to write AppDelegate.swift: %w", err)
	}

	fmt.Println("✅ Clix SDK successfully integrated into AppDelegate.swift")
	
	// Check if notification service extension exists, create if not
	extensionErrors := installNotificationServiceExtension(projectID)
	if len(extensionErrors) > 0 {
		installErrors = append(installErrors, extensionErrors...)
	} else {
		fmt.Println("✅ NotificationServiceExtension successfully configured")
	}
	
	// Report any errors that occurred during installation
	if len(installErrors) > 0 {
		fmt.Println("\n⚠️ Some issues occurred during installation:")
		for _, err := range installErrors {
			fmt.Println(" -", err)
		}
		fmt.Println("\nPlease address these issues manually or contact support.")
		return fmt.Errorf("installation completed with some issues")
	}
	
	return nil
}

// installNotificationServiceExtension checks if the notification service extension exists
// and creates it if it doesn't, configuring it according to requirements
func installNotificationServiceExtension(projectID string) []string {
	var errors []string

	// Find the project path
	projectPath, projectName, err := checkXcodeProject()
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to find Xcode project: %v", err))
		return errors
	}

	// Define paths
	extensionDir := filepath.Join(projectPath, "NotificationServiceExtension")
	infoPlist := filepath.Join(extensionDir, "Info.plist")
	serviceSwift := filepath.Join(extensionDir, "NotificationService.swift")
	appEntitlements := filepath.Join(projectPath, fmt.Sprintf("%s.entitlements", projectName))
	extensionEntitlements := filepath.Join(extensionDir, "NotificationServiceExtension.entitlements")

	// Create extension directory if it doesn't exist
	if _, err := os.Stat(extensionDir); os.IsNotExist(err) {
		err = os.MkdirAll(extensionDir, 0755)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to create NotificationServiceExtension directory: %v", err))
			return errors
		}
		fmt.Println("Created NotificationServiceExtension directory")
	}

	// Create or update NotificationService.swift
	serviceSwiftContent := fmt.Sprintf(`import Clix
import UserNotifications

/// NotificationService inherits all logic from ClixNotificationServiceExtension
/// No additional logic is needed unless you want to customize notification handling.
class NotificationService: ClixNotificationServiceExtension {

  // Initialize with your Clix project ID
  override init() {
    super.init()

    // Register your Clix project ID
    register(projectId: "%s")
  }

  override func didReceive(
    _ request: UNNotificationRequest,
    withContentHandler contentHandler: @escaping (UNNotificationContent) -> Void
  ) {
    // Call super to handle image downloading and send push received event
    super.didReceive(request, withContentHandler: contentHandler)
  }

  override func serviceExtensionTimeWillExpire() {
    super.serviceExtensionTimeWillExpire()
  }
}
`, projectID)

	err = os.WriteFile(serviceSwift, []byte(serviceSwiftContent), 0644)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to write NotificationService.swift: %v", err))
	} else {
		fmt.Println("Created/Updated NotificationService.swift")
	}

	// Create or update Info.plist with NSAppTransportSecurity
	infoPlistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDevelopmentRegion</key>
	<string>$(DEVELOPMENT_LANGUAGE)</string>
	<key>CFBundleDisplayName</key>
	<string>NotificationServiceExtension</string>
	<key>CFBundleExecutable</key>
	<string>$(EXECUTABLE_NAME)</string>
	<key>CFBundleIdentifier</key>
	<string>$(PRODUCT_BUNDLE_IDENTIFIER)</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>$(PRODUCT_NAME)</string>
	<key>CFBundlePackageType</key>
	<string>$(PRODUCT_BUNDLE_PACKAGE_TYPE)</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>NSExtension</key>
	<dict>
		<key>NSExtensionPointIdentifier</key>
		<string>com.apple.usernotifications.service</string>
		<key>NSExtensionPrincipalClass</key>
		<string>$(PRODUCT_MODULE_NAME).NotificationService</string>
	</dict>
	<key>NSAppTransportSecurity</key>
	<dict>
		<key>NSAllowsArbitraryLoads</key>
		<true/>
	</dict>
</dict>
</plist>
`

	err = os.WriteFile(infoPlist, []byte(infoPlistContent), 0644)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to write Info.plist: %v", err))
	} else {
		fmt.Println("Created/Updated Info.plist with NSAppTransportSecurity")
	}

	// Create or update entitlements files with app groups
	// App Group format: group.clix.{project_id}
	appGroupID := fmt.Sprintf("group.clix.%s", projectID)

	// Create app entitlements if it doesn't exist
	appEntitlementsContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>com.apple.security.application-groups</key>
	<array>
		<string>%s</string>
	</array>
	<key>aps-environment</key>
	<string>development</string>
</dict>
</plist>
`, appGroupID)

	if _, err := os.Stat(appEntitlements); os.IsNotExist(err) {
		err = os.WriteFile(appEntitlements, []byte(appEntitlementsContent), 0644)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to create app entitlements file: %v", err))
		} else {
			fmt.Println("Created app entitlements file with App Group")
		}
	} else {
		// If file exists, check if it has app groups and update if needed
		content, err := os.ReadFile(appEntitlements)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to read app entitlements file: %v", err))
		} else {
			if !strings.Contains(string(content), "com.apple.security.application-groups") {
				// Add app groups if missing
				updatedContent := strings.Replace(string(content), "<dict>", fmt.Sprintf(`<dict>
	<key>com.apple.security.application-groups</key>
	<array>
		<string>%s</string>
	</array>`, appGroupID), 1)
				err = os.WriteFile(appEntitlements, []byte(updatedContent), 0644)
				if err != nil {
					errors = append(errors, fmt.Sprintf("Failed to update app entitlements file: %v", err))
				} else {
					fmt.Println("Updated app entitlements file with App Group")
				}
			} else if !strings.Contains(string(content), appGroupID) {
				// Update app group ID if it's different
				errors = append(errors, "App entitlements file already has an App Group, but with a different ID.")
				errors = append(errors, fmt.Sprintf("Please manually update it to use: %s", appGroupID))
			}
		}
	}

	// Create extension entitlements
	extensionEntitlementsContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>com.apple.security.application-groups</key>
	<array>
		<string>%s</string>
	</array>
</dict>
</plist>
`, appGroupID)

	err = os.WriteFile(extensionEntitlements, []byte(extensionEntitlementsContent), 0644)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to create extension entitlements file: %v", err))
	} else {
		fmt.Println("Created/Updated extension entitlements file with App Group")
	}

	// Inform user about next steps in Xcode
	fmt.Println("\nℹ️ Next steps in Xcode:")
	fmt.Println("1. In Xcode, go to File > New > Target and select 'Notification Service Extension'")
	fmt.Println("2. Name it 'NotificationServiceExtension' (to match the files we created)")
	fmt.Println("3. In the extension target settings:")
	fmt.Println("   - Add the App Group capability and select '" + appGroupID + "'")
	fmt.Println("   - Enable Push Notifications capability")
	fmt.Println("4. Replace the default NotificationService.swift with our version")
	fmt.Println("   (The CLI has already created this file for you)")
	fmt.Println("5. Update the Info.plist with our version that includes NSAppTransportSecurity")

	return errors
}

func createAppDelegate(projectId, apiKey string) error {
	template := fmt.Sprintf(`import UIKit
import Clix
import Firebase

class AppDelegate: UIResponder, UIApplicationDelegate, UNUserNotificationCenterDelegate {
    func application(_ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {

        Clix.initialize(
			apiKey: "%s",
			projectId: "%s"
		)

        return true
    }
}
`, apiKey, projectId)

	appPath, err := FindAppPath()
	if err != nil {
		return err
	}
	appPath = filepath.Join(appPath, "AppDelegate.swift")

	err = os.WriteFile(appPath, []byte(template), 0644)
	if err != nil {
		return err
	}

	// Locate and modify <YourProjectName>App.swift
	projectDir := filepath.Dir(appPath)
	var appSwiftPath string
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
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
		return err
	}
	if appSwiftPath == "" {
		// Could not find App.swift with @main
		return nil
	}

	content, err := os.ReadFile(appSwiftPath)
	if err != nil {
		return err
	}
	contentStr := string(content)

	if strings.Contains(contentStr, "@UIApplicationDelegateAdaptor(AppDelegate.self)") {
		// Already contains the adaptor, no change needed
		return nil
	}

	// Find struct declaration line with 'struct ...App: App'
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "struct ") && strings.Contains(trimmed, ": App") && strings.Contains(contentStr, "@main") {
			// Insert '@UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate' just before first '{' in this line
			idx := strings.Index(line, "{")
			if idx != -1 {
				indent := ""
				for _, ch := range line[:idx] {
					if ch == ' ' || ch == '\t' {
						indent += string(ch)
					} else {
						indent = ""
					}
				}
				insertLine := indent + "    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate"
				// Insert after this line
				newLines := append(lines[:i+1], append([]string{insertLine}, lines[i+1:]...)...)
				contentStr = strings.Join(newLines, "\n")
				break
			}
		}
	}

	err = os.WriteFile(appSwiftPath, []byte(contentStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
