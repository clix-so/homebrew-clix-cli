package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
)

// DisplayIOSInstructions shows iOS installation instructions
func DisplayIOSInstructions(projectID string, apiKey string, verbose bool, dryRun bool) {
	// Automatically detect whether the project is using CocoaPods or SPM
	usingSPM, usingCocoaPods := detectPackageManager()

	if usingSPM {
		logx.NewLine()
		logx.Log().WithSpinner().Title().Println("Swift Package Manager (SPM) detected!")
		logx.Log().Branch().Println("📦 Please add the Clix SDK via SPM in Xcode:")
		logx.Log().Indent(2).Println("1. Open your Xcode project.")
		logx.Log().Indent(2).Println("2. Go to File > Add Package Dependencies")
		logx.Log().Indent(2).Println("3. Enter the URL below to the input on the right side")
		logx.Log().Indent(4).Code().Println("https://github.com/clix-so/clix-ios-sdk.git")
		logx.Log().Indent(2).Println("4. Select 'Up to Next Major' for the version rule")
		logx.NewLine()
		logx.Log().Println("Press Enter to continue...")
		_, _ = fmt.Scanln()
	} else if usingCocoaPods {
		logx.NewLine()
		logx.Log().WithSpinner().Title().Println("CocoaPods detected!")
		logx.Log().Branch().Println("Installing Clix SDK for iOS via CocoaPods")
		logx.NewLine()
	} else {
		// If neither is detected, ask the user
		useSPM := promptForSPM()
		if useSPM == "" || strings.ToLower(useSPM) == "y" {
			logx.NewLine()
			logx.Log().WithSpinner().Title().Println("Please add the Clix SDK via SPM in Xcode:")
			logx.Log().Indent(2).Println("1. Open your Xcode project.")
			logx.Log().Indent(2).Println("2. Go to File > Add Package Dependencies")
			logx.Log().Indent(2).Println("3. Enter the URL below to the input on the right side")
			logx.Log().Indent(4).Code().Println("https://github.com/clix-so/clix-ios-sdk.git")
			logx.Log().Indent(2).Println("4. Select 'Up to Next Major' for the version rule")
			logx.Log().Indent(2).Println("5. Click 'Add Package' to add the Clix SDK")
			logx.Log().Indent(2).Println("6. Add your main app to the target list")
			logx.NewLine()
			logx.Log().Println("Press Enter to continue...")
			_, _ = fmt.Scanln()
		} else {
			logx.NewLine()
			logx.Log().WithSpinner().Title().Println("Installing Clix SDK for iOS via CocoaPods")
			logx.NewLine()
		}
	}

	logx.NewLine()
	logx.Log().WithSpinner().Title().Println("Integrating Clix SDK for iOS…")
	logx.NewLine()

	logx.Log().Branch().Title().Println("Notification Service Extension & App Group Setup")
	logx.Log().Indent(2).Println("1. Open your Xcode project.")
	logx.Log().Indent(2).Println("2. Go to File > New > Target")
	logx.Log().Indent(2).Println("3. Select 'Notification Service Extension' and click Next.")
	logx.Log().Indent(2).Println("4. Name it 'NotificationServiceExtension' and click Finish.")
	logx.Log().Indent(2).Println("5. When prompted to activate the scheme, click 'Don't Activate'.")
	logx.Log().Indent(2).Println("6. Add Clix framework to NotificationServiceExtension target:")
	logx.Log().Indent(4).Println("   a. Select 'NotificationServiceExtension' target in the project navigator")
	logx.Log().Indent(4).Println("   b. Go to 'General' tab")
	logx.Log().Indent(4).Println("   c. Under 'Frameworks, Libraries, and Embedded Content', click '+'")
	logx.Log().Indent(4).Println("   d. Search for and add 'Clix' framework")
	logx.Log().Indent(4).Println("   e. Ensure 'Embed & Sign' is selected for the Clix framework")
	logx.NewLine()
	logx.Log().Println("Press Enter after you have created the NotificationServiceExtension…")
	_, _ = fmt.Scanln()

	logx.NewLine()
	logx.Log().WithSpinner().Title().Println("Configuring App Groups and NotificationServiceExtension")
	logx.NewLine()
	logx.Log().Branch().Println("Automating Xcode project configuration…")

	// Try to configure the Xcode project automatically
	err := ConfigureXcodeProject(projectID, verbose, dryRun)
	if err != nil {
		logx.Log().Branch().Failure().Println("Automatic configuration failed: " + err.Error())
		logx.Log().Branch().Println("Switching to manual configuration…")
		// Fall back to manual configuration
		logx.NewLine()
		logx.Log().Title().Println("App Group Configuration (Manual)")
		logx.Log().Indent(2).Println("1. Select your main app target.")
		logx.Log().Indent(2).Println("2. Go to the 'Signing & Capabilities' tab.")
		logx.Log().Indent(2).Println("3. Click '+' to add a capability.")
		logx.Log().Indent(2).Println("4. Search for and add 'App Groups'.")
		logx.Log().Indent(2).Println("5. Click '+' under App Groups to add a new group.")
		logx.Log().Indent(2).Println("6. Enter 'group.clix." + projectID + "' as the group name.")
		logx.NewLine()
		logx.Log().Println("Press Enter after you have configured App Groups for the main app…")
		_, _ = fmt.Scanln()

		logx.NewLine()
		logx.Log().Title().Println("NotificationServiceExtension Setup (Manual)")
		logx.Log().Indent(2).Println("1. Select the NotificationServiceExtension target.")
		logx.Log().Indent(2).Println("2. Go to the 'Signing & Capabilities' tab.")
		logx.Log().Indent(2).Println("3. Add the App Groups capability.")
		logx.Log().Indent(2).Println("4. Select the same group: 'group.clix." + projectID + "'.")
		logx.NewLine()
		logx.Log().Println("Press Enter after you have configured App Groups for the extension target…")
		_, _ = fmt.Scanln()

		logx.NewLine()
		logx.Log().Title().Println("4️⃣  Update NotificationServiceExtension Dependencies (Manual)")
		logx.Log().Indent(2).Println("1. Select the NotificationServiceExtension target.")
		logx.Log().Indent(2).Println("2. Go to the 'General' tab.")
		logx.Log().Indent(2).Println("3. Click '+' under 'Frameworks, Libraries, and Embedded Content'.")
		logx.Log().Indent(2).Println("4. Search for and add 'Clix' framework.")
		logx.Log().Indent(2).Println("5. Ensure 'Embed & Sign' is selected for the Clix framework.")
		logx.Log().Indent(2).Println("6. Verify that Clix appears in the frameworks list for NotificationServiceExtension.")
		logx.NewLine()
		logx.Log().Println("Press Enter after you have configured everything for the extension target…")
		_, _ = fmt.Scanln()

		// Add manual Background Modes steps
		logx.NewLine()
		logx.Log().Title().Println("Enable Background Modes on Main App (Manual)")
		logx.Log().Indent(2).Println("1. Select your MAIN app target.")
		logx.Log().Indent(2).Println("2. Open the 'Signing & Capabilities' tab.")
		logx.Log().Indent(2).Println("3. Click '+' and add 'Background Modes'.")
		logx.Log().Indent(2).Println("4. Check the boxes for:")
		logx.Log().Indent(4).Println("- Background fetch")
		logx.Log().Indent(4).Println("- Remote notifications")
		logx.NewLine()
		logx.Log().Println("Press Enter after you have enabled Background Modes…")
		_, _ = fmt.Scanln()
	} else {
		logx.Log().Branch().Success().Println("Xcode project configured successfully!")
		logx.Log().Indent(2).Println("- App Groups capability added to main app target")
		logx.Log().Indent(2).Println("- Background Modes ('Background fetch', 'Remote notifications') enabled on main app")
		logx.Log().Indent(2).Println("- App Groups capability added to NotificationServiceExtension target (if present)")
		logx.Log().Indent(2).Println("- Clix framework added to NotificationServiceExtension target with 'Embed & Sign' (if present)")
		logx.Log().Indent(2).Println("- NotificationServiceExtension is now ready to handle Clix push notifications")
		logx.NewLine()
		logx.Log().Println("Press Enter to continue...")
		_, _ = fmt.Scanln()
	}

	logx.NewLine()
	logx.Log().Branch().Success().Println("Clix SDK iOS setup instructions complete!")
	logx.NewLine()
	logx.Log().Indent(2).Code().Println("Run 'clix doctor --ios' to verify your setup.")
}

