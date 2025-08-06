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

var doctorIosFlag bool
var doctorAndroidFlag bool
var doctorExpoFlag bool
var doctorFlutterFlag bool

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check Clix SDK integration status",
	Long: `The doctor command checks if your project has all the required 
configurations for push notifications and Clix SDK integration.
It verifies each step of the setup process and provides guidance
for any issues found.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Automatically Detect the platform
		if !doctorIosFlag && !doctorAndroidFlag && !doctorExpoFlag && !doctorFlutterFlag {
			doctorIosFlag, doctorAndroidFlag, doctorExpoFlag, doctorFlutterFlag = utils.DetectAllPlatforms()

			if !doctorIosFlag && !doctorAndroidFlag && !doctorExpoFlag && !doctorFlutterFlag {
				fmt.Fprintln(os.Stderr, "❗ Could not detect platform. Please specify --ios, --android, --expo, or --flutter")
				os.Exit(1)
			}
		}

		if doctorIosFlag {
			fmt.Println("🔍 Checking Clix SDK integration for iOS...")
			err := ios.RunDoctor()
			if err != nil {
				fmt.Fprintln(os.Stderr, "❌ Doctor check failed:", err)
				os.Exit(1)
			}
		}

		if doctorAndroidFlag {
			fmt.Println("🔍 Checking Clix SDK integration for Android...")
			logx.Separatorln()
			android.RunDoctor("") // pass project root if needed, or ""
		}

		if doctorExpoFlag {
			err := expo.RunDoctor()
			if err != nil {
				fmt.Fprintln(os.Stderr, "❌ Doctor check failed:", err)
				os.Exit(1)
			}
		}

		if doctorFlutterFlag {
			fmt.Println("🔍 Checking Clix SDK integration for Flutter...")
			err := flutter.RunDoctor()
			if err != nil {
				fmt.Fprintln(os.Stderr, "❌ Doctor check failed:", err)
				os.Exit(1)
			}
		}

		if !doctorIosFlag && !doctorAndroidFlag && !doctorExpoFlag && !doctorFlutterFlag {
			fmt.Fprintln(os.Stderr, "❗ Please specify --ios, --android, --expo, or --flutter")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVar(&doctorIosFlag, "ios", false, "Check Clix for iOS")
	doctorCmd.Flags().BoolVar(&doctorAndroidFlag, "android", false, "Check Clix for Android")
	doctorCmd.Flags().BoolVar(&doctorExpoFlag, "expo", false, "Check Clix for React Native Expo")
	doctorCmd.Flags().BoolVar(&doctorFlutterFlag, "flutter", false, "Check Clix for Flutter")
}
