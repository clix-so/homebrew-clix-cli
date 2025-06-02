package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/android"
	"github.com/clix-so/clix-cli/pkg/ios"
	"github.com/clix-so/clix-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var iosFlag bool
var androidFlag bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Clix SDK into your project",
	Run: func(cmd *cobra.Command, args []string) {
		// Automatically Detect the platform
		if !iosFlag && !androidFlag {
			iosFlag, androidFlag = utils.DetectPlatform()

			if !iosFlag && !androidFlag {
				fmt.Fprintln(os.Stderr, "❗ Could not detect platform. Please specify --ios or --android")
				os.Exit(1)
			}
		}

		if iosFlag {
			handleIOSInstall()
		}

		if androidFlag {
			handleAndroidInstall()
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&iosFlag, "ios", false, "Install Clix for iOS")
	installCmd.Flags().BoolVar(&androidFlag, "android", false, "Install Clix for Android")
}

// Function to handle iOS installation
func handleIOSInstall() {
	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")

	// Automatically detect whether the project is using CocoaPods or SPM
	usingSPM, usingCocoaPods := detectPackageManager()

	if usingSPM {
		fmt.Println("\n==================================================")
		fmt.Println("📦 Swift Package Manager (SPM) detected!")
		fmt.Println("📦 Please add the Clix SDK via SPM in Xcode:")
		fmt.Println("--------------------------------------------------")
		fmt.Println("1. Open your Xcode project.")
		fmt.Println("2. Go to File > Add Package Dependencies")
		fmt.Println("3. Enter the URL below to the input on the right side")
		fmt.Println("   https://github.com/clix-so/clix-ios-sdk.git")
		fmt.Println("4. Select 'Up to Next Major' for the version rule")
		fmt.Println("==================================================")
		fmt.Println("Press Enter to continue...")
		_, _ = fmt.Scanln()
	} else if usingCocoaPods {
		fmt.Println("\n==================================================")
		fmt.Println("📦 CocoaPods detected!")
		fmt.Println("🤖 Installing Clix SDK for iOS via CocoaPods")
		fmt.Println("==================================================")
		err := utils.RunShellCommand("pod", "Clix")
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ Failed to run 'pod Clix':", err)
			return
		}
	} else {
		// If neither is detected, ask the user
		useSPM := utils.Prompt("Could not automatically detect package manager. Are you using Swift Package Manager (SPM)? (Y/n)")
		if useSPM == "" || strings.ToLower(useSPM) == "y" {
			fmt.Println("\n==================================================")
			fmt.Println("📦 Please add the Clix SDK via SPM in Xcode:")
			fmt.Println("--------------------------------------------------")
			fmt.Println("1. Open your Xcode project.")
			fmt.Println("2. Go to File > Add Package Dependencies")
			fmt.Println("3. Enter the URL below to the input on the right side")
			fmt.Println("   https://github.com/clix-so/clix-ios-sdk.git")
			fmt.Println("4. Select 'Up to Next Major' for the version rule")
			fmt.Println("5. Click 'Add Package' to add the Clix SDK")
			fmt.Println("6. Add your main app to the target list")
			fmt.Println("==================================================")
			fmt.Println("Press Enter to continue...")
			_, _ = fmt.Scanln()
		} else {
			fmt.Println("\n==================================================")
			fmt.Println("🤖 Installing Clix SDK for iOS via CocoaPods")
			fmt.Println("==================================================")
			err := utils.RunShellCommand("pod", "Clix")
			if err != nil {
				fmt.Fprintln(os.Stderr, "❌ Failed to run 'pod Clix':", err)
				return
			}
		}
	}

	fmt.Println("\n==================================================")
	fmt.Println("📱 Integrating Clix SDK for iOS...")
	fmt.Println("==================================================")

	fmt.Println("1️⃣  Notification Service Extension & App Group Setup")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. In Xcode, go to File > New > Target > Notification Service Extension.")
	fmt.Println("2. Name it 'NotificationServiceExtension'.")
	fmt.Println("3. After creation, you should see a 'NotificationService.swift' file added.")
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have added the extension...")
	_, _ = fmt.Scanln()

	fmt.Println("\n2️⃣  Main App Setup")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. Select your main app target in Xcode.")
	fmt.Println("2. Go to the 'Signing & Capabilities' tab.")
	fmt.Println("3. Click the '+ Capability' button to add a capability.")
	fmt.Println("4. Search for and add 'Push Notifications'.")
	fmt.Println("5. Search for and add 'App Groups'.")
	fmt.Printf("6. Add the App Group: 'group.clix.%s'.\n", projectID)
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have configured App Groups for the main app...")
	_, _ = fmt.Scanln()

	fmt.Println("\n3️⃣  NotificationServiceExtension Setup")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. Select the NotificationServiceExtension target.")
	fmt.Println("2. Go to the 'Signing & Capabilities' tab.")
	fmt.Println("3. Add the App Groups capability.")
	fmt.Printf("4. Select the same group: 'group.clix.%s'.\n", projectID)
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have configured App Groups for the extension target...")
	_, _ = fmt.Scanln()

	fmt.Println("\n4️⃣  Update NotificationServiceExtension Dependencies")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. Select the NotificationServiceExtension target.")
	fmt.Println("2. Go to the 'General' tab.")
	fmt.Println("3. Click '+' under 'Frameworks, Libraries, and Embedded Content'.")
	fmt.Println("4. Search for and add 'Clix'.")
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have configured everything for the extension target...")
	_, _ = fmt.Scanln()

	fmt.Println("\n==================================================")
	fmt.Println("🚀 Clix SDK iOS setup instructions complete!")
	fmt.Println("==================================================")
	fmt.Println("Run 'clix-cli doctor --ios' to verify your setup.")

	err := ios.InstallClixIOS(projectID, apiKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Failed:", err)
		return
	}

	extensionErrors := ios.UpdateNotificationServiceExtension(projectID)
	if len(extensionErrors) > 0 {
		fmt.Fprintln(os.Stderr, "❌ Failed to update NotificationServiceExtension:", extensionErrors)
	} else {
		fmt.Println("✅ NotificationServiceExtension successfully configured")
	}

	fmt.Println("\n🔍 Running doctor to verify Clix SDK and push notification setup...")
	doctorErr := ios.RunDoctor()
	if doctorErr != nil {
		fmt.Fprintln(os.Stderr, "❌ Doctor check failed:", doctorErr)
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

// Function to handle Android installation
func handleAndroidInstall() {
	fmt.Println("🤖 Installing Clix SDK for Android...")
	apiKey := utils.Prompt("Enter your Public API Key")
	projectID := utils.Prompt("Enter your Project ID")
	android.HandleAndroidInstall(apiKey, projectID)

	fmt.Println("\n🔍 Running doctor to verify Clix SDK and push notification setup...")
	android.RunDoctor("")
}
