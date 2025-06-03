package android

import (
	"fmt"

	"github.com/clix-so/clix-cli/pkg/utils"
)

// stringContainsImportClix checks if the given file content contains the import statement for so.clix.core.Clix
func stringContainsImportClix(content string) bool {
	return len(content) > 0 && Contains(content, "import so.clix.core.Clix")
}

// StringContainsClixInitializeInOnCreate checks if Clix.initialize(this, ...) is called inside onCreate
func StringContainsClixInitializeInOnCreate(content string) bool {
	// Simple heuristic: check for 'void onCreate' or 'fun onCreate', then 'Clix.initialize(this'
	onCreateIdx := -1
	if idx := IndexOf(content, "void onCreate"); idx != -1 {
		onCreateIdx = idx
	} else if idx := IndexOf(content, "fun onCreate"); idx != -1 {
		onCreateIdx = idx
	}
	if onCreateIdx == -1 {
		return false
	}
	// Check 200 chars after onCreate for Clix.initialize(this
	endIdx := onCreateIdx + 200
	if endIdx > len(content) {
		endIdx = len(content)
	}
	return Contains(content[onCreateIdx:endIdx], "Clix.initialize(this")
}

// Contains is a helper for substring check (to avoid strings package import)
func Contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && (func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})()
}

// IndexOf returns the index of the first instance of substr in s, or -1 if not present
func IndexOf(s, substr string) int {
	if len(substr) == 0 || len(s) < len(substr) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}



// RunDoctor runs all Android doctor checks.
func RunDoctor(projectRoot string) {
	utils.TitlelnWithSpinner("Checking Gradle repository settings...")
	if !CheckGradleRepository(projectRoot) {
		utils.Indentln("To fix this, add the following to settings.gradle(.kts) or build.gradle(.kts):", 3)
		fmt.Println()
		utils.Grayln(`   repositories {
	   mavenCentral()
   }`)
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking for Clix SDK dependency...")
	if !CheckGradleDependency(projectRoot) {
		utils.Indentln("To fix this, add the following to app/build.gradle(.kts):", 3)
		fmt.Println()
		utils.Grayln(`   dependencies {
       implementation("so.clix:clix-android-sdk:0.0.2")
   }`)
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking for Google Services plugin...")
	if !CheckGradlePlugin(projectRoot) {
		utils.Indentln("To fix this, add the following to build.gradle(.kts):", 3)
		fmt.Println()
		utils.Grayln(`   plugins {
       id("com.google.gms.google-services") version "4.4.2"
   }`)
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking Clix SDK initialization...")
	CheckClixCoreImport(projectRoot)
	fmt.Println()

	utils.TitlelnWithSpinner("Checking permission request...")
	if !CheckAndroidMainActivityPermissions(projectRoot) {
		utils.Indentln("To fix this, add the following to MainActivity.java or MainActivity.kt:", 3)
		fmt.Println()
		utils.Grayln(`   ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.POST_NOTIFICATIONS), 1001)`)
	}
	fmt.Println()

	utils.TitlelnWithSpinner("Checking google-services.json...")
	CheckGoogleServicesJSON(projectRoot)
	fmt.Println()
}
