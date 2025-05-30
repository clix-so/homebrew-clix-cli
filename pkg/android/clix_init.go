package android

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AddClixInitializationToApplication inserts Clix SDK initialization code into the Application.kt if missing
func AddClixInitializationToApplication(projectRoot, apiKey, projectID string) bool {
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	found := false
	err := filepath.Walk(kotlinDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), "Application.kt") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil
			}
			content := string(data)
			if strings.Contains(content, "Clix.initialize(") {
				found = true
				return nil
			}
			// Insert imports at the top if missing
			importBlock := "import so.clix.Clix\nimport so.clix.ClixConfig\nimport so.clix.ClixLogLevel\n"
			if !strings.Contains(content, "import so.clix.Clix") {
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.HasPrefix(line, "package ") {
						// Insert after package
						lines = append(lines[:i+1], append([]string{importBlock}, lines[i+1:]...)...)
						break
					}
				}
				content = strings.Join(lines, "\n")
			}
			// Insert initialization in onCreate
			initBlock := `override fun onCreate() {
        super.onCreate()
        // Project ID: ` + projectID + `
        lifecycleScope.launch {
            try {
                val config =
            ClixConfig(
                projectId = "` + projectID + `",
                apiKey = "` + apiKey + `",
            )
        Clix.initialize(this, config)
            } catch (e: Exception) {
                // Handle initialization failure
            }
        }
    }`
			if strings.Contains(content, "override fun onCreate()") {
				// Replace existing onCreate with template
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "override fun onCreate()") {
						// Replace block (simple heuristic: next 2~20 lines)
						end := i + 1
						for ; end < len(lines) && end-i < 20; end++ {
							if strings.Contains(lines[end], "}") {
								break
							}
						}
						lines = append(lines[:i], append([]string{initBlock}, lines[end+1:]...)...)
						break
					}
				}
				content = strings.Join(lines, "\n")
			} else {
				// Insert initBlock before last '}'
				idx := strings.LastIndex(content, "}")
				if idx != -1 {
					content = content[:idx] + initBlock + "\n}" + content[idx+1:]
				}
			}
			ioutil.WriteFile(path, []byte(content), 0644)
			found = true
		}
		return nil
	})
	return found && err == nil
}
