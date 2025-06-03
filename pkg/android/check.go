package android

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/utils"
)

// CheckGradleRepository checks if mavenCentral() is present in settings.gradle(.kts) or build.gradle(.kts)
func CheckGradleRepository(projectRoot string) bool {
	gradleFiles := []string{
		filepath.Join(projectRoot, "settings.gradle"),
		filepath.Join(projectRoot, "settings.gradle.kts"),
		filepath.Join(projectRoot, "build.gradle"),
		filepath.Join(projectRoot, "build.gradle.kts"),
	}

	found := false
	for _, file := range gradleFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		content := string(data)
		if Contains(content, "repositories") && Contains(content, "mavenCentral()") {
			found = true
			break
		}
	}

	if found {
		utils.Successln("Gradle repositories are properly configured.")
		return true
	}

	utils.Failureln("Gradle repository settings are missing.")
	return false
}

// CheckGradleDependency checks if so.clix:clix-android-sdk is present in app/build.gradle(.kts)
func CheckGradleDependency(projectRoot string) bool {
	appBuildGradleFilePath := GetAppBuildGradlePath(projectRoot)

	if appBuildGradleFilePath == "" {
		utils.Failureln("app/build.gradle(.kts) not found.")
		return false
	}

	found := false
	data, err := ioutil.ReadFile(appBuildGradleFilePath)
	if err != nil {
		utils.Failureln("Failed to read app/build.gradle(.kts)")
		return false
	}

	content := string(data)
	if Contains(content, "implementation(\"so.clix:clix-android-sdk:") {
		found = true
	} else if Contains(content, "implementation(libs.clix.android.sdk)") {
		found = true
	}
	

	if (found) {
		utils.Successln("Clix SDK dependency found.")
		return true
	}

	utils.Failureln("Clix SDK dependency is missing.")
	return false
}

// CheckGradlePlugin checks if com.google.gms:google-services is present in app/build.gradle(.kts)
func CheckGradlePlugin(projectRoot string) bool {
	gradleFiles := []string{
		filepath.Join(projectRoot, "app", "build.gradle"),
		filepath.Join(projectRoot, "app", "build.gradle.kts"),
	}

	found := false
	for _, file := range gradleFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		content := string(data)
		if Contains(content, "alias(libs.plugins.gms") {
			found = true
			break
		}
		if Contains(content, "id(\"com.google.gms.google-services\")") {
			found = true
			break
		}
	}

	if (found) {
		utils.Successln("Google services plugin found in app/build.gradle(.kts).")
		return true
	}

	fmt.Println("❌ Google services plugin not found in app/build.gradle(.kts)")
	return false
}

// CheckClixCoreImport checks if any Application class imports so.clix.core.Clix (Java or Kotlin).
func CheckClixCoreImport(projectRoot string) (bool, string) {
	manifestPath := filepath.Join(projectRoot, "app", "src", "main", "AndroidManifest.xml")
	appName, err := extractApplicationClassName(manifestPath)

	if err != nil {
		utils.Failureln("Failed to read AndroidManifest.xml")
		return false, "unknown"
	}

	if appName == "" {
		utils.Failureln("No Application class found in AndroidManifest.xml")
		return false, "missing-application"
	}

	appPath := strings.TrimPrefix(appName, ".")
	appPath = strings.ReplaceAll(appPath, ".", string(filepath.Separator))

	sourceDir := GetSourceDirPath(projectRoot)
	if sourceDir == "" {
		utils.Failureln("Source directory not found.")
		return false, "unknown"
	}

	ktPath := filepath.Join(sourceDir, appPath + ".kt")
	javaPath := filepath.Join(sourceDir, appPath + ".java")

	if _, err := os.Stat(javaPath); err == nil {
		appPath = javaPath
	} else if _, err := os.Stat(ktPath); err == nil {
		appPath = ktPath
	} else {
		utils.Failureln("Application class not found in expected locations.")
		return false, "unknown"
	}

	importFound := false
	initializeFound := false

	data, err := os.ReadFile(appPath)
	if err != nil {
		utils.Failureln("Failed to read Application class file")
		return false, "unknown"
	}
	content := string(data)
	if stringContainsImportClix(content) {
		importFound = true
	}
	if StringContainsClixInitializeInOnCreate(content) {
		initializeFound = true
	}

	if importFound {
		utils.Successln("so.clix.core.Clix is imported in Application class.")
	} else {
		utils.Failureln("so.clix.core.Clix is not imported in any Application class.")
	}

	if initializeFound {
		utils.Successln("Clix.initialize(this, ...) is called in onCreate() of Application class.")
	} else {
		utils.Failureln("Clix.initialize(this, ...) is NOT called in onCreate() of any Application class.")
	}

	if !importFound || !initializeFound {
		return false, "missing-content"
	}

	return true, ""
}

// CheckAndroidMainActivityPermissions checks MainActivity for permission request code, prints instructions if missing
func CheckAndroidMainActivityPermissions(projectRoot string) bool {
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	mainActivityFiles := []string{}

	findMainActivity := func(root string) {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (info.Name() == "MainActivity.java" || info.Name() == "MainActivity.kt") {
				mainActivityFiles = append(mainActivityFiles, path)
			}
			return nil
		})
	}

	findMainActivity(javaDir)
	findMainActivity(kotlinDir)

	if len(mainActivityFiles) == 0 {
		utils.Warnln("No MainActivity.java or MainActivity.kt found. Please ensure you have a MainActivity.") // TODO: add following action
		return false
	}

	permissionPattern := []string{
		"requestPermissions(",
		"ActivityCompat.requestPermissions(",
		"ContextCompat.checkSelfPermission(",
		"Manifest.permission.",
	}

	found := false
	for _, file := range mainActivityFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(data)
		for _, pat := range permissionPattern {
			if Contains(content, pat) {
				found = true
				break
			}
		}
	}

	if found {
		utils.Successln("MainActivity contains code requesting permissions.")
		return true
	}

	utils.Failureln("MainActivity does not contain code requesting permissions.")
	return false
}

// CheckGoogleServicesJSON checks if google-services.json exists in the correct location
func CheckGoogleServicesJSON(projectRoot string) bool {
	gsPath := filepath.Join(projectRoot, "app", "google-services.json")
	if _, err := os.Stat(gsPath); os.IsNotExist(err) {
		utils.Failureln("Missing google-services.json at app/google-services.json")
		utils.Indentln("See https://docs.clix.so/firebase-setting for setup instructions.", 3)
		return false
	}

	utils.Successln("google-services.json found")
	return true
}