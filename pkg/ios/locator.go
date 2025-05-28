package ios

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func FindAppPath() (string, error) {
	// Check if there is a .xcodeproj folder in the current directory
	entries, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	var projectName string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcodeproj") {
			projectName = strings.TrimSuffix(entry.Name(), ".xcodeproj")
			break
		}
	}

	if projectName == "" {
		return "", errors.New("❌ No .xcodeproj found. Please run this command from the root of your Xcode project")
	}

	return filepath.Join(".", projectName), nil
}
