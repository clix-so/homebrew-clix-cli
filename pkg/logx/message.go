package logx

const (
	// Gradle Repository
	MsgGradleRepoSuccess = "Gradle repositories are properly configured."
	MsgGradleRepoFailure = "Gradle repository settings are missing."
	MsgGradleRepoFixFailure = "Could not fix automatically. Please add the following manually to settings.gradle(.kts) or build.gradle(.kts):"

	// Clix SDK Dependency
	MsgClixDependencySuccess = "Clix SDK dependency found."
	MsgClixDependencyFailure = "Clix SDK dependency is missing."
	MsgAppBuildGradleNotFound = "app/build.gradle(.kts) not found."
	MsgAppBuildGradleReadFail = "Failed to read app/build.gradle(.kts)"
	MsgClixDependencyFixFailure = "Could not fix automatically. Please add the following manually to app/build.gradle(.kts):"

	// Google Services Plugin
	MsgGmsPluginFound    = "Google services plugin found in app/build.gradle(.kts)."
	MsgGmsPluginNotFound = "Google services plugin not found in app/build.gradle(.kts)."
	MsgGmsPluginFixFailure = "Could not fix automatically. Please add the following manually to build.gradle(.kts):"

	// Clix Application Import & Initialization
	MsgManifestReadFail           = "Failed to read AndroidManifest.xml"
	MsgApplicationClassNotDefined = "No Application class found in AndroidManifest.xml"
	MsgApplicationClassMissing    = "Application class not found in expected locations."
	MsgApplicationFileReadFail    = "Failed to read Application class file"
	MsgClixImportSuccess          = "so.clix.core.Clix is imported in Application class."
	MsgClixImportMissing          = "so.clix.core.Clix is NOT imported in Application class."
	MsgClixInitSuccess            = "Clix.initialize(this, ...) is called in onCreate() of Application class."
	MsgClixInitMissing            = "Clix.initialize(this, ...) is NOT called in onCreate() of Application class."
	MsgAppCreateSuccess       = "Fixed: Application class created successfully"
	MsgAppCannotFix    = "Could not fix automatically. Please follow the guide below to set up your Application class:"
	MsgAppManualGuideLink     = "https://docs.clix.so/sdk-quickstart-android#setup-clix-manual-installation"
	MsgAppInitFixSuccess      = "Fixed: Clix SDK initialization added to Application class"
	MsgAppInitFixFailure      = "Could not fix automatically. Please ensure your Application class initializes Clix SDK."

	// MainActivity Permissions
	MsgMainActivityNotFound     = "No MainActivity.java or MainActivity.kt found. Please ensure you have a MainActivity."
	MsgPermissionFound          = "MainActivity contains code requesting permissions."
	MsgPermissionMissing        = "MainActivity does not contain code requesting permissions."
	MsgPermissionFixFailure   = "Could not fix automatically. Please add the following to your MainActivity.kt or MainActivity.java:"

	// google-services.json
	MsgGoogleJsonMissing        = "Missing google-services.json at app/google-services.json"
	MsgGoogleJsonGuideLink      = "See https://docs.clix.so/firebase-setting for setup instructions."
	MsgGoogleJsonFound          = "google-services.json found"

	// General
	MsgSourceDirNotFound        = "Source directory not found."
	MsgWorkingDirectoryNotFound = "Could not determine working directory."

	// Titles for checks
	TitleGradleRepoCheck        = "Checking Gradle repository settings..."
	TitleClixDependencyCheck    = "Checking for Clix SDK dependency..."
	TitleGmsPluginCheck         = "Checking for Google Services plugin..."
	TitleClixInitializationCheck = "Checking Clix SDK initialization..."
	TitlePermissionCheck        = "Checking permission request..."
	TitleGoogleServicesJsonCheck = "Checking google-services.json..."

	// Fix instructions
	FixGradleRepo = "To fix this, add the following to settings.gradle(.kts) or build.gradle(.kts):"
	FixClixDependency = "To fix this, add the following to app/build.gradle(.kts):"
	FixGmsPlugin = "To fix this, add the following to build.gradle(.kts):"
	FixPermissionRequest = "To fix this, add the following to MainActivity.java or MainActivity.kt:"

	// Auto-fix messages
	MsgAutoFixSuccess = "Fixed: Automatically added"

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
	CodePermissionRequest = `ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.POST_NOTIFICATIONS), 1001)`
)
