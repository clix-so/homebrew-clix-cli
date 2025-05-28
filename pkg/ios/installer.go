package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InstallClixIOS(apiKey, projectID string) error {
	appPath, err := FindAppPath()
	if err != nil {
		return err
	}
	appPath = filepath.Join(appPath, "AppDelegate.swift")
	if _, err := os.Stat(appPath); err != nil {
		// If AppDelegate.swift not found, create one and return its result
		fmt.Println("AppDelegate.swift not found, creating one...")
		return createAppDelegate(apiKey)
	}

	content, err := os.ReadFile(appPath)
	if err != nil {
		return err
	}

	updated := string(content)

	// 1. Import
	if !strings.Contains(updated, "import Clix") {
		updated = strings.Replace(updated, "import UIKit", "import UIKit\nimport Clix", 1)
	}
	if !strings.Contains(updated, "import Firebase") {
		// Add Firebase import after last import statement
		lines := strings.Split(updated, "\n")
		insertIdx := 0
		for i, line := range lines {
			if strings.HasPrefix(line, "import ") {
				insertIdx = i + 1
			}
		}
		lines = append(lines[:insertIdx], append([]string{"import Firebase"}, lines[insertIdx:]...)...)
		updated = strings.Join(lines, "\n")
	}

	// 2. Clix.initialize
	if strings.Contains(updated, "didFinishLaunchingWithOptions") && !strings.Contains(updated, "Clix.initialize") {
		updated = strings.Replace(updated,
			"return true",
			fmt.Sprintf(`
        Clix.initialize(
			apiKey: "%s"
		)

        return true`, apiKey),
			1)
	}

	// 3. FirebaseApp.configure
	if strings.Contains(updated, "didFinishLaunchingWithOptions") && !strings.Contains(updated, "FirebaseApp.configure") {
		updated = strings.Replace(updated,
			"return true",
			"FirebaseApp.configure()\n\n        return true",
			1)
	}

	err = os.WriteFile(appPath, []byte(updated), 0644)
	if err != nil {
		return fmt.Errorf("failed to write AppDelegate.swift: %w", err)
	}

	fmt.Println("✅ Clix SDK successfully integrated into AppDelegate.swift")
	return nil
}

func createAppDelegate(apiKey string) error {
	template := fmt.Sprintf(`import UIKit
import Clix
import Firebase

class AppDelegate: UIResponder, UIApplicationDelegate, UNUserNotificationCenterDelegate {
    func application(_ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {

        Clix.initialize(
			apiKey: "%s"
		)

        return true
    }
}
`, apiKey)

	appPath, err := FindAppPath()
	if err != nil {
		return err
	}
	appPath = filepath.Join(appPath, "AppDelegate.swift")

	err = os.WriteFile(appPath, []byte(template), 0644)
	if err != nil {
		return err
	}

	// Locate and modify <YourProjectName>App.swift
	projectDir := filepath.Dir(appPath)
	var appSwiftPath string
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "App.swift") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(content), "@main") {
				appSwiftPath = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if appSwiftPath == "" {
		// Could not find App.swift with @main
		return nil
	}

	content, err := os.ReadFile(appSwiftPath)
	if err != nil {
		return err
	}
	contentStr := string(content)

	if strings.Contains(contentStr, "@UIApplicationDelegateAdaptor(AppDelegate.self)") {
		// Already contains the adaptor, no change needed
		return nil
	}

	// Find struct declaration line with 'struct ...App: App'
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "struct ") && strings.Contains(trimmed, ": App") && strings.Contains(contentStr, "@main") {
			// Insert '@UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate' just before first '{' in this line
			idx := strings.Index(line, "{")
			if idx != -1 {
				indent := ""
				for _, ch := range line[:idx] {
					if ch == ' ' || ch == '\t' {
						indent += string(ch)
					} else {
						indent = ""
					}
				}
				insertLine := indent + "    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate"
				// Insert after this line
				newLines := append(lines[:i+1], append([]string{insertLine}, lines[i+1:]...)...)
				contentStr = strings.Join(newLines, "\n")
				break
			}
		}
	}

	err = os.WriteFile(appSwiftPath, []byte(contentStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
