package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
	"github.com/clix-so/clix-cli/pkg/utils"
)

// HandleIOSInstall guides the user through the iOS installation process
func HandleIOSInstall(projectID, apiKey string) {
	// Automatically detect whether the project is using CocoaPods or SPM
	logx.Log().Println("🔍 Detecting package manager...")
	usingSPM, usingCocoaPods := detectPackageManager()

	if usingSPM {
		logx.Log().Println("📦 Swift Package Manager (SPM) detected!")
		logx.Log().Println("📦 Please add the Clix SDK via SPM in Xcode:")
		logx.Separatorln()
		logx.Log().Println("1. Open your Xcode project.")
		logx.Log().Println("2. Go to File > Add Package Dependencies")
		logx.Log().Println("3. Enter the URL below to the input on the right side")
		logx.Log().Println("   https://github.com/clix-so/clix-ios-sdk.git")
		logx.Log().Println("4. Select 'Up to Next Major' for the version rule")
		logx.Separatorln()
		utils.Prompt("Press Enter to continue...")
	} else if usingCocoaPods {
		logx.Log().Println("📦 CocoaPods detected!")
		logx.Log().Println("🤖 Installing Clix SDK for iOS via CocoaPods")
		logx.Separatorln()
		err := utils.RunShellCommand("pod", "Clix")
		if err != nil {
			logx.Log().Failure().Println(fmt.Sprintf("Failed to run 'pod Clix': %s", err))
			return
		}
	} else {
		// If neither is detected, ask the user
		useSPM := utils.Prompt("Could not automatically detect package manager. Are you using Swift Package Manager (SPM)? (Y/n)")
		if useSPM == "" || strings.ToLower(useSPM) == "y" {
			logx.Log().Println("📦 Please add the Clix SDK via SPM in Xcode:")
			logx.Separatorln()
			logx.Log().Println("1. Open your Xcode project.")
			logx.Log().Println("2. Go to File > Add Package Dependencies")
			logx.Log().Println("3. Enter the URL below to the input on the right side")
			logx.Log().Println("   https://github.com/clix-so/clix-ios-sdk.git")
			logx.Log().Println("4. Select 'Up to Next Major' for the version rule")
			logx.Log().Println("5. Click 'Add Package' to add the Clix SDK")
			logx.Log().Println("6. Add your main app to the target list")
			logx.Separatorln()
			utils.Prompt("Press Enter to continue...")
		} else {
			logx.Log().Println("🤖 Installing Clix SDK for iOS via CocoaPods")
			logx.Separatorln()
			err := utils.RunShellCommand("pod", "Clix")
			if err != nil {
				logx.Log().Failure().Println(fmt.Sprintf("Failed to run 'pod Clix': %s", err))
				return
			}
		}
	}

	logx.NewLine()
	logx.Log().Println("📱 Integrating Clix SDK for iOS...")
	logx.Separatorln()

	logx.Log().Println("1️⃣  Notification Service Extension & App Group Setup")
	logx.Separatorln()
	logx.Log().Println("1. In Xcode, go to File > New > Target > Notification Service Extension.")
	logx.Log().Println("📝 Updating NotificationServiceExtension...")
	logx.Log().Println("2. Name it 'NotificationServiceExtension'.")
	logx.Log().Println("3. After creation, you should see a 'NotificationService.swift' file added.")
	logx.Separatorln()
	utils.Prompt("Press Enter after you have added the extension...")

	logx.NewLine()
	logx.Log().Println("2️⃣  Main App Setup")
	logx.Separatorln()
	logx.Log().Println("1. Select your main app target in Xcode.")
	logx.Log().Println("2. Go to the 'Signing & Capabilities' tab.")
	logx.Log().Println("3. Click the '+ Capability' button to add a capability.")
	logx.Log().Println("4. Search for and add 'Push Notifications'.")
	logx.Log().Println("5. Search for and add 'App Groups'.")
	logx.Log().Println(fmt.Sprintf("6. Add the App Group: 'group.clix.%s'.", projectID))
	logx.Separatorln()
	utils.Prompt("Press Enter after you have configured App Groups for the main app...")

	logx.NewLine()
	logx.Log().Println("3️⃣  NotificationServiceExtension Setup")
	logx.Separatorln()
	logx.Log().Println("1. Select the NotificationServiceExtension target.")
	logx.Log().Println("2. Go to the 'Signing & Capabilities' tab.")
	logx.Log().Println("3. Add the App Groups capability.")
	logx.Log().Println(fmt.Sprintf("4. Select the same group: 'group.clix.%s'.", projectID))
	logx.Separatorln()
	utils.Prompt("Press Enter after you have configured App Groups for the extension target...")

	logx.NewLine()
	logx.Log().Println("4️⃣  Update NotificationServiceExtension Dependencies")
	logx.Separatorln()
	logx.Log().Println("1. Select the NotificationServiceExtension target.")
	logx.Log().Println("2. Go to the 'General' tab.")
	logx.Log().Println("3. Click '+' under 'Frameworks, Libraries, and Embedded Content'.")
	logx.Log().Println("4. Search for and add 'Clix'.")
	logx.Separatorln()
	utils.Prompt("Press Enter after you have configured everything for the extension target...")

	logx.NewLine()
	logx.Log().Println("🚀 Clix SDK iOS setup instructions complete!")
	logx.Separatorln()
	logx.Log().Println("Running installation command to install Clix SDK for iOS.")

	err := InstallClixIOS(projectID, apiKey)
	if err != nil {
		logx.Log().Failure().Println(fmt.Sprintf("Failed: %s", err))
		logx.Log().Println("Please follow the manual installation guide: https://docs.clix.so/sdk-quickstart-ios")
		return
	}

	extensionErrors := UpdateNotificationServiceExtension(projectID)
	if len(extensionErrors) > 0 {
		logx.Log().Failure().Println(fmt.Sprintf("Failed to update NotificationServiceExtension: %v", extensionErrors))
		logx.Log().Println("Please follow the manual installation guide: https://docs.clix.so/sdk-quickstart-ios")
	} else {
		logx.Log().Success().Println("NotificationServiceExtension successfully configured")
	}

	logx.NewLine()
	logx.Log().Println("🔍 Running doctor check to verify Clix SDK and push notification setup...")
	doctorErr := RunDoctor()
	if doctorErr != nil {
		logx.Log().Failure().Println(fmt.Sprintf("Doctor check failed: %s", doctorErr))
	}
}

