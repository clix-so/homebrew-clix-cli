package cmd

import (
	"fmt"
	"os"
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

	useSPM := utils.Prompt("Are you using Swift Package Manager (SPM)? (Y/n)")
	if useSPM == "" || strings.ToLower(useSPM) == "y" {
		fmt.Println("\n==================================================")
		fmt.Println("📦 Please add the Clix SDK via SPM in Xcode:")
		fmt.Println("--------------------------------------------------")
		fmt.Println("1. Open your Xcode project.")
		fmt.Println("2. Go to File > Add Packages...")
		fmt.Println("3. Enter the URL below:")
		fmt.Println("   https://github.com/clix-so/clix-ios-sdk.git")
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

	fmt.Println("\n==================================================")
	fmt.Println("📱 Integrating Clix SDK for iOS...")
	fmt.Println("==================================================\n")

	fmt.Println("1️⃣  Notification Service Extension & App Group Setup")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. In Xcode, go to File > New > Target > Notification Service Extension.")
	fmt.Println("2. Name it 'NotificationServiceExtension'.")
	fmt.Println("3. After creation, you should see a 'NotificationService.swift' file added.")
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have added the extension...")
	_, _ = fmt.Scanln()

	fmt.Println("\n2️⃣  App Groups Setup (Main App & Extension)")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. Select your main app target in Xcode.")
	fmt.Println("2. Go to the 'Signing & Capabilities' tab.")
	fmt.Println("3. Click the '+' button to add a capability.")
	fmt.Println("4. Search for and add 'App Groups'.")
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have configured App Groups for the main app...")
	_, _ = fmt.Scanln()

	fmt.Println("\n3️⃣  Push Notifications Setup (Main App & Extension)")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. Select your main app target in Xcode.")
	fmt.Println("2. Go to the 'Signing & Capabilities' tab.")
	fmt.Println("3. Click the '+' button to add a capability.")
	fmt.Println("4. Search for and add 'Push Notifications'.")
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have configured Push Notifications for the main app...")
	_, _ = fmt.Scanln()

	fmt.Println("\n4️⃣  Update NotificationService.swift")
	fmt.Println("--------------------------------------------------")
	fmt.Println("1. Select the NotificationServiceExtension target.")
	fmt.Println("2. Add the App Groups capability.")
	fmt.Printf("3. Select the same group: 'group.clix.%s'.", projectID)
	fmt.Println("4. Add the Push Notifications capability as well.")
	fmt.Println("--------------------------------------------------")
	fmt.Print("Press Enter after you have configured everything for the extension target...")
	fmt.Println("Press Enter after you have configured everything for the extension target...")
	_, _ = fmt.Scanln()

	fmt.Println("\n==================================================")
	fmt.Println("🚀 Clix SDK iOS setup instructions complete!")
	fmt.Println("==================================================")
	fmt.Println("Run 'clix-cli doctor --ios' to verify your setup.\n")

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

// Function to handle Android installation
func handleAndroidInstall() {
	fmt.Println("🤖 Installing Clix SDK for Android...")
	utils.Separatorln()

	// apiKey := utils.Prompt("Enter your Public API Key")
	// projectID := utils.Prompt("Enter your Project ID")
	apiKey := "clix_pk_N1Oc6_lQOG4-xc30_6lHFEd6GGM8Nw"
	projectID := "1b198dde-66ee-45c4-9eeb-6222129d25aa"

	android.HandleAndroidInstall(apiKey, projectID)
	fmt.Println()

	// fmt.Println("\n🔍 Running doctor to verify Clix SDK and push notification setup...")
	// android.RunDoctor("")
}
