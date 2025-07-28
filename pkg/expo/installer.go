package expo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
	"github.com/clix-so/clix-cli/pkg/utils"
)

// AppConfig represents the app.json configuration
type AppConfig struct {
	Expo ExpoConfig `json:"expo"`
}

type ExpoConfig struct {
	Name              string           `json:"name"`
	Slug              string           `json:"slug"`
	Version           string           `json:"version"`
	Orientation       string           `json:"orientation"`
	Icon              string           `json:"icon"`
	UserInterfaceStyle string          `json:"userInterfaceStyle"`
	Splash            map[string]any   `json:"splash"`
	Updates           map[string]any   `json:"updates"`
	AssetBundlePatterns []string       `json:"assetBundlePatterns"`
	IOS               map[string]any   `json:"ios"`
	Android           map[string]any   `json:"android"`
	Web               map[string]any   `json:"web"`
	Plugins           []any            `json:"plugins"`
}

// HandleExpoInstall guides the user through the React Native Expo installation process
func HandleExpoInstall(apiKey, projectID string) {
	projectRoot, err := os.Getwd()
	if err != nil {
		logx.Log().Failure().Println("Failed to get current working directory")
		return
	}

	fmt.Println("🚀 Installing Clix SDK for React Native Expo...")
	logx.Separatorln()

	// Check if this is an Expo project
	if !CheckExpoProject(projectRoot) {
		logx.Log().Failure().Println("This doesn't appear to be an Expo project. Please ensure you're in the root of an Expo project.")
		return
	}

	// Step 1: Install expo-dev-client
	logx.Log().WithSpinner().Title().Println("📦 Installing expo-dev-client...")
	if err := utils.RunShellCommand("npx", "expo", "install", "expo-dev-client"); err != nil {
		logx.Log().Branch().Failure().Println("Failed to install expo-dev-client")
		logx.Log().Indent(6).Code().Println("npx expo install expo-dev-client")
		return
	}
	logx.Log().Branch().Success().Println("expo-dev-client installed successfully")
	logx.NewLine()

	// Step 2: Install Firebase modules
	logx.Log().WithSpinner().Title().Println("🔥 Installing Firebase modules...")
	if err := utils.RunShellCommand("npx", "expo", "install", "@react-native-firebase/app", "@react-native-firebase/messaging", "expo-build-properties"); err != nil {
		logx.Log().Branch().Failure().Println("Failed to install Firebase modules")
		logx.Log().Indent(6).Code().Println("npx expo install @react-native-firebase/app @react-native-firebase/messaging expo-build-properties")
		return
	}
	logx.Log().Branch().Success().Println("Firebase modules installed successfully")
	logx.NewLine()

	// Step 3: Install Clix dependencies
	logx.Log().WithSpinner().Title().Println("📱 Installing Clix dependencies...")
	if err := utils.RunShellCommand("npx", "expo", "install", "@clix-so/react-native-sdk", "@notifee/react-native", "react-native-device-info", "react-native-get-random-values", "react-native-mmkv", "uuid"); err != nil {
		logx.Log().Branch().Failure().Println("Failed to install Clix dependencies")
		logx.Log().Indent(6).Code().Println("npx expo install @clix-so/react-native-sdk @notifee/react-native react-native-device-info react-native-get-random-values react-native-mmkv uuid")
		return
	}
	logx.Log().Branch().Success().Println("Clix dependencies installed successfully")
	logx.NewLine()

	// Step 4: Check Firebase configuration files
	logx.Log().WithSpinner().Title().Println("🔧 Checking Firebase configuration files...")
	hasAndroidConfig := CheckFirebaseConfig(projectRoot, "android")
	hasIOSConfig := CheckFirebaseConfig(projectRoot, "ios")

	if !hasAndroidConfig || !hasIOSConfig {
		logx.Log().Branch().Failure().Println("Firebase configuration files missing")
		if !hasAndroidConfig {
			logx.Log().Indent(6).Code().Println("Missing: google-services.json (place in project root)")
		}
		if !hasIOSConfig {
			logx.Log().Indent(6).Code().Println("Missing: GoogleService-Info.plist (place in project root)")
		}
		logx.Log().Indent(6).Code().Println("Download these files from Firebase Console")
		return
	}
	logx.Log().Branch().Success().Println("Firebase configuration files found")
	logx.NewLine()

	// Step 5: Update app.json with Firebase plugin
	logx.Log().WithSpinner().Title().Println("⚙️  Updating app.json configuration...")
	if err := UpdateAppConfig(projectRoot); err != nil {
		logx.Log().Branch().Failure().Println("Failed to update app.json")
		logx.Log().Indent(6).Code().Println(err.Error())
		return
	}
	logx.Log().Branch().Success().Println("app.json updated successfully")
	logx.NewLine()

	// Step 6: Create Clix initialization file
	logx.Log().WithSpinner().Title().Println("🔨 Creating Clix initialization...")
	if err := CreateClixInitialization(projectRoot, apiKey, projectID); err != nil {
		logx.Log().Branch().Failure().Println("Failed to create Clix initialization")
		logx.Log().Indent(6).Code().Println(err.Error())
		return
	}
	logx.Log().Branch().Success().Println("Clix initialization created successfully")
	logx.NewLine()

	// Step 7: Final instructions
	fmt.Println("🎉 Clix SDK installation completed!")
	fmt.Println("==================================================")
	fmt.Println("Next steps:")
	fmt.Println("1. Import and call Clix.initialize() in your App.js/App.tsx")
	fmt.Println("2. Run 'npx expo prebuild --clean' to generate native code")
	fmt.Println("3. Run 'npx expo run:android' or 'npx expo run:ios' to test")
	fmt.Println("4. Run 'clix doctor --expo' to verify your setup")
	fmt.Println("==================================================")
}