// detectPackageManager detects whether the iOS project is using CocoaPods or Swift Package Manager (SPM)
func detectPackageManager() (usingSPM bool, usingCocoaPods bool) {
	// Check for Podfile which indicates CocoaPods
	_, podfileErr := os.Stat("Podfile")
	if podfileErr == nil {
		usingCocoaPods = true
	}

	// Check for Package.swift which indicates SPM
	_, packageSwiftErr := os.Stat("Package.swift")
	if packageSwiftErr == nil {
		usingSPM = true
	}

	// Check for .xcodeproj files with SPM dependencies
	files, err := os.ReadDir(".")
	if err == nil {
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".xcodeproj") {
				// Check if project.pbxproj contains SPM references
				pbxprojPath := filepath.Join(f.Name(), "project.pbxproj")
				data, err := os.ReadFile(pbxprojPath)
				if err == nil {
					content := string(data)
					if strings.Contains(content, "XCRemoteSwiftPackageReference") {
						usingSPM = true
					}
				}
			}

			// Check for .xcworkspace which typically indicates CocoaPods
			if strings.HasSuffix(f.Name(), ".xcworkspace") && !strings.HasSuffix(f.Name(), "xcodeproj.xcworkspace") {
				usingCocoaPods = true
			}
		}
	}

	// If both are detected, prioritize the one that seems more actively used
	if usingSPM && usingCocoaPods {
		// Check if Podfile.lock exists, which indicates active use of CocoaPods
		_, podfileLockErr := os.Stat("Podfile.lock")
		if podfileLockErr == nil {
			// Podfile.lock exists, prioritize CocoaPods
			usingSPM = false
			usingCocoaPods = true
		} else {
			// No Podfile.lock, prioritize SPM
			usingSPM = true
			usingCocoaPods = false
		}
	}

	return
}

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
		logx.Log().Println("AppDelegate.swift not found, creating one...")
		return createAppDelegate(projectID, apiKey)
	}

	content, err := os.ReadFile(appPath)
	if err != nil {
		return err
	}

	updated := string(content)

	// 1. Add required imports
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

	// 2. Update class declaration to inherit from ClixAppDelegate
	if !strings.Contains(updated, "ClixAppDelegate") {
		lines := strings.Split(updated, "\n")
		for i, line := range lines {
			if strings.Contains(line, "class AppDelegate") {
				// Replace the class declaration line
				indent := ""
				for _, ch := range line {
					if ch == ' ' || ch == '\t' {
						indent += string(ch)
					} else {
						break
					}
				}
				lines[i] = indent + "class AppDelegate: ClixAppDelegate {"
				break
			}
		}
		updated = strings.Join(lines, "\n")
	}

	// 3. Update didFinishLaunchingWithOptions method to include override keyword
	if strings.Contains(updated, "didFinishLaunchingWithOptions") && !strings.Contains(updated, "override func application") {
		lines := strings.Split(updated, "\n")
		for i, line := range lines {
			if strings.Contains(line, "func application") && strings.Contains(line, "didFinishLaunchingWithOptions") {
				// Add override keyword
				indent := ""
				for _, ch := range line {
					if ch == ' ' || ch == '\t' {
						indent += string(ch)
					} else {
						break
					}
				}
				lines[i] = indent + "override " + strings.TrimSpace(line)
				break
			}
		}
		updated = strings.Join(lines, "\n")
	}

	// 4. Add FirebaseApp.configure and Clix.initialize before super.application or return true
	if strings.Contains(updated, "didFinishLaunchingWithOptions") {
		// Check if Firebase is already configured
		hasFirebaseConfig := strings.Contains(updated, "FirebaseApp.configure()")
		hasClixInit := strings.Contains(updated, "Clix.initialize")

		// If we need to add either Firebase config or Clix init
		if !hasFirebaseConfig || !hasClixInit {
			// Find the return statement
			lines := strings.Split(updated, "\n")
			returnLineIndex := -1
			returnStatement := ""

			for i, line := range lines {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "return ") {
					returnLineIndex = i
					returnStatement = trimmed
					break
				}
			}

			if returnLineIndex != -1 {
				// Get indentation from the return line
				indent := ""
				for _, ch := range lines[returnLineIndex] {
					if ch == ' ' || ch == '\t' {
						indent += string(ch)
					} else {
						break
					}
				}

				// Build the insertion content
				var insertContent strings.Builder

				// Add Firebase configuration if needed
				if !hasFirebaseConfig {
					insertContent.WriteString(indent + "FirebaseApp.configure()\n\n")
				}

				// Add Clix initialization if needed
				if !hasClixInit {
					insertContent.WriteString(fmt.Sprintf(indent+"Task {\n"+
						indent+"    await Clix.initialize(\n"+
						indent+"        config: ClixConfig(\n"+
						indent+"            projectId: \"%s\",\n"+
						indent+"            apiKey: \"%s\"\n"+
						indent+"        )\n"+
						indent+"    )\n"+
						indent+"}\n\n", projectID, apiKey))
				}

				// Replace the return line with our insertions + the original return statement
				lines[returnLineIndex] = insertContent.String() + indent + returnStatement
				updated = strings.Join(lines, "\n")
			}
		}

		// Ensure super.application is called after Firebase and Clix initialization
		if strings.Contains(updated, "return true") && !strings.Contains(updated, "return super.application") {
			updated = strings.Replace(updated,
				"return true",
				"return super.application(application, didFinishLaunchingWithOptions: launchOptions)",
				1)
		}
	}

	err = os.WriteFile(appPath, []byte(updated), 0644)
	if err != nil {
		return fmt.Errorf("failed to write AppDelegate.swift: %w", err)
	}

	logx.Log().Success().Println("Clix SDK installed successfully!")

	// Report any errors that occurred during installation
	if len(installErrors) > 0 {
		logx.Log().Println("\n⚠️ Some issues occurred during installation:")
		for _, err := range installErrors {
			fmt.Println(" -", err)
		}
		fmt.Println("\nPlease address these issues manually or contact support.")
		return fmt.Errorf("installation completed with some issues")
	}

	return nil
}

