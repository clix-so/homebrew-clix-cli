package cmd

import (
	"fmt"
	"os"

	"github.com/clix-so/clix-cli/pkg/ios"
	"github.com/clix-so/clix-cli/pkg/android"
	"github.com/clix-so/clix-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	uninstallIOS     bool
	uninstallAndroid bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall clix from devices",
	Run: func(cmd *cobra.Command, args []string) {
		// Automatically Detect the platform
		if !uninstallIOS && !uninstallAndroid {
			uninstallIOS, uninstallAndroid, _ = utils.DetectPlatform()

			if !uninstallIOS && !uninstallAndroid {
				fmt.Fprintln(os.Stderr, "❗ Could not detect platform. Please specify --ios or --android")
				os.Exit(1)
			}
		}

		if uninstallIOS {
			ios.UninstallClixIOS()
		}
		if uninstallAndroid {
			err := android.UninstallClixAndroid()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to uninstall Clix SDK from Android: %v\n", err)
			} else {
				fmt.Println("Clix SDK uninstalled from Android project.")
			}
		}
		if !uninstallIOS && !uninstallAndroid {
			fmt.Println("Please specify at least one platform with --ios or --android.")
		}
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallIOS, "ios", false, "Uninstall clix from iOS device")
	uninstallCmd.Flags().BoolVar(&uninstallAndroid, "android", false, "Uninstall clix from Android device")
	rootCmd.AddCommand(uninstallCmd)
}
