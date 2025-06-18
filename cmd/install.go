package cmd

import (
	"fmt"
	"os"

	"github.com/clix-so/clix-cli/pkg/android"
	"github.com/clix-so/clix-cli/pkg/ios"
	"github.com/clix-so/clix-cli/pkg/logx"
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

	// Display installation instructions from the ios package
	ios.DisplayIOSInstructions(projectID, apiKey)

	// Install Clix iOS SDK
	err := ios.InstallClixIOS(projectID, apiKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Failed:", err)
		return
	}

	// Update NotificationServiceExtension
	extensionErrors := ios.UpdateNotificationServiceExtension(projectID)
	if len(extensionErrors) > 0 {
		fmt.Fprintln(os.Stderr, "❌ Failed to update NotificationServiceExtension:", extensionErrors)
	} else {
		fmt.Println("✅ NotificationServiceExtension successfully configured")
	}

	// Run doctor to verify setup
	fmt.Println("\n🔍 Running doctor to verify Clix SDK and push notification setup...")
	doctorErr := ios.RunDoctor()
	if doctorErr != nil {
		fmt.Fprintln(os.Stderr, "❌ Doctor check failed:", doctorErr)
	}
}



// Function to handle Android installation
func handleAndroidInstall() {
	fmt.Println("🤖 Installing Clix SDK for Android...")
	logx.Separatorln()

	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")
	
	logx.NewLine()
	android.HandleAndroidInstall(apiKey, projectID)
}
