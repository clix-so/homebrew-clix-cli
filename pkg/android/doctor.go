package android

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clix-so/clix-cli/pkg/utils"
)

// CheckClixCoreImport checks if any Application class imports so.clix.core.Clix (Java or Kotlin).
func CheckClixCoreImport(projectRoot string) bool {
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	appFiles := []string{}

	// Helper: recursively find *Application.java and *Application.kt
	findAppFiles := func(root string) {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			
			extension := filepath.Ext(path)

			isFile := !info.IsDir()
			isJavaApplication := extension == ".java" && len(info.Name()) >= len("Application.java") && info.Name()[len(info.Name())-len("Application.java"):] == "Application.java"
			isKotlinApplication := extension == ".kt" && len(info.Name()) >= len("Application.kt") && info.Name()[len(info.Name())-len("Application.kt"):] == "Application.kt"

			if isFile && (isJavaApplication || isKotlinApplication) {
				appFiles = append(appFiles, path)
			}
			return nil
		})
	}

	findAppFiles(javaDir)
	findAppFiles(kotlinDir)

	if len(appFiles) == 0 {
		utils.Warnln("No Application class found under app/src/main/java or app/src/main/kotlin.") // TODO: add following action
		return false
	}

	importFound := false
	initializeFound := false
	for _, file := range appFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(data)
		if stringContainsImportClix(content) {
			importFound = true
		}
		if StringContainsClixInitializeInOnCreate(content) {
			initializeFound = true
		}
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
		return false
	}

	return true
}

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

	// 3. Check permission request in MainActivity
	utils.TitlelnWithSpinner("Checking permission request...")
	if !CheckAndroidMainActivityPermissions(projectRoot) {
		utils.Indentln("To fix this, add the following to MainActivity.java or MainActivity.kt:", 3)
		fmt.Println()
		utils.Grayln(`ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.POST_NOTIFICATIONS), 1001)`)
	}
	fmt.Println()

	// 4. Check google-services.json existence
	utils.TitlelnWithSpinner("Checking google-services.json...")
	CheckGoogleServicesJSON(projectRoot)
	fmt.Println()
}


