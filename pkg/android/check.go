package android

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/clix-so/clix-cli/pkg/logx"
)

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
		logx.Log().Success().Println(logx.MsgGradleRepoSuccess)
		return true
	}

	logx.Log().Failure().Println(logx.MsgGradleRepoFailure)
	return false
}

// CheckGradleDependency checks if so.clix:clix-android-sdk is present in app/build.gradle(.kts)
func CheckGradleDependency(projectRoot string) bool {
	appBuildGradleFilePath := GetAppBuildGradlePath(projectRoot)

	if appBuildGradleFilePath == "" {
		logx.Log().Failure().Println(logx.MsgAppBuildGradleNotFound)
		return false
	}

	found := false
	data, err := ioutil.ReadFile(appBuildGradleFilePath)
	if err != nil {
		logx.Log().Failure().Println(logx.MsgAppBuildGradleReadFail)
		return false
	}

	content := string(data)
	if Contains(content, "implementation(\"so.clix:clix-android-sdk:") {
		found = true
	} else if Contains(content, "implementation(libs.clix.android.sdk)") {
		found = true
	}
	

	if (found) {
		logx.Log().Success().Println(logx.MsgClixDependencySuccess)
		return true
	}

	logx.Log().Failure().Println(logx.MsgClixDependencyFailure)
	return false
}

// CheckGradlePlugin checks if com.google.gms:google-services is present in app/build.gradle(.kts)
func CheckGradlePlugin(projectRoot string) bool {
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
		if Contains(content, "alias(libs.plugins.gms") {
			found = true
			break
		}
		if Contains(content, "id(\"com.google.gms.google-services\")") {
			found = true
			break
		}
	}

	if (found) {
		logx.Log().Success().Println(logx.MsgGmsPluginFound)
		return true
	}

	logx.Log().Failure().Println(logx.MsgGmsPluginNotFound)
	return false
}

// CheckClixCoreImport checks if any Application class imports so.clix.core.Clix (Java or Kotlin).
func CheckClixCoreImport(projectRoot string) (bool, string) {
	manifestPath := filepath.Join(projectRoot, "app", "src", "main", "AndroidManifest.xml")
	appName, err := extractApplicationClassName(manifestPath)

	if err != nil {
		logx.Log().Failure().Println(logx.MsgManifestReadFail)
		return false, "unknown"
	}

	if appName == "" {
		logx.Log().Failure().Println(logx.MsgApplicationClassNotDefined)
		return false, "missing-application"
	}

	appPath := strings.TrimPrefix(appName, ".")
	appPath = strings.ReplaceAll(appPath, ".", string(filepath.Separator))

	sourceDir := GetSourceDirPath(projectRoot)
	if sourceDir == "" {
		logx.Log().Failure().Println(logx.MsgSourceDirNotFound)
		return false, "unknown"
	}

	ktPath := filepath.Join(sourceDir, appPath + ".kt")
	javaPath := filepath.Join(sourceDir, appPath + ".java")

	if _, err := os.Stat(javaPath); err == nil {
		appPath = javaPath
	} else if _, err := os.Stat(ktPath); err == nil {
		appPath = ktPath
	} else {
		logx.Log().Failure().Println(logx.MsgApplicationClassMissing)
		return false, "unknown"
	}

	importFound := false
	initializeFound := false

	data, err := os.ReadFile(appPath)
	if err != nil {
		logx.Log().Failure().Println(logx.MsgApplicationFileReadFail)
		return false, "unknown"
	}
	content := string(data)
	if stringContainsImportClix(content) {
		importFound = true
	}
	if StringContainsClixInitializeInOnCreate(content) {
		initializeFound = true
	}

	if importFound {
		logx.Log().Success().Println(logx.MsgClixImportSuccess)
	} else {
		logx.Log().Failure().Println(logx.MsgClixImportMissing)
	}

	if initializeFound {
		logx.Log().Success().Println(logx.MsgClixInitSuccess)
	} else {
		logx.Log().Failure().Println(logx.MsgClixInitMissing)
	}

	if !importFound || !initializeFound {
		return false, "missing-content"
	}

	return true, ""
}

// CheckAndroidMainActivityPermissions checks MainActivity for permission request code, prints instructions if missing
func CheckAndroidMainActivityPermissions(projectRoot string) bool {
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	mainActivityFiles := []string{}

	findMainActivity := func(root string) {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (info.Name() == "MainActivity.java" || info.Name() == "MainActivity.kt") {
				mainActivityFiles = append(mainActivityFiles, path)
			}
			return nil
		})
	}

	findMainActivity(javaDir)
	findMainActivity(kotlinDir)

	if len(mainActivityFiles) == 0 {
		logx.Log().Failure().Println(logx.MsgMainActivityNotFound)
		return false
	}

	permissionPattern := []string{
		"requestPermissions(",
		"ActivityCompat.requestPermissions(",
		"ContextCompat.checkSelfPermission(",
		"Manifest.permission.",
	}

	found := false
	for _, file := range mainActivityFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(data)
		for _, pat := range permissionPattern {
			if Contains(content, pat) {
				found = true
				break
			}
		}
	}

	if found {
		logx.Log().Success().Println(logx.MsgPermissionFound)
		return true
	}

	logx.Log().Failure().Println(logx.MsgPermissionMissing)
	return false
}

// CheckGoogleServicesJSON checks if google-services.json exists in the correct location
func CheckGoogleServicesJSON(projectRoot string) bool {
	gsPath := filepath.Join(projectRoot, "app", "google-services.json")
	if _, err := os.Stat(gsPath); os.IsNotExist(err) {
		logx.Log().Failure().Println(logx.MsgGoogleJsonMissing)
		logx.Log().Indent(3).Println(logx.MsgGoogleJsonGuideLink)
		return false
	}

	logx.Log().Success().Println(logx.MsgGoogleJsonFound)
	return true
}