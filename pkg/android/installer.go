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
		logx.Log().Indent(6).Code().Println(logx.CodeClixDependency)
	}
	fmt.Println()

	logx.Log().WithSpinner().Println(logx.TitleGmsPluginCheck)
	pluginOK := CheckGradlePlugin(projectRoot)
	if !pluginOK {
		if AddGradlePlugin(projectRoot) {
			logx.Log().Branch().Success().Println(logx.MsgAutoFixSuccess)
			pluginOK = true
		} else {
			logx.Log().Branch().Failure().Println(logx.MsgGmsPluginFixFailure)
		}
		logx.Log().Indent(6).Code().Println(logx.CodeGmsPlugin)
	}
	fmt.Println()

	logx.Log().WithSpinner().Println(logx.TitleClixInitializationCheck)
	appOK, code := CheckClixCoreImport(projectRoot)
	if !appOK {
		if code == "missing-application" {
			ok, _ := AddApplication(projectRoot, apiKey, projectID)
			if ok {
				logx.Log().Branch().Success().Println(logx.MsgAppCreateSuccess)
				appOK = true
			} else {
				logx.Log().Branch().Failure().Println(logx.MsgAppCannotFix)
			}
		} else if code == "missing-content" {
			ok := AddClixInitializationToApplication(projectRoot, apiKey, projectID)
			if ok {
				logx.Log().Branch().Success().Println(logx.MsgAppInitFixSuccess)
				appOK = true
			} else {
				logx.Log().Branch().Failure().Println(logx.MsgAppInitFixFailure)
			}
		} else {
			logx.Log().Branch().Failure().Println(logx.MsgAppCannotFix)
		    logx.Log().Indent(6).Println(logx.MsgAppManualGuideLink)
		}
	}
	fmt.Println()

	logx.Log().WithSpinner().Println(logx.TitlePermissionCheck)
	mainActivityOK := CheckAndroidMainActivityPermissions(projectRoot)
	if !mainActivityOK {
		logx.Log().Branch().Failure().Println(logx.MsgPermissionMissing)
		fmt.Println()
		logx.Log().Indent(6).Code().Println(logx.CodePermissionRequest)
	}
	fmt.Println()

	if repoOK && depOK && appOK && mainActivityOK {
		logx.Log().Success().Println("Clix SDK installation checklist complete! Your Android project is ready.")
	} else {
		logx.Log().Failure().Println("Please address the above issues and re-run 'clix install --android' or 'clix doctor --android'.")
	}
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
			newContent := content[:insertAt] + "\n    implementation(\"so.clix:clix-android-sdk:0.0.2\")" + content[insertAt:]
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

// AddClixInitializationToApplication inserts Clix SDK initialization code into the Application.kt if missing
// FIXME(@nyanxyz): Does not work properly yet
func AddClixInitializationToApplication(projectRoot, apiKey, projectID string) bool {
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	found := false
	err := filepath.Walk(kotlinDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), "Application.kt") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil
			}
			content := string(data)
			if strings.Contains(content, "Clix.initialize(") {
				found = true
				return nil
			}
			// Insert imports at the top if missing
			importBlock := "import so.clix.Clix\nimport so.clix.ClixConfig\nimport so.clix.ClixLogLevel\n"
			if !strings.Contains(content, "import so.clix.Clix") {
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.HasPrefix(line, "package ") {
						// Insert after package
						lines = append(lines[:i+1], append([]string{importBlock}, lines[i+1:]...)...)
						break
					}
				}
				content = strings.Join(lines, "\n")
			}
			// Insert initialization in onCreate
			initBlock := `override fun onCreate() {
        super.onCreate()
        // Project ID: ` + projectID + `
        lifecycleScope.launch {
            try {
                val config =
            ClixConfig(
                projectId = "` + projectID + `",
                apiKey = "` + apiKey + `",
            )
        Clix.initialize(this, config)
            } catch (e: Exception) {
                // Handle initialization failure
            }
        }
    }`
			if strings.Contains(content, "override fun onCreate()") {
				// Replace existing onCreate with template
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "override fun onCreate()") {
						// Replace block (simple heuristic: next 2~20 lines)
						end := i + 1
						for ; end < len(lines) && end-i < 20; end++ {
							if strings.Contains(lines[end], "}") {
								break
							}
						}
						lines = append(lines[:i], append([]string{initBlock}, lines[end+1:]...)...)
						break
					}
				}
				content = strings.Join(lines, "\n")
			} else {
				// Insert initBlock before last '}'
				idx := strings.LastIndex(content, "}")
				if idx != -1 {
					content = content[:idx] + initBlock + "\n}" + content[idx+1:]
				}
			}
			ioutil.WriteFile(path, []byte(content), 0644)
			found = true
		}
		return nil
	})
	return found && err == nil
}
