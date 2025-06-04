package android

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
)

// HandleAndroidInstall guides the user through the Android installation checklist.
func HandleAndroidInstall(apiKey, projectID string) {
	projectRoot, err := os.Getwd()
	if err != nil {
		logx.Log().Failure().Println(logx.MsgWorkingDirectoryNotFound)
		return
	}

	logx.Log().WithSpinner().Title().Println(logx.TitleGoogleServicesJsonCheck)
	if !CheckGoogleServicesJSON(projectRoot) {
		logx.Log().Branch().Failure().Println(logx.MsgGoogleJsonFixFailure)
		logx.Log().Indent(6).Code().Println(logx.GoogleServicesJsonLink)
		return
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleGradleRepoCheck)
	repoOK := CheckGradleRepository(projectRoot)
	if !repoOK {
		if AddGradleRepository(projectRoot) {
			logx.Log().Branch().Success().Println(logx.MsgAutoFixSuccess)
		} else {
			logx.Log().Branch().Failure().Println(logx.MsgGradleRepoFixFailure)
		}
		logx.NewLine()
		logx.Log().Indent(6).Code().Println(logx.CodeGradleRepo)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Println(logx.TitleClixDependencyCheck)
	depOK := CheckGradleDependency(projectRoot)
	if !depOK {
		if AddGradleDependency(projectRoot) {
			logx.Log().Branch().Success().Println(logx.MsgAutoFixSuccess)
			depOK = true
		} else {
			logx.Log().Branch().Failure().Println(logx.MsgClixDependencyFixFailure)
		}
		logx.NewLine()
		logx.Log().Indent(6).Code().Println(logx.CodeClixDependency)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Println(logx.TitleGmsPluginCheck)
	pluginOK := CheckGradlePlugin(projectRoot)
	if !pluginOK {
		if AddGradlePlugin(projectRoot) {
			logx.Log().Branch().Success().Println(logx.MsgAutoFixSuccess)
			pluginOK = true
		} else {
			logx.Log().Branch().Failure().Println(logx.MsgGmsPluginFixFailure)
		}
		logx.NewLine()
		logx.Log().Indent(6).Code().Println(logx.CodeGmsPlugin)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Println(logx.TitleClixInitializationCheck)
	appOK, code := CheckClixCoreImport(projectRoot)
	if !appOK {
		if code == "missing-application" {
			ok, _ := AddApplication(projectRoot, apiKey, projectID)
			if ok {
				appOK = true
			}
		}
	}
	if appOK {
		logx.Log().Branch().Success().Println(logx.MsgAppCreateSuccess)
	} else {
		logx.Log().Branch().Failure().Println(logx.MsgAppFixFailure)
		logx.Log().Indent(6).Code().Println(logx.ClixInitializationLink)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Println(logx.TitlePermissionCheck)
	mainActivityOK := CheckAndroidMainActivityPermissions(projectRoot)
	if !mainActivityOK {
		logx.Log().Branch().Failure().Println(logx.MsgPermissionFixFailure)
		logx.Log().Indent(6).Code().Println(logx.PermissionRequestLink)
	}
	logx.NewLine()
}


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
		if Contains(content, "implementation(\"so.clix:clix-android-sdk") {
			return true // already present
		}
		// Try to insert after 'dependencies {' or at end
		if idx := IndexOf(content, "dependencies {"); idx != -1 {
			insertAt := idx + len("dependencies {")
			newContent := content[:insertAt] + "\n    implementation(\"so.clix:clix-android-sdk:0.0.4\")" + content[insertAt:]
			err = ioutil.WriteFile(file, []byte(newContent), 0644)
			return err == nil
		}
	}
	return false
}

// AddGradlePlugin tries to insert the Google services plugin into app/build.gradle(.kts)
func AddGradlePlugin(projectRoot string) bool {
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
		if Contains(content, "id(\"com.google.gms.google-services\")") {
			return true // already present
		}
		// Try to insert after 'dependencies {' or at end
		if idx := IndexOf(content, "plugins {"); idx != -1 {
			insertAt := idx + len("plugins {")
			newContent := content[:insertAt] + "\n    id(\"com.google.gms.google-services\") version \"4.4.2\"" + content[insertAt:]
			err = ioutil.WriteFile(file, []byte(newContent), 0644)
			return err == nil
		}
	}
	return false
}

func AddApplicationFile(projectRoot, apiKey, projectID string) (bool, string) {
	sourceDir := GetSourceDirPath(projectRoot)
	if sourceDir == "" {
		return false, "Could not find source directory for Android project."
	}

	appBuildGradlePath := GetAppBuildGradlePath(projectRoot)
	if appBuildGradlePath == "" {
		return false, "Could not find app/build.gradle(.kts) file."
	}

	packageName := GetPackageName(projectRoot)
	if packageName == "" {
		return false, "Could not extract package name from app/build.gradle(.kts)."
	}

	filePath := filepath.Join(sourceDir, "BasicApplication.kt")
	code := `package %s

import android.app.Application
import so.clix.core.Clix
import so.clix.core.ClixConfig

class BasicApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        Clix.initialize(this, ClixConfig(
			projectId = "%s",
            apiKey = "%s",
        ))
    }
}`
	code = fmt.Sprintf(code, packageName, projectID, apiKey)

	err := os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return false, err.Error()
	}

	return true, "Application class created successfully"
}

func AddApplicationNameToMenifest(projectRoot string) (bool, string) {
	manifestPath := GetAndroidManifestPath(projectRoot)

	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return false, "Failed to read AndroidManifest.xml"
	}
	content := string(data)

	newContent := strings.Replace(content, "<application", `<application android:name=".BasicApplication"`, 1)
	err = ioutil.WriteFile(manifestPath, []byte(newContent), 0644)
	if err != nil {
		return false, "Failed to write AndroidManifest.xml"
	}

	return true, "Application name added to AndroidManifest.xml"
}

func AddApplication(projectRoot, apiKey, projectID string) (bool, string) {
	// Step 1: Create BasicApplication.kt file
	ok, message := AddApplicationFile(projectRoot, apiKey, projectID)
	if !ok {
		return false, message
	}

	// Step 2: Add application name to AndroidManifest.xml
	ok, message = AddApplicationNameToMenifest(projectRoot)
	if !ok {
		return false, message
	}

	return true, "Application class setup complete"
}
