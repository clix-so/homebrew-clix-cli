package flutter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
	"github.com/clix-so/clix-cli/pkg/utils"
)

// PubspecConfig represents the pubspec.yaml configuration
type PubspecConfig struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Version      string            `yaml:"version"`
	Environment  map[string]string `yaml:"environment"`
	Dependencies map[string]any    `yaml:"dependencies"`
}

// HandleFlutterInstall guides the user through the Flutter SDK installation process
func HandleFlutterInstall(apiKey, projectID string) {
	projectRoot, err := os.Getwd()
	if err != nil {
		logx.Log().Failure().Println("Failed to get current working directory")
		return
	}

	logx.Log().Title().Println("Installing Clix SDK for Flutter…")
	logx.Separatorln()

	// Check if this is a Flutter project
	if !CheckFlutterProject(projectRoot) {
		logx.Log().Failure().Println("This doesn't appear to be a Flutter project. Please ensure you're in the root of a Flutter project.")
		return
	}

	// Step 1: Check and install Firebase CLI
	logx.Log().WithSpinner().Title().Println("Checking Firebase CLI…")
	if err := CheckAndInstallFirebaseCLI(); err != nil {
		logx.Log().Branch().Failure().Println("Failed to setup Firebase CLI")
		logx.Log().Indent(6).Code().Println(err.Error())
		return
	}
	logx.Log().Branch().Success().Println("Firebase CLI is ready")
	logx.NewLine()

	// Step 2: Check and install FlutterFire CLI
	logx.Log().WithSpinner().Title().Println("Checking FlutterFire CLI…")
	if err := CheckAndInstallFlutterFireCLI(); err != nil {
		logx.Log().Branch().Failure().Println("Failed to setup FlutterFire CLI")
		logx.Log().Indent(6).Code().Println(err.Error())
		return
	}
	logx.Log().Branch().Success().Println("FlutterFire CLI is ready")
	logx.NewLine()

	// Step 3: Configure Firebase project
	logx.Log().WithSpinner().Title().Println("Configuring Firebase project…")
	if err := ConfigureFirebaseProject(); err != nil {
		logx.Log().Branch().Failure().Println("Failed to configure Firebase project")
		logx.Log().Indent(6).Code().Println(err.Error())
		logx.Log().Warn().Println("Please manually run 'flutterfire configure' to setup Firebase")
		return
	}
	logx.Log().Branch().Success().Println("Firebase project configured")
	logx.NewLine()

	// Step 4: Add dependencies to pubspec.yaml
	logx.Log().WithSpinner().Title().Println("Adding dependencies to pubspec.yaml…")
	if err := AddFlutterDependencies(projectRoot); err != nil {
		logx.Log().Branch().Failure().Println("Failed to add dependencies to pubspec.yaml")
		logx.Log().Indent(6).Code().Println(err.Error())
		logx.Log().Warn().Println("Please manually add the following dependencies to your pubspec.yaml:")
		logx.Log().Indent(2).Println("dependencies:")
		logx.Log().Indent(4).Code().Println("clix_flutter: ^0.0.1")
		logx.Log().Indent(4).Code().Println("firebase_core: ^3.6.0")
		logx.Log().Indent(4).Code().Println("firebase_messaging: ^15.1.3")
		return
	}
	logx.Log().Branch().Success().Println("Dependencies added to pubspec.yaml")
	logx.NewLine()

	// Step 5: Install dependencies
	logx.Log().WithSpinner().Title().Println("Installing Flutter dependencies…")
	if err := utils.RunShellCommand("flutter", "pub", "get"); err != nil {
		logx.Log().Branch().Failure().Println("Failed to install Flutter dependencies")
		logx.Log().Indent(6).Code().Println("flutter pub get")
		return
	}
	logx.Log().Branch().Success().Println("Flutter dependencies installed successfully")
	logx.NewLine()

	// Step 6: Verify Firebase configuration
	logx.Log().WithSpinner().Title().Println("Verifying Firebase configuration…")
	if err := VerifyFirebaseConfig(projectRoot); err != nil {
		logx.Log().Branch().Failure().Println("Firebase configuration verification failed")
		logx.Log().Indent(6).Code().Println(err.Error())
		return
	}
	logx.Log().Branch().Success().Println("Firebase configuration verified")
	logx.NewLine()

	// Step 7: Update main.dart with Clix initialization
	logx.Log().WithSpinner().Title().Println("Updating main.dart with Clix initialization…")
	if err := UpdateMainDart(projectRoot, projectID, apiKey); err != nil {
		logx.Log().Branch().Failure().Println("Failed to update main.dart")
		logx.Log().Indent(6).Code().Println(err.Error())
		logx.Log().Warn().Println("Please manually add the following to your main.dart:")
		logx.Log().Indent(3).Code().Println("import 'package:firebase_core/firebase_core.dart';")
		logx.Log().Indent(3).Code().Println("import 'package:clix_flutter/clix_flutter.dart';")
		logx.Log().Indent(3).Println("// Add Firebase.initializeApp() and Clix.initialize() in main() before runApp()")
	} else {
		logx.Log().Branch().Success().Println("main.dart updated successfully")
	}
	logx.NewLine()

	// Step 8: iOS-specific setup instructions
	logx.Log().Title().Println("iOS-specific setup required")
	logx.Separatorln()
	logx.Log().Indent(2).Println("1. Open ios/Runner.xcworkspace in Xcode")
	logx.Log().Indent(2).Println("2. Select Runner target > Signing & Capabilities")
	logx.Log().Indent(2).Println("3. Add 'Push Notifications' capability")
	logx.Log().Indent(2).Println("4. Add 'Background Modes' capability")
	logx.Log().Indent(2).Println("5. Enable 'Remote notifications' in Background Modes")
	logx.Separatorln()

	// Step 9: Final instructions
	logx.Log().Success().Println("Clix SDK Flutter installation completed!")
	logx.Separatorln()
	logx.Log().Title().Println("Next steps")
	logx.Log().Indent(2).Println("1. Configure iOS push notifications in Xcode (as shown above)")
	logx.Log().Indent(2).Println("2. Upload your iOS Service Account Key to Clix console")
	logx.Log().Indent(2).Println("3. Run 'flutter run' to test your app")
	logx.Log().Indent(2).Println("4. Run 'clix doctor --flutter' to verify your setup")
	logx.Separatorln()
}

