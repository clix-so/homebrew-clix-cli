package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UninstallClixIOS() error {
	appPath, err := FindAppPath()
	if err != nil {
		return err
	}
	appPath = filepath.Join(appPath, "AppDelegate.swift")
	if _, err := os.Stat(appPath); err != nil {
		return fmt.Errorf("failed to find AppDelegate.swift: %w", err)
	}

	content, err := os.ReadFile(appPath)
	if err != nil {
		return fmt.Errorf("failed to read AppDelegate.swift: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var cleaned []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "import Clix") ||
			strings.Contains(trimmed, "import Firebase") ||
			strings.Contains(trimmed, "FirebaseApp.configure()") ||
			strings.Contains(trimmed, "Clix.initialize") ||
			strings.Contains(trimmed, "UNUserNotificationCenter.current().delegate = self") ||
			strings.HasPrefix(trimmed, "Clix.") {
			continue
		}
		cleaned = append(cleaned, line)
	}

	err = os.WriteFile(appPath, []byte(strings.Join(cleaned, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write AppDelegate.swift: %w", err)
	}

	fmt.Println("✅ All the Clix SDK code has been removed from the project.")
	return nil
}
