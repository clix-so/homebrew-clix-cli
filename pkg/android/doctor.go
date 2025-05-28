package android

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckClixCoreImport checks if any Application class imports so.clix.core.Clix (Java or Kotlin).
func CheckClixCoreImport(projectRoot string) error {
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	appFiles := []string{}

	// Helper: recursively find *Application.java and *Application.kt
	findAppFiles := func(root string) {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (filepath.Ext(path) == ".java" || filepath.Ext(path) == ".kt") &&
				len(info.Name()) >= len("Application.java") &&
				info.Name()[len(info.Name())-len("Application.java"):] == "Application.java" ||
				info.Name()[len(info.Name())-len("Application.kt"):] == "Application.kt" {
				appFiles = append(appFiles, path)
			}
			return nil
		})
	}

	findAppFiles(javaDir)
	findAppFiles(kotlinDir)

	if len(appFiles) == 0 {
		fmt.Println("[WARN] No Application class found under app/src/main/java or app/src/main/kotlin.")
		return nil
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
		fmt.Println("[OK] so.clix.core.Clix is imported in Application class.")
	} else {
		fmt.Println("[FAIL] so.clix.core.Clix is not imported in any Application class.")
	}

	if initializeFound {
		fmt.Println("[OK] Clix.initialize(this, ...) is called in onCreate() of Application class.")
	} else {
		fmt.Println("[FAIL] Clix.initialize(this, ...) is NOT called in onCreate() of any Application class.")
	}

	if !importFound || !initializeFound {
		return fmt.Errorf("Application class missing required import or initialization")
	}
	return nil
}

// stringContainsImportClix checks if the given file content contains the import statement for so.clix.core.Clix
func stringContainsImportClix(content string) bool {
	return (len(content) > 0 && (Contains(content, "import so.clix.core.Clix") || Contains(content, "import so.clix.core.Clix;")))
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
func CheckGoogleServicesJSON(projectRoot string) error {
	gsPath := filepath.Join(projectRoot, "app", "google-services.json")
	if _, err := os.Stat(gsPath); os.IsNotExist(err) {
		fmt.Println("[FAIL] google-services.json not found at app/google-services.json.")
		return fmt.Errorf("google-services.json not found")
	}
	fmt.Println("[OK] google-services.json found at app/google-services.json.")
	return nil
}

// CheckMainActivityPermissionRequest checks if MainActivity requests permissions (Java or Kotlin)
func CheckMainActivityPermissionRequest(projectRoot string) error {
	mainActivityFiles := []string{}
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")

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
		fmt.Println("[WARN] No MainActivity.java or MainActivity.kt found.")
		return nil
	}

	permissionPattern := []string{
		"requestPermissions(", // AndroidX, API 23+
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
		fmt.Println("[OK] MainActivity contains code requesting permissions.")
		return nil
	}
	fmt.Println("[FAIL] MainActivity does NOT contain code requesting permissions.")
	return fmt.Errorf("No permission request code found in MainActivity")
}

// RunDoctor runs all Android doctor checks.
func RunDoctor(projectRoot string) {
	errs := []error{}

	// 1. Check Gradle repository and dependency
	if !CheckGradleRepository(projectRoot) {
		errs = append(errs, fmt.Errorf("repositories { mavenCentral() } not found in Gradle config"))
	}
	if !CheckGradleDependency(projectRoot) {
		errs = append(errs, fmt.Errorf("Clix SDK dependency not found in app/build.gradle(.kts)"))
	}

	// 2. Check Clix core import in Application class
	if err := CheckClixCoreImport(projectRoot); err != nil {
		errs = append(errs, err)
	}

	// 3. Check permission request in MainActivity
	if err := CheckMainActivityPermissionRequest(projectRoot); err != nil {
		errs = append(errs, err)
	}

	// 4. Check google-services.json existence
	if err := CheckGoogleServicesJSON(projectRoot); err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		fmt.Println("\n[Android Doctor] All checks passed!")
	} else {
		fmt.Println("\n[Android Doctor] Some checks failed:")
		for _, err := range errs {
			fmt.Println(" -", err)
		}
	}
}


