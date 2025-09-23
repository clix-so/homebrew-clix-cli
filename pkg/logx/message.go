package logx

const (
	// Gradle Repository
	TitleGradleRepoCheck = "Checking Gradle repository configuration..."
	MsgGradleRepoSuccess = "Gradle repositories are properly configured."
	MsgGradleRepoFailure = "Gradle repository configuration is missing."
	FixGradleRepo = "To resolve this, add the following to settings.gradle(.kts) or build.gradle(.kts):"
	MsgGradleRepoFixFailure = "Automatic fix failed. Please add the following to settings.gradle(.kts) or build.gradle(.kts):"
	CodeGradleRepo = `repositories {
	mavenCentral()
}`

	// Clix SDK Dependency
	TitleClixDependencyCheck = "Checking Clix SDK dependency..."
	MsgClixDependencySuccess = "Clix SDK dependency detected."
	MsgClixDependencyFailure = "Clix SDK dependency is missing."
	FixClixDependency = "To resolve this, add the following to app/build.gradle(.kts):"
	MsgClixDependencyFixFailure = "Automatic fix failed. Please add the following to app/build.gradle(.kts):"
	CodeClixDependency = `dependencies {
	implementation("so.clix:clix-android-sdk:1.1.2")
}`

	// Google Services Plugin
	TitleGmsPluginCheck = "Checking for Google Services plugin..."
	MsgGmsPluginFound = "Google Services plugin detected."
	MsgGmsPluginNotFound = "Google Services plugin is missing."
	FixGmsPlugin = "To resolve this, add the following to build.gradle(.kts):"
	MsgGmsPluginFixFailure = "Automatic fix failed. Please add the following to app/build.gradle(.kts):"
	CodeGmsPlugin = `plugins {
	id("com.google.gms.google-services") version "4.4.2"
}`

	// Clix Application Import & Initialization
	TitleClixInitializationCheck = "Checking Clix SDK initialization..."
	MsgManifestReadFail = "Unable to read AndroidManifest.xml."
	MsgApplicationClassNotDefined = "Application class is not defined in AndroidManifest.xml."
	MsgApplicationClassMissing = "Application class is missing from expected locations."
	MsgApplicationFileReadFail = "Unable to read Application class file."
	MsgClixInitSuccess = "Clix SDK initialization detected."
	MsgClixInitMissing = "Clix SDK initialization is missing from the Application class."
	MsgAppCreateSuccess = "Fixed: Application class created."
	FixClixInitialization = "To resolve this, follow the guide below:"
	MsgAppFixFailure = "Automatic fix failed. Please follow the guide below:"
	ClixInitializationLink = "https://docs.clix.so/sdk-quickstart-android#2-initialize-clix-with-config"

	// MainActivity Permissions
	TitlePermissionCheck = "Checking permission request implementation..."
	MsgMainActivityNotFound = "MainActivity.java or MainActivity.kt is missing."
	MsgPermissionFound = "MainActivity contains permission request code."
	MsgPermissionMissing = "MainActivity does not contain code to request permissions."
	FixPermissionRequest = "To resolve this, follow the guide below:"
	MsgPermissionFixFailure = "Automatic fix failed. Please follow the guide below:"
	PermissionRequestLink = "https://docs.clix.so/sdk-quickstart-android#2-initialize-clix-with-config"

	// google-services.json
	TitleGoogleServicesJsonCheck = "Checking for google-services.json..."
	MsgGoogleJsonMissing = "google-services.json is missing from app/google-services.json."
	MsgGoogleJsonFound = "google-services.json is present."
	FixGoogleServicesJson = "To resolve this, follow the guide below:"
	MsgGoogleJsonFixFailure = "Automatic fix failed. Please follow the guide below:"
	GoogleServicesJsonLink = "https://docs.clix.so/firebase-setting"

	// Auto-fix messages
	MsgAutoFixSuccess = "Fixed: Changes applied automatically."

	// General
	MsgSourceDirNotFound = "Source directory is missing."
	MsgWorkingDirectoryNotFound = "Working directory is missing."
	MsgAppBuildGradleNotFound = "app/build.gradle(.kts) is missing."
	MsgAppBuildGradleReadFail = "Unable to read app/build.gradle(.kts)."
)