// CheckExpoProject checks if the current directory is an Expo project
func CheckExpoProject(projectRoot string) bool {
	appJSONPath := filepath.Join(projectRoot, "app.json")
	if _, err := os.Stat(appJSONPath); err != nil {
		return false
	}

	// Check if package.json exists and contains expo
	packageJSONPath := filepath.Join(projectRoot, "package.json")
	if data, err := os.ReadFile(packageJSONPath); err == nil {
		return strings.Contains(string(data), "expo")
	}

	return false
}

// CheckFirebaseConfig checks if Firebase configuration files exist
func CheckFirebaseConfig(projectRoot, platform string) bool {
	var fileName string
	switch platform {
	case "android":
		fileName = "google-services.json"
	case "ios":
		fileName = "GoogleService-Info.plist"
	default:
		return false
	}

	configPath := filepath.Join(projectRoot, fileName)
	_, err := os.Stat(configPath)
	return err == nil
}

// UpdateAppConfig updates the app.json file with Firebase plugin configuration
func UpdateAppConfig(projectRoot string) error {
	appJSONPath := filepath.Join(projectRoot, "app.json")
	
	data, err := os.ReadFile(appJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read app.json: %v", err)
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse app.json: %v", err)
	}

	// Check if required plugins are already present
	hasFirebaseAppPlugin := false
	hasFirebaseMessagingPlugin := false
	buildPropertiesPluginIndex := -1

	for i, plugin := range config.Expo.Plugins {
		if pluginStr, ok := plugin.(string); ok {
			if pluginStr == "@react-native-firebase/app" {
				hasFirebaseAppPlugin = true
			}
			if pluginStr == "@react-native-firebase/messaging" {
				hasFirebaseMessagingPlugin = true
			}
			if pluginStr == "expo-build-properties" {
				buildPropertiesPluginIndex = i
			}
		}
		if pluginArray, ok := plugin.([]any); ok && len(pluginArray) > 0 {
			if pluginStr, ok := pluginArray[0].(string); ok {
				if pluginStr == "@react-native-firebase/app" {
					hasFirebaseAppPlugin = true
				}
				if pluginStr == "@react-native-firebase/messaging" {
					hasFirebaseMessagingPlugin = true
				}
				if pluginStr == "expo-build-properties" {
					buildPropertiesPluginIndex = i
				}
			}
		}
	}

	// Add missing plugins
	if !hasFirebaseAppPlugin {
		config.Expo.Plugins = append(config.Expo.Plugins, "@react-native-firebase/app")
	}

	if !hasFirebaseMessagingPlugin {
		config.Expo.Plugins = append(config.Expo.Plugins, "@react-native-firebase/messaging")
	}

	// Handle expo-build-properties plugin
	notifeeRepo := "../../node_modules/@notifee/react-native/android/libs"
	
	if buildPropertiesPluginIndex == -1 {
		// Plugin doesn't exist, add complete configuration
		buildPropertiesPlugin := []any{
			"expo-build-properties",
			map[string]any{
				"ios": map[string]any{
					"useFrameworks": "static",
				},
				"android": map[string]any{
					"extraMavenRepos": []string{notifeeRepo},
				},
			},
		}
		config.Expo.Plugins = append(config.Expo.Plugins, buildPropertiesPlugin)
	} else {
		// Plugin exists, ensure it has the correct configuration
		plugin := config.Expo.Plugins[buildPropertiesPluginIndex]
		
		// Handle string plugin format - convert to array format
		if pluginStr, ok := plugin.(string); ok && pluginStr == "expo-build-properties" {
			config.Expo.Plugins[buildPropertiesPluginIndex] = []any{
				"expo-build-properties",
				map[string]any{
					"ios": map[string]any{
						"useFrameworks": "static",
					},
					"android": map[string]any{
						"extraMavenRepos": []string{notifeeRepo},
					},
				},
			}
		} else if pluginArray, ok := plugin.([]any); ok && len(pluginArray) >= 1 {
			// Handle array plugin format
			var pluginConfig map[string]any
			if len(pluginArray) >= 2 {
				if existingConfig, ok := pluginArray[1].(map[string]any); ok {
					pluginConfig = existingConfig
				} else {
					pluginConfig = make(map[string]any)
					pluginArray = append(pluginArray, pluginConfig)
				}
			} else {
				pluginConfig = make(map[string]any)
				pluginArray = append(pluginArray, pluginConfig)
			}
				// Ensure iOS configuration
				if iosConfig, exists := pluginConfig["ios"]; exists {
					if iosMap, ok := iosConfig.(map[string]any); ok {
						iosMap["useFrameworks"] = "static"
					}
				} else {
					pluginConfig["ios"] = map[string]any{
						"useFrameworks": "static",
					}
				}

				// Ensure Android configuration
				if androidConfig, exists := pluginConfig["android"]; exists {
					if androidMap, ok := androidConfig.(map[string]any); ok {
						// Handle extraMavenRepos
						if existingRepos, exists := androidMap["extraMavenRepos"]; exists {
							// Check if notifee repo already exists
							hasNotifeeRepo := false
							if repoArray, ok := existingRepos.([]any); ok {
								for _, repo := range repoArray {
									if repoStr, ok := repo.(string); ok && repoStr == notifeeRepo {
										hasNotifeeRepo = true
										break
									}
								}
								if !hasNotifeeRepo {
									androidMap["extraMavenRepos"] = append(repoArray, notifeeRepo)
								}
							} else if repoSlice, ok := existingRepos.([]string); ok {
								for _, repo := range repoSlice {
									if repo == notifeeRepo {
										hasNotifeeRepo = true
										break
									}
								}
								if !hasNotifeeRepo {
									newRepos := make([]any, len(repoSlice)+1)
									for j, repo := range repoSlice {
										newRepos[j] = repo
									}
									newRepos[len(repoSlice)] = notifeeRepo
									androidMap["extraMavenRepos"] = newRepos
								}
							}
						} else {
							androidMap["extraMavenRepos"] = []string{notifeeRepo}
						}
					}
				} else {
					pluginConfig["android"] = map[string]any{
						"extraMavenRepos": []string{notifeeRepo},
					}
				}
			
			// Update the plugin array in the config
			config.Expo.Plugins[buildPropertiesPluginIndex] = pluginArray
		}
	}

	// Update Android configuration
	if config.Expo.Android == nil {
		config.Expo.Android = make(map[string]any)
	}
	config.Expo.Android["googleServicesFile"] = "./google-services.json"
	
	// Check for Android package name
	if _, exists := config.Expo.Android["package"]; !exists {
		fmt.Println("\n📱 Android package name is required for Firebase configuration.")
		fmt.Println("This should match the package name in your google-services.json file.")
		fmt.Println("Example: com.yourcompany.yourapp")
		androidPackage := utils.Prompt("Enter your Android package name")
		if androidPackage != "" {
			config.Expo.Android["package"] = androidPackage
		}
	}

	// Update iOS configuration
	if config.Expo.IOS == nil {
		config.Expo.IOS = make(map[string]any)
	}
	config.Expo.IOS["googleServicesFile"] = "./GoogleService-Info.plist"
	
	// Check for iOS bundle identifier
	if _, exists := config.Expo.IOS["bundleIdentifier"]; !exists {
		fmt.Println("\n🍎 iOS bundle identifier is required for Firebase configuration.")
		fmt.Println("This should match the bundle ID in your GoogleService-Info.plist file.")
		fmt.Println("Example: com.yourcompany.yourapp")
		iosBundleId := utils.Prompt("Enter your iOS bundle identifier")
		if iosBundleId != "" {
			config.Expo.IOS["bundleIdentifier"] = iosBundleId
		}
	}
	
	// Add iOS push notification entitlements
	if entitlements, exists := config.Expo.IOS["entitlements"]; exists {
		if entMap, ok := entitlements.(map[string]any); ok {
			entMap["aps-environment"] = "production"
		}
	} else {
		config.Expo.IOS["entitlements"] = map[string]any{
			"aps-environment": "production",
		}
	}

	// Add iOS info.plist settings for background modes
	if infoPlist, exists := config.Expo.IOS["infoPlist"]; exists {
		if plistMap, ok := infoPlist.(map[string]any); ok {
			plistMap["UIBackgroundModes"] = []string{"remote-notification"}
		}
	} else {
		config.Expo.IOS["infoPlist"] = map[string]any{
			"UIBackgroundModes": []string{"remote-notification"},
		}
	}

	// Write updated config back to file
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal app.json: %v", err)
	}

	if err := os.WriteFile(appJSONPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write app.json: %v", err)
	}

	return nil
}

