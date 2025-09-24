package cmd

import (
	"fmt"
	"os"

	"github.com/clix-so/clix-cli/pkg/android"
	"github.com/clix-so/clix-cli/pkg/expo"
	"github.com/clix-so/clix-cli/pkg/flutter"
	"github.com/clix-so/clix-cli/pkg/ios"
	"github.com/clix-so/clix-cli/pkg/logx"
	"github.com/clix-so/clix-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var iosFlag bool
var androidFlag bool
var verboseFlag bool
var dryRunFlag bool
var expoFlag bool
var flutterFlag bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Clix SDK into your project",
	Run: func(cmd *cobra.Command, args []string) {
		// Automatically Detect the platform
		if !iosFlag && !androidFlag && !expoFlag && !flutterFlag {
			iosFlag, androidFlag, expoFlag, flutterFlag = utils.DetectAllPlatforms()

			if !iosFlag && !androidFlag && !expoFlag && !flutterFlag {
				logx.Log().Warn().Println("Could not detect platform. Please specify --ios, --android, --expo, or --flutter")
				os.Exit(1)
			}
		}

		if iosFlag {
			handleIOSInstall()
		}

		if androidFlag {
			handleAndroidInstall()
		}

		if expoFlag {
			handleExpoInstall()
		}

		if flutterFlag {
			handleFlutterInstall()
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&iosFlag, "ios", false, "Install Clix for iOS")
	installCmd.Flags().BoolVar(&androidFlag, "android", false, "Install Clix for Android")
	installCmd.Flags().BoolVar(&verboseFlag, "verbose", false, "Show verbose output during installation")
	installCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Show what would be changed without making changes")
	installCmd.Flags().BoolVar(&expoFlag, "expo", false, "Install Clix for React Native Expo")
	installCmd.Flags().BoolVar(&flutterFlag, "flutter", false, "Install Clix for Flutter")
}

// Function to handle iOS installation
func handleIOSInstall() {
	logx.Log().Title().Println("Installing Clix SDK for iOS…")
	logx.Separatorln()
	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")

	// Display installation instructions from the ios package
	ios.DisplayIOSInstructions(projectID, apiKey, verboseFlag, dryRunFlag)

	// Install Clix iOS SDK
	err := ios.InstallClixIOS(projectID, apiKey)
	if err != nil {
		logx.Log().Failure().Println(fmt.Sprintf("Failed: %v", err))
		return
	}

	// Update NotificationServiceExtension
	extensionErrors := ios.UpdateNotificationServiceExtension(projectID)
	if len(extensionErrors) > 0 {
		logx.Log().Failure().Println("Failed to update NotificationServiceExtension:")
		for _, err := range extensionErrors {
			logx.Log().Indent(2).Println("- " + err.Error())
		}
	} else {
		logx.Log().Success().Println("NotificationServiceExtension successfully configured")
	}

	// Run doctor to verify setup
	logx.NewLine()
	logx.Log().Title().Println("Running doctor to verify Clix SDK and push notification setup…")
	doctorErr := ios.RunDoctor()
	if doctorErr != nil {
		logx.Log().Failure().Println(fmt.Sprintf("Doctor check failed: %v", doctorErr))
	}
}

// Function to handle Android installation
func handleAndroidInstall() {
	logx.Log().Title().Println("Installing Clix SDK for Android…")
	logx.Separatorln()

	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")

	logx.NewLine()
	android.HandleAndroidInstall(apiKey, projectID)
}

// Function to handle React Native Expo installation
func handleExpoInstall() {
	logx.Log().Title().Println("Installing Clix SDK for React Native Expo…")
	logx.Separatorln()

	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")

	logx.NewLine()
	expo.HandleExpoInstall(apiKey, projectID)
}

// Function to handle Flutter installation
func handleFlutterInstall() {
	logx.Log().Title().Println("Installing Clix SDK for Flutter…")
	logx.Separatorln()

	projectID := utils.Prompt("Enter your Project ID")
	apiKey := utils.Prompt("Enter your Public API Key")

	logx.NewLine()
	flutter.HandleFlutterInstall(apiKey, projectID)
}
