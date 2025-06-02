package android

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clix-so/clix-cli/pkg/utils"
)

// HandleAndroidInstall guides the user through the Android installation checklist.
func HandleAndroidInstall(apiKey, projectID string) {
	projectRoot, err := os.Getwd()
	if err != nil {
		utils.Failureln("Could not determine working directory.") // TODO: revisit this
		return
	}

	utils.TitlelnWithSpinner("Checking google-services.json...")
	if !CheckGoogleServicesJSON(projectRoot) {
		return
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking Gradle repository settings...")
	repoOK := CheckGradleRepository(projectRoot)
	if !repoOK {
		if AddGradleRepository(projectRoot) {
			utils.BranchSuccessln("Fixed: Automatically added")
		} else {
			utils.BranchFailureln("Could not fix automatically. Please add the following manually to settings.gradle(.kts) or build.gradle(.kts):")
		}
		utils.Grayln(`      repositories {
          mavenCentral()
      }`)
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking for Clix SDK dependency...")
	depOK := CheckGradleDependency(projectRoot)
	if !depOK {
		if AddGradleDependency(projectRoot) {
			utils.BranchSuccessln("Fixed: Automatically added")
			depOK = true
		} else {
			utils.BranchFailureln("Could not fix automatically. Please add the following manually to app/build.gradle(.kts):")
		}
		utils.Grayln(`      dependencies {
          implementation("so.clix:clix-android-sdk:1.0.0")
      }`)
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking for Google Services plugin...")
	pluginOK := CheckGradlePlugin(projectRoot)
	if !pluginOK {
		if AddGradlePlugin(projectRoot) {
			utils.BranchSuccessln("Fixed: Automatically added")
			pluginOK = true
		} else {
			utils.BranchFailureln("Could not fix automatically. Please add the following manually to build.gradle(.kts):")
		}
		utils.Grayln(`      plugins {
          id("com.google.gms.google-services") version "4.4.2"
      }`)
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking Clix SDK initialization...")
	appOK := CheckClixCoreImport(projectRoot)
	if !appOK {
		if AddClixInitializationToApplication(projectRoot, apiKey, projectID) {
			utils.BranchSuccessln("Fixed: Automatically fixed")
			appOK = true
		} else {
			utils.BranchFailureln("Could not fix automatically. Please add the following to your Application(.kt or .java):")
			// TODO: add example code snippet
		}
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking permission request...")
	mainActivityOK := CheckAndroidMainActivityPermissions(projectRoot)
	if !mainActivityOK {
		fmt.Println(`ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.POST_NOTIFICATIONS), 1001)`)
	}
	fmt.Println()

	if repoOK && depOK && appOK && mainActivityOK {
		utils.Successln("Clix SDK installation checklist complete! Your Android project is ready.")
	} else {
		utils.Failureln("Please address the above issues and re-run 'clix install --android' or 'clix doctor --android'.")
	}
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

