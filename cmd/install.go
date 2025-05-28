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
	useSPM := utils.Prompt("Are you using Swift Package Manager (SPM)? (Y/n)")
	if useSPM == "" || strings.ToLower(useSPM) == "y" {
		fmt.Println("📦 Please add the Notifly SDK via SPM in Xcode by following these steps:")
		fmt.Println("1. Open your Xcode project.")
		fmt.Println("2. Go to File > Add Packages...")
		fmt.Println("3. Enter the URL below:")
		fmt.Println("   https://github.com/clix-so/clix-ios-sdk.git")
		fmt.Println("Press Enter to continue...")
		_, _ = fmt.Scanln()
	} else {
		fmt.Println("🤖 Installing Clix SDK for iOS")

		err := utils.RunShellCommand("pod", "Clix")
		if err != nil {
			fmt.Fprintln(os.Stderr, "❌ Failed to run 'pod Clix':", err)
			return
		}
	}

	fmt.Println("📱 Integrating Clix SDK for iOS...")
	apiKey := utils.Prompt("Enter your Public API Key")
	projectID := utils.Prompt("Enter your Project ID")

	err := ios.InstallClixIOS(apiKey, projectID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Failed:", err)
		return
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
	apiKey := utils.Prompt("Enter your Public API Key")
	projectID := utils.Prompt("Enter your Project ID")
	android.HandleAndroidInstall(apiKey, projectID)

	fmt.Println("\n🔍 Running doctor to verify Clix SDK and push notification setup...")
	android.RunDoctor("")
}
