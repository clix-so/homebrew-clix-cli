package ios

import (
	"os"
	"path/filepath"
)

// checkNotificationServiceExtension checks if NotificationServiceExtension target and files exist
func checkNotificationServiceExtension(projectPath string) []string {
	var errors []string

	extensionDir := filepath.Join(projectPath, "NotificationServiceExtension")
	infoPlist := filepath.Join(extensionDir, "Info.plist")
	serviceSwift := filepath.Join(extensionDir, "NotificationService.swift")

	if _, err := os.Stat(extensionDir); err != nil {
		errors = append(errors, "❌ NotificationServiceExtension directory not found.")
		errors = append(errors, "  └ Please add a Notification Service Extension target in Xcode.")
	}

	if _, err := os.Stat(infoPlist); err != nil {
		errors = append(errors, "❌ NotificationServiceExtension/Info.plist not found.")
		errors = append(errors, "  └ Please ensure Info.plist exists in NotificationServiceExtension.")
	}

	if _, err := os.Stat(serviceSwift); err != nil {
		errors = append(errors, "❌ NotificationService.swift not found in NotificationServiceExtension.")
		errors = append(errors, "  └ Please ensure NotificationService.swift exists in NotificationServiceExtension.")
	}

	return errors
}