// promptForSPM asks the user if they are using Swift Package Manager
func promptForSPM() string {
	logx.Log().Println("Could not automatically detect package manager. Are you using Swift Package Manager (SPM)? (Y/n) ")
	var response string
	_, _ = fmt.Scanln(&response)
	return response
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
	installErrors := []string{} // Initialize with empty slice to avoid nil
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
			// Handle both single-line and multi-line function declarations
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "func application") {
				// Check if this line or subsequent lines contain didFinishLaunchingWithOptions
				foundMethod := strings.Contains(line, "didFinishLaunchingWithOptions")
				if !foundMethod {
					// Check next few lines for didFinishLaunchingWithOptions (multiline case)
					for j := i + 1; j < len(lines) && j < i+5; j++ {
						if strings.Contains(lines[j], "didFinishLaunchingWithOptions") {
							foundMethod = true
							break
						}
						// Stop searching if we hit another function or closing brace
						nextTrimmed := strings.TrimSpace(lines[j])
						if strings.HasPrefix(nextTrimmed, "func ") || nextTrimmed == "}" {
							break
						}
					}
				}

				if foundMethod {
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
		}
		updated = strings.Join(lines, "\n")
	}

	// 3.1. Add override keyword to other AppDelegate lifecycle methods
	appDelegateMethods := []string{
		"applicationDidBecomeActive",
		"applicationWillResignActive",
		"applicationDidEnterBackground",
		"applicationWillEnterForeground",
		"applicationWillTerminate",
		"applicationDidReceiveMemoryWarning",
	}

	for _, methodName := range appDelegateMethods {
		if strings.Contains(updated, methodName) && !strings.Contains(updated, "override func "+methodName) {
			lines := strings.Split(updated, "\n")
			for i, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if strings.HasPrefix(trimmedLine, "func "+methodName) {
					// Check if this is the method we're looking for (single-line or multiline)
					foundMethod := strings.Contains(line, methodName)
					if !foundMethod {
						// Check next few lines for method name (multiline case)
						for j := i + 1; j < len(lines) && j < i+3; j++ {
							if strings.Contains(lines[j], methodName) {
								foundMethod = true
								break
							}
							// Stop searching if we hit another function or closing brace
							nextTrimmed := strings.TrimSpace(lines[j])
							if strings.HasPrefix(nextTrimmed, "func ") || nextTrimmed == "}" {
								break
							}
						}
					}

					if foundMethod {
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
			}
			updated = strings.Join(lines, "\n")
		}
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

	logx.Log().Success().Println("Clix SDK successfully integrated into AppDelegate.swift")

	// Report any errors that occurred during installation
	if len(installErrors) > 0 {
		logx.NewLine()
		logx.Log().Warn().Println("Some issues occurred during installation:")
		// Since Go slices are never nil when initialized, we don't need this check
		// Just iterate over the slice, which will be empty if there are no errors
		for _, err := range installErrors {
			logx.Log().Indent(2).Println("- " + err)
		}
		logx.NewLine()
		logx.Log().Info().Println("Please address these issues manually or contact support.")
		return fmt.Errorf("installation completed with some issues")
	}

	return nil
}

// UpdateNotificationServiceExtension: only guide and patch files (do not auto-generate)
func UpdateNotificationServiceExtension(projectID string) []error {
	errors := []error{} // Initialize with empty slice to avoid niling

	// Find the project path
	projectPath, err := FindAppPath()
	if err != nil {
		errors = append(errors, fmt.Errorf("Failed to find Xcode project: %w", err))
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
			errors = append(errors, fmt.Errorf("Failed to write NotificationService.swift: %w", err))
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
			errors = append(errors, fmt.Errorf("Failed to read Info.plist: %w", err))
		} else {
			infoStr := string(content)
			if !strings.Contains(infoStr, "NSAppTransportSecurity") {
				insertKey := `<key>NSAppTransportSecurity</key><dict><key>NSAllowsArbitraryLoads</key><true/></dict>`
				updated := strings.Replace(infoStr, "<dict>", "<dict>\n\t"+insertKey, 1)
				err = os.WriteFile(infoPlist, []byte(updated), 0644)
				if err != nil {
					errors = append(errors, fmt.Errorf("Failed to update Info.plist: %w", err))
				} else {
					fmt.Println("Inserted NSAppTransportSecurity into Info.plist")
				}
			}
		}
	} else {
		err = os.WriteFile(infoPlist, []byte(infoPlistContent), 0644)
		if err != nil {
			errors = append(errors, fmt.Errorf("Failed to write Info.plist: %w", err))
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
