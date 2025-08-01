package flutter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
)

// RunDoctor checks the Flutter project setup for Clix SDK
func RunDoctor() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}

	fmt.Println("🏥 Clix Flutter Doctor")
	logx.Separatorln()

	allChecks := []struct {
		name string
		fn   func(string) (bool, string)
	}{
		{"Flutter Project Detection", checkFlutterProject},
		{"Clix Flutter SDK Dependency", checkClixDependency},
		{"Firebase Core Dependency", checkFirebaseCoreDependency},
		{"Firebase Messaging Dependency", checkFirebaseMessagingDependency},
		{"Firebase Android Config", checkFirebaseAndroidConfig},
		{"Firebase iOS Config", checkFirebaseIOSConfig},
		{"Main.dart Configuration", checkMainDartConfiguration},
	}

	allPassed := true
	
	for _, check := range allChecks {
		passed, message := check.fn(projectRoot)
		if passed {
			logx.Log().Branch().Success().Println(fmt.Sprintf("%s ✅", check.name))
		} else {
			logx.Log().Branch().Failure().Println(fmt.Sprintf("%s ❌", check.name))
			if message != "" {
				logx.Log().Indent(6).Code().Println(message)
			}
			allPassed = false
		}
	}

	logx.NewLine()
	
	if allPassed {
		fmt.Println("🎉 All checks passed! Your Flutter project is properly configured for Clix SDK.")
	} else {
		fmt.Println("❗ Some checks failed. Please fix the issues above and run 'clix doctor --flutter' again.")
		return fmt.Errorf("flutter doctor checks failed")
	}

	return nil
}

// checkFlutterProject verifies this is a Flutter project
func checkFlutterProject(projectRoot string) (bool, string) {
	pubspecPath := filepath.Join(projectRoot, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); err != nil {
		return false, "pubspec.yaml not found - not a Flutter project"
	}

	data, err := os.ReadFile(pubspecPath)
	if err != nil {
		return false, "failed to read pubspec.yaml"
	}

	content := string(data)
	if !strings.Contains(content, "flutter:") {
		return false, "flutter dependency not found in pubspec.yaml"
	}

	return true, ""
}

// checkClixDependency verifies Clix Flutter SDK dependency
func checkClixDependency(projectRoot string) (bool, string) {
	pubspecPath := filepath.Join(projectRoot, "pubspec.yaml")
	data, err := os.ReadFile(pubspecPath)
	if err != nil {
		return false, "failed to read pubspec.yaml"
	}

	content := string(data)
	if !strings.Contains(content, "clix_flutter:") {
		return false, "Add 'clix_flutter: ^0.0.1' to dependencies in pubspec.yaml"
	}

	return true, ""
}

// checkFirebaseCoreDependency verifies Firebase Core dependency
func checkFirebaseCoreDependency(projectRoot string) (bool, string) {
	pubspecPath := filepath.Join(projectRoot, "pubspec.yaml")
	data, err := os.ReadFile(pubspecPath)
	if err != nil {
		return false, "failed to read pubspec.yaml"
	}

	content := string(data)
	if !strings.Contains(content, "firebase_core:") {
		return false, "Add 'firebase_core: ^3.6.0' to dependencies in pubspec.yaml"
	}

	return true, ""
}

// checkFirebaseMessagingDependency verifies Firebase Messaging dependency
func checkFirebaseMessagingDependency(projectRoot string) (bool, string) {
	pubspecPath := filepath.Join(projectRoot, "pubspec.yaml")
	data, err := os.ReadFile(pubspecPath)
	if err != nil {
		return false, "failed to read pubspec.yaml"
	}

	content := string(data)
	if !strings.Contains(content, "firebase_messaging:") {
		return false, "Add 'firebase_messaging: ^15.1.3' to dependencies in pubspec.yaml"
	}

	return true, ""
}

// checkFirebaseAndroidConfig verifies Firebase Android configuration
func checkFirebaseAndroidConfig(projectRoot string) (bool, string) {
	configPath := filepath.Join(projectRoot, "android", "app", "google-services.json")
	if _, err := os.Stat(configPath); err != nil {
		return false, "google-services.json not found at android/app/google-services.json"
	}

	return true, ""
}

// checkFirebaseIOSConfig verifies Firebase iOS configuration
func checkFirebaseIOSConfig(projectRoot string) (bool, string) {
	configPath := filepath.Join(projectRoot, "ios", "Runner", "GoogleService-Info.plist")
	if _, err := os.Stat(configPath); err != nil {
		return false, "GoogleService-Info.plist not found at ios/Runner/GoogleService-Info.plist"
	}

	return true, ""
}


// checkMainDartConfiguration verifies main.dart has proper Clix setup
func checkMainDartConfiguration(projectRoot string) (bool, string) {
	mainPath := filepath.Join(projectRoot, "lib", "main.dart")
	if _, err := os.Stat(mainPath); err != nil {
		return false, "main.dart not found at lib/main.dart"
	}

	data, err := os.ReadFile(mainPath)
	if err != nil {
		return false, "failed to read main.dart"
	}

	content := string(data)
	issues := []string{}

	if !strings.Contains(content, "firebase_core") {
		issues = append(issues, "Missing firebase_core import")
	}

	if !strings.Contains(content, "clix_flutter") {
		issues = append(issues, "Missing clix_flutter import")
	}

	if !strings.Contains(content, "Firebase.initializeApp()") {
		issues = append(issues, "Missing Firebase.initializeApp() call")
	}

	if !strings.Contains(content, "Clix.initialize") {
		issues = append(issues, "Missing Clix.initialize call")
	}

	if !strings.Contains(content, "WidgetsFlutterBinding.ensureInitialized()") {
		issues = append(issues, "Missing WidgetsFlutterBinding.ensureInitialized() call")
	}

	if len(issues) > 0 {
		return false, strings.Join(issues, "; ")
	}

	return true, ""
}