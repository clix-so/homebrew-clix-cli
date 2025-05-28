package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checkFirebaseIntegration checks for 'import Firebase' and 'FirebaseApp.configure()' in AppDelegate.swift
func checkFirebaseIntegration(appDelegatePath string) []string {
	content, err := os.ReadFile(appDelegatePath)
	if err != nil {
		return []string{fmt.Sprintf("❌ Error reading AppDelegate.swift: %s", err)}
	}
	var errors []string
	contentStr := string(content)
	if !strings.Contains(contentStr, "import Firebase") {
		errors = append(errors, "❌ Missing 'import Firebase' in AppDelegate.swift")
		errors = append(errors, "  └ Add 'import Firebase' at the top of your AppDelegate.swift file")
	}
	if !strings.Contains(contentStr, "FirebaseApp.configure") {
		errors = append(errors, "❌ Missing 'FirebaseApp.configure' call in AppDelegate.swift")
		errors = append(errors, "  └ Add 'FirebaseApp.configure()' or 'FirebaseApp.configure(options:)' in your didFinishLaunchingWithOptions method")
	}
	return errors
}

// checkGoogleServicePlist checks if GoogleService-Info.plist exists in the project directory
func checkGoogleServicePlist(projectPath string) error {
	plistPath := filepath.Join(projectPath, "GoogleService-Info.plist")
	_, err := os.Stat(plistPath)
	if err != nil {
		return fmt.Errorf("❌ GoogleService-Info.plist not found in project directory.\n  └ Download it from Firebase Console and add it to your Xcode project root.")
	}
	return nil
}