// UpdateNotificationServiceExtension: only guide and patch files (do not auto-generate)
func UpdateNotificationServiceExtension(projectID string) []string {
	var errors []string

	// Find the project path
	projectPath, err := FindAppPath()
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to find Xcode project: %v", err))
		return errors
	}
	// We don't need to extract project name here as it's not used in this function

	// Assume NotificationServiceExtension is already added in Xcode
	// Get the directory one level above the project root
	projectRoot := filepath.Dir(projectPath)                                 // First get project root
	parentDir := filepath.Dir(projectRoot)                                   // Then get one level above
	extensionDir := filepath.Join(parentDir, "NotificationServiceExtension") // One level above project root
	serviceSwift := filepath.Join(extensionDir, "NotificationService.swift")
	infoPlist := filepath.Join(extensionDir, "Info.plist")

	// Debug info for extension directory path
	fmt.Printf("Looking for extension at: %s\n", extensionDir)

	// Patch code if NotificationService.swift file exists
	if _, err := os.Stat(serviceSwift); err == nil {
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
			fmt.Println("Created or updated NotificationService.swift")
		}
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

	// Check if Info.plist exists and update NSAppTransportSecurity if needed
	if _, err := os.Stat(infoPlist); err == nil {
		content, err := os.ReadFile(infoPlist)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to read Info.plist: %v", err))
		} else {
			infoStr := string(content)
			if !strings.Contains(infoStr, "NSAppTransportSecurity") {
				insertKey := `<key>NSAppTransportSecurity</key><dict><key>NSAllowsArbitraryLoads</key><true/></dict>`
				updated := strings.Replace(infoStr, "<dict>", "<dict>\n\t"+insertKey, 1)
				err = os.WriteFile(infoPlist, []byte(updated), 0644)
				if err != nil {
					errors = append(errors, fmt.Sprintf("Failed to update Info.plist: %v", err))
				} else {
					fmt.Println("Inserted NSAppTransportSecurity into Info.plist")
				}
			}
		}
	} else {
		err = os.WriteFile(infoPlist, []byte(infoPlistContent), 0644)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to write Info.plist: %v", err))
		} else {
			fmt.Println("Created Info.plist with NSAppTransportSecurity")
		}
	}

	return errors
}

func createAppDelegate(projectId, apiKey string) error {
	template := fmt.Sprintf(`import UIKit
import Clix
import Firebase

class AppDelegate: ClixAppDelegate {
    override func application(_ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {

        FirebaseApp.configure()

        Task {
            await Clix.initialize(
                config: ClixConfig(
                    projectId: "%s",
                    apiKey: "%s"
                )
            )
        }

        return super.application(application, didFinishLaunchingWithOptions: launchOptions)
    }
}
`, projectId, apiKey)

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
