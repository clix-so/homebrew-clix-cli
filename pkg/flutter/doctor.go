package flutter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
	"github.com/clix-so/clix-cli/pkg/utils"
	"github.com/clix-so/clix-cli/pkg/versions"
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
		{"Firebase CLI Installation", checkFirebaseCLI},
		{"FlutterFire CLI Installation", checkFlutterFireCLI},
		{"Firebase Options Configuration", checkFirebaseOptions},
		{"Clix Flutter SDK Dependency", checkClixDependency},
		{"Firebase Core Dependency", checkFirebaseCoreDependency},
		{"Firebase Messaging Dependency", checkFirebaseMessagingDependency},
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
		return false, "Add 'clix_flutter: " + versions.FlutterClixSDKVersion + "' to dependencies in pubspec.yaml"
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
		return false, "Add 'firebase_core: " + versions.FlutterFirebaseCoreVersion + "' to dependencies in pubspec.yaml"
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
		return false, "Add 'firebase_messaging: " + versions.FlutterFirebaseMessagingVersion + "' to dependencies in pubspec.yaml"
	}

	return true, ""
}

// checkFirebaseCLI verifies Firebase CLI is installed
func checkFirebaseCLI(projectRoot string) (bool, string) {
	if err := utils.RunShellCommand("firebase", "--version"); err != nil {
		return false, "Firebase CLI not installed. Run: npm install -g firebase-tools"
	}
	return true, ""
}

// checkFlutterFireCLI verifies FlutterFire CLI is installed
func checkFlutterFireCLI(projectRoot string) (bool, string) {
	if err := utils.RunShellCommand("flutterfire", "--version"); err != nil {
		return false, "FlutterFire CLI not installed. Run: dart pub global activate flutterfire_cli"
	}
	return true, ""
}

// checkFirebaseOptions verifies firebase_options.dart exists
func checkFirebaseOptions(projectRoot string) (bool, string) {
	configPath := filepath.Join(projectRoot, "lib", "firebase_options.dart")
	if _, err := os.Stat(configPath); err != nil {
		return false, "firebase_options.dart not found. Run: flutterfire configure"
	}

	// Check if Firebase config files exist (they should be created by flutterfire configure)
	androidConfigPath := filepath.Join(projectRoot, "android", "app", "google-services.json")
	iosConfigPath := filepath.Join(projectRoot, "ios", "Runner", "GoogleService-Info.plist")

	var missingFiles []string
	if _, err := os.Stat(androidConfigPath); err != nil {
		missingFiles = append(missingFiles, "android/app/google-services.json")
	}
	if _, err := os.Stat(iosConfigPath); err != nil {
		missingFiles = append(missingFiles, "ios/Runner/GoogleService-Info.plist")
	}

	if len(missingFiles) > 0 {
		return false, fmt.Sprintf("missing Firebase config files: %v. Run: flutterfire configure", missingFiles)
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

	if !strings.Contains(content, "firebase_options.dart") {
		issues = append(issues, "Missing firebase_options.dart import")
	}

	if !strings.Contains(content, "clix_flutter") {
		issues = append(issues, "Missing clix_flutter import")
	}

	if !strings.Contains(content, "Firebase.initializeApp") {
		issues = append(issues, "Missing Firebase.initializeApp call")
	}

	if !strings.Contains(content, "DefaultFirebaseOptions.currentPlatform") {
		issues = append(issues, "Missing DefaultFirebaseOptions.currentPlatform in Firebase.initializeApp")
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
