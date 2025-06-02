package android

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// AddGradleRepository tries to insert mavenCentral() into settings.gradle(.kts) or build.gradle(.kts)
func AddGradleRepository(projectRoot string) bool {
	gradleFiles := []string{
		filepath.Join(projectRoot, "settings.gradle"),
		filepath.Join(projectRoot, "settings.gradle.kts"),
		filepath.Join(projectRoot, "build.gradle"),
		filepath.Join(projectRoot, "build.gradle.kts"),
	}
	for _, file := range gradleFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(data)
		if Contains(content, "repositories") && Contains(content, "mavenCentral()") {
			return true // already present
		}
		// Try to insert after 'repositories {' or at end
		if idx := IndexOf(content, "repositories {"); idx != -1 {
			insertAt := idx + len("repositories {")
			newContent := content[:insertAt] + "\n    mavenCentral()" + content[insertAt:]
			err = ioutil.WriteFile(file, []byte(newContent), 0644)
			return err == nil
		}
	}
	return false
}

// AddGradleDependency tries to insert the Clix SDK dependency into app/build.gradle(.kts)
func AddGradleDependency(projectRoot string) bool {
	gradleFiles := []string{
		filepath.Join(projectRoot, "app", "build.gradle"),
		filepath.Join(projectRoot, "app", "build.gradle.kts"),
	}
	for _, file := range gradleFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(data)
		if Contains(content, "implementation(\"so.clix:clix-android-sdk:0.0.2\")") {
			return true // already present
		}
		// Try to insert after 'dependencies {' or at end
		if idx := IndexOf(content, "dependencies {"); idx != -1 {
			insertAt := idx + len("dependencies {")
			newContent := content[:insertAt] + "\n    implementation(\"so.clix:clix-android-sdk:0.0.2\")" + content[insertAt:]
			err = ioutil.WriteFile(file, []byte(newContent), 0644)
			return err == nil
		}
	}
	return false
}

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
		fmt.Println("[OK] repositories { mavenCentral() } found in Gradle config.")
		return true
	}

	fmt.Println("[FAIL] repositories { mavenCentral() } not found in settings.gradle(.kts) or build.gradle(.kts). Please add:")
	fmt.Println(`repositories {\n    mavenCentral()\n}`)
	return false
}

// CheckGradleDependency checks if so.clix:clix-android-sdk is present in app/build.gradle(.kts)
func CheckGradleDependency(projectRoot string) bool {
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
		if Contains(content, "implementation(\"so.clix:clix-android-sdk:") {
			found = true
			break
		}
		if Contains(content, "implementation(libs.clix.android.sdk)") {
			found = true
			break
		}
	}

	if found {
		fmt.Println("[OK] Clix SDK dependency found in app/build.gradle(.kts).")
		return true
	}

	fmt.Println("[FAIL] Clix SDK dependency not found in app/build.gradle(.kts). Please add:")
	fmt.Println(`dependencies {\n    implementation(\"so.clix:clix-android-sdk:0.0.2\")\n}`)
	return false
}
