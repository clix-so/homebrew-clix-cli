package android

import (
	"os"
	"path/filepath"
	"strings"
)

// app/build.gradle or app/build.gradle.kts
func GetAppBuildGradlePath(projectRoot string) string {
	gradlePaths := []string{
		filepath.Join(projectRoot, "app", "build.gradle"),
		filepath.Join(projectRoot, "app", "build.gradle.kts"),
	}

	for _, path := range gradlePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// app/src/main/java or app/src/main/kotlin
func GetBaseDirPath(projectRoot string) string {
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")

	if _, err := os.Stat(javaDir); err == nil {
		return javaDir
	} else if _, err := os.Stat(kotlinDir); err == nil {
		return kotlinDir
	} else {
		return ""
	}
}

// app/src/main/java/com/example/app
func GetSourceDirPath(projectRoot string) string {
	baseDir := GetBaseDirPath(projectRoot)
	if baseDir == "" {
		return ""
	}

	packageName := GetPackageName(projectRoot)
	if packageName == "" {
		return ""
	}

	packagePath := strings.ReplaceAll(packageName, ".", string(filepath.Separator))
	sourceDir := filepath.Join(baseDir, packagePath)

	return sourceDir
}

// app/src/main/AndroidManifest.xml
func GetAndroidManifestPath(projectRoot string) string {
	return filepath.Join(projectRoot, "app", "src", "main", "AndroidManifest.xml")
}

// gradle/libs.versions.toml
func GetVersionCatalogPath(projectRoot string) string {
	p := filepath.Join(projectRoot, "gradle", "libs.versions.toml")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return ""
}

func HasVersionCatalog(projectRoot string) bool {
	return GetVersionCatalogPath(projectRoot) != ""
}