// CreateClixInitialization creates a TypeScript file with Clix initialization code
func CreateClixInitialization(projectRoot, apiKey, projectID string) error {
	// Check if the project uses TypeScript
	isTypeScript := false
	if _, err := os.Stat(filepath.Join(projectRoot, "tsconfig.json")); err == nil {
		isTypeScript = true
	}

	var fileName string
	var content string

	if isTypeScript {
		fileName = "clix-config.ts"
		content = fmt.Sprintf(`import Clix from '@clix-so/react-native-sdk';

export const initializeClix = async (): Promise<void> => {
  try {
    await Clix.initialize({
      projectId: '%s',
      apiKey: '%s'
    });
    console.log('Clix SDK initialized successfully');
  } catch (error) {
    console.error('Failed to initialize Clix SDK:', error);
  }
};
`, projectID, apiKey)
	} else {
		fileName = "clix-config.js"
		content = fmt.Sprintf(`import Clix from '@clix-so/react-native-sdk';

export const initializeClix = async () => {
  try {
    await Clix.initialize({
      projectId: '%s',
      apiKey: '%s'
    });
    console.log('Clix SDK initialized successfully');
  } catch (error) {
    console.error('Failed to initialize Clix SDK:', error);
  }
};
`, projectID, apiKey)
	}

	filePath := filepath.Join(projectRoot, fileName)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create %s: %v", fileName, err)
	}

	return nil
}