// CheckFlutterProject checks if the current directory is a Flutter project
func CheckFlutterProject(projectRoot string) bool {
	pubspecPath := filepath.Join(projectRoot, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); err != nil {
		return false
	}

	// Check if pubspec.yaml contains flutter dependency
	if data, err := os.ReadFile(pubspecPath); err == nil {
		content := string(data)
		return strings.Contains(content, "flutter:") || strings.Contains(content, "flutter_test:")
	}

	return false
}

// CheckFirebaseConfig checks if Firebase configuration files exist
func CheckFirebaseConfig(projectRoot, platform string) bool {
	var configPath string
	switch platform {
	case "android":
		configPath = filepath.Join(projectRoot, "android", "app", "google-services.json")
	case "ios":
		configPath = filepath.Join(projectRoot, "ios", "Runner", "GoogleService-Info.plist")
	default:
		return false
	}

	_, err := os.Stat(configPath)
	return err == nil
}

// AddFlutterDependencies adds required dependencies to pubspec.yaml
func AddFlutterDependencies(projectRoot string) error {
	pubspecPath := filepath.Join(projectRoot, "pubspec.yaml")

	data, err := os.ReadFile(pubspecPath)
	if err != nil {
		return fmt.Errorf("failed to read pubspec.yaml: %v", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var result []string

	dependenciesFound := false
	dependenciesAdded := false

	requiredDeps := map[string]string{
		"clix_flutter":       "^0.0.1",
		"firebase_core":      "^3.6.0",
		"firebase_messaging": "^15.1.3",
	}

	for _, line := range lines {
		result = append(result, line)

		// Find dependencies section
		if strings.TrimSpace(line) == "dependencies:" {
			dependenciesFound = true
		}

		// Add dependencies after finding the dependencies section
		if dependenciesFound && !dependenciesAdded {
			// Check if this line starts a new section (not indented under dependencies)
			if strings.TrimSpace(line) != "dependencies:" &&
				strings.TrimSpace(line) != "" &&
				!strings.HasPrefix(line, "  ") &&
				!strings.HasPrefix(line, "\t") {
				// We've moved to a new section, add dependencies before this line
				result = result[:len(result)-1] // Remove the current line

				// Add required dependencies
				for dep, version := range requiredDeps {
					if !strings.Contains(content, dep+":") {
						result = append(result, fmt.Sprintf("  %s: %s", dep, version))
					}
				}
				result = append(result, line) // Add back the current line
				dependenciesAdded = true
			}
		}
	}

	// If dependencies section was found but we reached the end, add dependencies
	if dependenciesFound && !dependenciesAdded {
		for dep, version := range requiredDeps {
			if !strings.Contains(content, dep+":") {
				result = append(result, fmt.Sprintf("  %s: %s", dep, version))
			}
		}
	}

	// If no dependencies section found, add it
	if !dependenciesFound {
		return fmt.Errorf("dependencies section not found in pubspec.yaml")
	}

	// Write updated pubspec.yaml
	updatedContent := strings.Join(result, "\n")
	if err := os.WriteFile(pubspecPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write pubspec.yaml: %v", err)
	}

	return nil
}

// UpdateMainDart updates the main.dart file to include Clix and Firebase initialization
func UpdateMainDart(projectRoot, projectID, apiKey string) error {
	mainDartPath := filepath.Join(projectRoot, "lib", "main.dart")

	// Check if main.dart exists
	if _, err := os.Stat(mainDartPath); err != nil {
		return fmt.Errorf("main.dart not found at %s", mainDartPath)
	}

	// Read existing main.dart
	content, err := os.ReadFile(mainDartPath)
	if err != nil {
		return fmt.Errorf("failed to read main.dart: %v", err)
	}

	mainContent := string(content)

	// Check if Clix is already initialized
	if strings.Contains(mainContent, "Clix.initialize") {
		return nil // Already integrated
	}

	// Add necessary imports and initialization
	modifiedContent, err := addClixToMainDart(mainContent, projectID, apiKey)
	if err != nil {
		return fmt.Errorf("failed to modify main.dart: %v", err)
	}

	// Write the modified content back
	if err := os.WriteFile(mainDartPath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("failed to write modified main.dart: %v", err)
	}

	return nil
}

// addClixToMainDart adds Clix and Firebase imports and initialization to main.dart
func addClixToMainDart(content, projectID, apiKey string) (string, error) {
	lines := strings.Split(content, "\n")
	var result []string

	// Imports to add
	firebaseImport := "import 'package:firebase_core/firebase_core.dart';"
	firebaseOptionsImport := "import 'firebase_options.dart';"
	clixImport := "import 'package:clix_flutter/clix_flutter.dart';"

	importsAdded := false
	mainFunctionModified := false

	for i, line := range lines {
		// Add imports after existing import statements
		if !importsAdded && strings.HasPrefix(strings.TrimSpace(line), "import ") {
			result = append(result, line)

			// Check if this is the last import
			isLastImport := true
			for j := i + 1; j < len(lines); j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine == "" || strings.HasPrefix(nextLine, "//") {
					continue
				}
				if strings.HasPrefix(nextLine, "import ") {
					isLastImport = false
					break
				}
				break
			}

			if isLastImport {
				if !strings.Contains(content, "firebase_core") {
					result = append(result, firebaseImport)
				}
				if !strings.Contains(content, "firebase_options.dart") {
					result = append(result, firebaseOptionsImport)
				}
				if !strings.Contains(content, "clix_flutter") {
					result = append(result, clixImport)
				}
				importsAdded = true
			}
		} else if strings.Contains(line, "void main()") && !mainFunctionModified {
			// Modify main function to be async and add initialization
			if strings.Contains(line, "async") {
				// Already async, just add the line
				result = append(result, line)
			} else {
				// Make it async
				modifiedLine := strings.Replace(line, "void main()", "void main() async", 1)
				result = append(result, modifiedLine)
			}

			// Add initialization code after the opening brace
			result = append(result, "  WidgetsFlutterBinding.ensureInitialized();")
			result = append(result, "")
			result = append(result, "  await Firebase.initializeApp(")
			result = append(result, "    options: DefaultFirebaseOptions.currentPlatform,")
			result = append(result, "  );")
			result = append(result, "")
			result = append(result, "  await Clix.initialize(const ClixConfig(")
			result = append(result, fmt.Sprintf("    projectId: '%s',", projectID))
			result = append(result, fmt.Sprintf("    apiKey: '%s',", apiKey))
			result = append(result, "  ));")
			result = append(result, "")

			mainFunctionModified = true
		} else {
			result = append(result, line)
		}
	}

	// If imports weren't added at the beginning, add them
	if !importsAdded {
		imports := []string{firebaseImport, firebaseOptionsImport, clixImport, ""}
		result = append(imports, result...)
	}

	return strings.Join(result, "\n"), nil
}

// CheckAndInstallFirebaseCLI checks if Firebase CLI is installed and installs it if needed
func CheckAndInstallFirebaseCLI() error {
	// Check if Firebase CLI is already installed
	if err := utils.RunShellCommand("firebase", "--version"); err == nil {
		return nil
	}

	fmt.Println("Firebase CLI not found. Installing Firebase CLI...")

	// Try to install Firebase CLI via npm
	if err := utils.RunShellCommand("npm", "install", "-g", "firebase-tools"); err != nil {
		return fmt.Errorf("failed to install Firebase CLI via npm. Please install manually: npm install -g firebase-tools")
	}

	// Verify installation
	if err := utils.RunShellCommand("firebase", "--version"); err != nil {
		return fmt.Errorf("Firebase CLI installation failed. Please install manually: npm install -g firebase-tools")
	}

	return nil
}

// CheckAndInstallFlutterFireCLI checks if FlutterFire CLI is installed and installs it if needed
func CheckAndInstallFlutterFireCLI() error {
	// Check if FlutterFire CLI is already installed
	if err := utils.RunShellCommand("flutterfire", "--version"); err == nil {
		return nil
	}

	fmt.Println("FlutterFire CLI not found. Installing FlutterFire CLI...")

	// Install FlutterFire CLI
	if err := utils.RunShellCommand("dart", "pub", "global", "activate", "flutterfire_cli"); err != nil {
		return fmt.Errorf("failed to install FlutterFire CLI. Please install manually: dart pub global activate flutterfire_cli")
	}

	// Verify installation
	if err := utils.RunShellCommand("flutterfire", "--version"); err != nil {
		return fmt.Errorf("FlutterFire CLI installation failed. Please ensure dart is in PATH and run: dart pub global activate flutterfire_cli")
	}

	return nil
}

// ConfigureFirebaseProject runs flutterfire configure to setup Firebase project
func ConfigureFirebaseProject() error {
	// Check if firebase_options.dart already exists
	if _, err := os.Stat("lib/firebase_options.dart"); err == nil {
		fmt.Println("Firebase project appears to be already configured (firebase_options.dart found)")
		return nil
	}

	fmt.Println("Running 'flutterfire configure' to setup Firebase project...")
	fmt.Println("Please follow the interactive prompts to:")
	fmt.Println("1. Select your Firebase project")
	fmt.Println("2. Choose platforms (iOS and Android)")
	fmt.Println("3. Configure bundle IDs")

	// Run flutterfire configure interactively
	if err := utils.RunShellCommand("flutterfire", "configure"); err != nil {
		return fmt.Errorf("flutterfire configure failed. Please run manually: flutterfire configure")
	}

	// Verify firebase_options.dart was created
	if _, err := os.Stat("lib/firebase_options.dart"); err != nil {
		return fmt.Errorf("firebase_options.dart was not created. Please run 'flutterfire configure' manually")
	}

	return nil
}

// VerifyFirebaseConfig verifies that Firebase is properly configured
func VerifyFirebaseConfig(projectRoot string) error {
	// Check if firebase_options.dart exists
	firebaseOptionsPath := filepath.Join(projectRoot, "lib", "firebase_options.dart")
	if _, err := os.Stat(firebaseOptionsPath); err != nil {
		return fmt.Errorf("firebase_options.dart not found. Please run 'flutterfire configure' first")
	}

	// Check if Firebase configuration files exist in their proper locations
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
		return fmt.Errorf("missing Firebase config files: %v. These should be automatically created by 'flutterfire configure'", missingFiles)
	}

	return nil
}
