package logx

const (
	// Gradle Repository
	MsgGradleRepoSuccess = "Gradle repositories are properly configured."
	MsgGradleRepoFailure = "Gradle repository configuration is missing."
	MsgGradleRepoFixFailure = "Automatic fix failed. Please add the following to settings.gradle(.kts) or build.gradle(.kts):"

	// Clix SDK Dependency
	MsgClixDependencySuccess = "Clix SDK dependency detected."
	MsgClixDependencyFailure = "Clix SDK dependency is missing."
	MsgClixDependencyFixFailure = "Automatic fix failed. Please add the following to app/build.gradle(.kts):"

	// Google Services Plugin
	MsgGmsPluginFound = "Google Services plugin detected."
	MsgGmsPluginNotFound = "Google Services plugin is missing."
	MsgGmsPluginFixFailure = "Automatic fix failed. Please add the following to app/build.gradle(.kts):"

	// Clix Application Import & Initialization
	MsgManifestReadFail = "Unable to read AndroidManifest.xml."
	MsgApplicationClassNotDefined = "Application class is not defined in AndroidManifest.xml."
	MsgApplicationClassMissing = "Application class is missing from expected locations."
	MsgApplicationFileReadFail = "Unable to read Application class file."
	MsgClixInitSuccess = "Clix SDK initialization detected."
	MsgClixInitMissing = "Clix SDK initialization is missing from the Application class."
	MsgAppCreateSuccess = "Fixed: Application class created."
	MsgAppCannotFix = "Automatic fix failed. Please follow the guide below to set up your Application class:"
	MsgAppManualGuideLink = "https://docs.clix.so/sdk-quickstart-android#setup-clix-manual-installation"

	// MainActivity Permissions
	MsgMainActivityNotFound = "MainActivity.java or MainActivity.kt is missing."
	MsgPermissionFound = "MainActivity contains permission request code."
	MsgPermissionMissing = "MainActivity does not contain code to request permissions."
	MsgPermissionFixFailure = "Automatic fix failed. Please add the following to your MainActivity.java or MainActivity.kt:"

	// google-services.json
	MsgGoogleJsonMissing = "google-services.json is missing from app/google-services.json."
	MsgGoogleJsonGuideLink = "See https://docs.clix.so/firebase-setting for setup instructions."
	MsgGoogleJsonFound = "google-services.json is present."

	// General
	MsgSourceDirNotFound = "Source directory is missing."
	MsgWorkingDirectoryNotFound = "Working directory is missing."
	MsgAppBuildGradleNotFound = "app/build.gradle(.kts) is missing."
	MsgAppBuildGradleReadFail = "Unable to read app/build.gradle(.kts)."

	// Titles for checks
	TitleGradleRepoCheck = "Checking Gradle repository configuration..."
	TitleClixDependencyCheck = "Checking Clix SDK dependency..."
	TitleGmsPluginCheck = "Checking for Google Services plugin..."
	TitleClixInitializationCheck = "Checking Clix SDK initialization..."
	TitlePermissionCheck = "Checking permission request implementation..."
	TitleGoogleServicesJsonCheck = "Checking for google-services.json..."

	// Fix instructions
	FixGradleRepo = "To resolve this, add the following to settings.gradle(.kts) or build.gradle(.kts):"
	FixClixDependency = "To resolve this, add the following to app/build.gradle(.kts):"
	FixGmsPlugin = "To resolve this, add the following to build.gradle(.kts):"
	FixPermissionRequest = "To resolve this, add the following to MainActivity.java or MainActivity.kt:"

	// Auto-fix messages
	MsgAutoFixSuccess = "Fixed: Changes applied automatically."

	// Code snippets for fixes
	CodeGradleRepo = `repositories {
	mavenCentral()
}`
	CodeClixDependency = `dependencies {
	implementation("so.clix:clix-android-sdk:0.0.2")
}`
	CodeGmsPlugin = `plugins {
	id("com.google.gms.google-services") version "4.4.2"
}`
)
