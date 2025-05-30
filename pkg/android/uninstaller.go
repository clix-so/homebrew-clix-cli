package android

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UninstallClixAndroid removes Clix SDK initialization from Application.kt and related files
func UninstallClixAndroid() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	removed := false
	err = filepath.Walk(kotlinDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), "Application.kt") {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			lines := strings.Split(string(data), "\n")
			var cleaned []string
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if strings.Contains(trimmed, "Clix.initialize") ||
					strings.Contains(trimmed, "ClixConfig(") ||
					strings.Contains(trimmed, "import so.clix.Clix") ||
					strings.Contains(trimmed, "import so.clix.ClixConfig") ||
					strings.Contains(trimmed, "import so.clix.ClixLogLevel") {
					continue
				}
				cleaned = append(cleaned, line)
			}
			err = os.WriteFile(path, []byte(strings.Join(cleaned, "\n")), 0644)
			if err == nil {
				removed = true
			}
		}
		return nil
	})
	if removed {
		fmt.Println("✅ All the Clix SDK code has been removed from the Android Application class.")
		return nil
	}
	return fmt.Errorf("No Application.kt with Clix SDK code found.")
}
