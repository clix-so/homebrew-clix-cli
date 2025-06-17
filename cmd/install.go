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

// handleIOSInstall handles the iOS installation process by delegating to the ios package
func handleIOSInstall() {
	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")

	logx.Log().Println("🍏 Installing Clix SDK for iOS...")
	logx.Separatorln()

	ios.HandleIOSInstall(projectID, apiKey)
}

// handleAndroidInstall handles the Android installation process by delegating to the android package
func handleAndroidInstall() {
	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")

	logx.Log().Println("🤖 Installing Clix SDK for Android...")
	logx.Separatorln()

	android.HandleAndroidInstall(apiKey, projectID)
}
