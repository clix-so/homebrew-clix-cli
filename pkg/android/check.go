package android

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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

		content := stripGradleComments(string(data))
		if hasMavenCentral(content) {
			found = true
			break
		}
	}

	// As a fallback, scan application module's build.gradle(.kts) in case repositories are declared there
	if !found {
		appBuildGradle := findApplicationModuleBuildGradle(projectRoot)
		if appBuildGradle != "" {
			if b, err := ioutil.ReadFile(appBuildGradle); err == nil {
				if hasMavenCentral(stripGradleComments(string(b))) {
					found = true
				}
			}
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

	// Fallback: try to locate an Android application module if app module isn't found
	if appBuildGradleFilePath == "" {
		appBuildGradleFilePath = findApplicationModuleBuildGradle(projectRoot)
	}

	if appBuildGradleFilePath == "" {
		logx.Log().Failure().Println(logx.MsgAppBuildGradleNotFound)
		return false
	}

	data, err := ioutil.ReadFile(appBuildGradleFilePath)
	if err != nil {
		logx.Log().Failure().Println(logx.MsgAppBuildGradleReadFail)
		return false
	}

	content := string(data)
	noComments := stripGradleComments(content)

	if containsClixDependency(noComments) || containsClixVersionCatalogAlias(noComments) {
		logx.Log().Success().Println(logx.MsgClixDependencySuccess)
		return true
	}

	logx.Log().Failure().Println(logx.MsgClixDependencyFailure)
	return false
}

// stripGradleComments removes // and /* */ comments in Groovy/KTS files
func stripGradleComments(s string) string {
	// Remove block comments first
	reBlock := regexp.MustCompile(`(?s)/\*.*?\*/`)
	s = reBlock.ReplaceAllString(s, "")
	// Remove line comments
	reLine := regexp.MustCompile(`(?m)^\s*//.*$`)
	s = reLine.ReplaceAllString(s, "")
	return s
}

// containsClixDependency matches various Gradle notations for the dependency, with or without version
func containsClixDependency(s string) bool {
	// e.g., implementation("so.clix:clix-android-sdk"), api('so.clix:clix-android-sdk:1.2.3')
	reDirect := regexp.MustCompile(`(?m)^\s*(implementation|api|compileOnly|runtimeOnly|debugImplementation|releaseImplementation)\s*\(\s*['"]so\.clix:clix-android-sdk(?:[:][^'"\)]+)?['"]\s*\)`)
	// e.g., add("implementation", "so.clix:clix-android-sdk")
	reAdd := regexp.MustCompile(`(?m)^\s*add\s*\(\s*['"](implementation|api|compileOnly|runtimeOnly|debugImplementation|releaseImplementation)['"]\s*,\s*['"]so\.clix:clix-android-sdk(?:[:][^'"\)]+)?['"]\s*\)`)

	if reDirect.MatchString(s) || reAdd.MatchString(s) {
		return true
	}
	return false
}

// containsClixVersionCatalogAlias detects version catalog aliases like libs.clix.android.sdk
func containsClixVersionCatalogAlias(s string) bool {
	reLibs := regexp.MustCompile(`(?m)^\s*(implementation|api|debugImplementation|releaseImplementation)\s*\(\s*libs\.[A-Za-z0-9_.-]*clix[A-Za-z0-9_.-]*\s*\)`)
	return reLibs.MatchString(s)
}

// findApplicationModuleBuildGradle walks the project to find a module using the Android application plugin
func findApplicationModuleBuildGradle(root string) string {
	var result string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == "build" || name == ".git" || name == ".gradle" || name == "gradle" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, "build.gradle") || strings.HasSuffix(path, "build.gradle.kts") {
			b, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			txt := string(b)
			// Heuristics: presence of com.android.application plugin
			if strings.Contains(txt, "com.android.application") {
				result = path
				return fs.SkipAll
			}
		}
		return nil
	})
	return result
}

// CheckGradlePlugin checks if com.google.gms:google-services is present in app/build.gradle(.kts)
func CheckGradlePlugin(projectRoot string) bool {
	// Prefer application module
	moduleGradle := GetAppBuildGradlePath(projectRoot)
	if moduleGradle == "" {
		moduleGradle = findApplicationModuleBuildGradle(projectRoot)
	}

	// Check in the application module first
	if moduleGradle != "" {
		if b, err := ioutil.ReadFile(moduleGradle); err == nil {
			c := stripGradleComments(string(b))
			if containsGoogleServicesPlugin(c) {
				logx.Log().Success().Println(logx.MsgGmsPluginFound)
				return true
			}
		}
	}

	// Fallback: some projects apply plugin from root to subprojects
	rootGradles := []string{
		filepath.Join(projectRoot, "build.gradle"),
		filepath.Join(projectRoot, "build.gradle.kts"),
	}
	for _, f := range rootGradles {
		if b, err := ioutil.ReadFile(f); err == nil {
			c := stripGradleComments(string(b))
			if containsGoogleServicesPluginInSubprojects(c) {
				logx.Log().Success().Println(logx.MsgGmsPluginFound)
				return true
			}
		}
	}

	logx.Log().Failure().Println(logx.MsgGmsPluginNotFound)
	return false
}

// hasMavenCentral detects mavenCentral() or maven { url "https://repo.maven.apache.org/maven2" } within any repositories block
func hasMavenCentral(s string) bool {
	// repositories { ... mavenCentral() ... }
	reRepoBlock := regexp.MustCompile(`(?s)repositories\s*\{.*?\}`)
	blocks := reRepoBlock.FindAllString(s, -1)
	for _, b := range blocks {
		if strings.Contains(b, "mavenCentral()") {
			return true
		}
		// maven { url "https://repo.maven.apache.org/maven2" } or repo1
		reMavenUrl := regexp.MustCompile(`maven\s*\{[^}]*url[^\n\r\"]*[\"']https?://repo(\.maven|\d+)\.apache\.org/maven2[\"'][^}]*\}`)
		if reMavenUrl.MatchString(b) {
			return true
		}
	}
	// Also handle settings.gradle(.kts) dependencyResolutionManagement { repositories { ... } }
	if strings.Contains(s, "dependencyResolutionManagement") && strings.Contains(s, "mavenCentral()") {
		return true
	}
	if strings.Contains(s, "pluginManagement") && strings.Contains(s, "mavenCentral()") {
		return true
	}
	return false
}

// containsGoogleServicesPlugin detects various ways to apply the Google Services plugin inside a module build.gradle(.kts)
func containsGoogleServicesPlugin(s string) bool {
	// plugins { id("com.google.gms.google-services") }
	reId := regexp.MustCompile(`(?s)plugins\s*\{[^}]*id\s*\(\s*['"]com\.google\.gms\.google-services['"]\s*\)[^}]*\}`)
	if reId.MatchString(s) {
		return true
	}
	// Groovy: apply plugin: 'com.google.gms.google-services'
	reApplyGroovy := regexp.MustCompile(`apply\s+plugin\s*:\s*['"]com\.google\.gms\.google-services['"]`)
	if reApplyGroovy.MatchString(s) {
		return true
	}
	// Kotlin DSL: apply(plugin = "com.google.gms.google-services")
	reApplyKts := regexp.MustCompile(`apply\s*\(\s*plugin\s*=\s*['"]com\.google\.gms\.google-services['"]\s*\)`)
	if reApplyKts.MatchString(s) {
		return true
	}
	// Version catalog alias, broaden to common patterns containing gms/google/services
	reAlias := regexp.MustCompile(`(?s)plugins\s*\{[^}]*alias\s*\(\s*libs\.plugins\.[A-Za-z0-9_.-]*(gms|google)[A-Za-z0-9_.-]*services[A-Za-z0-9_.-]*\s*\)[^}]*\}`)
	return reAlias.MatchString(s)
}

// containsGoogleServicesPluginInSubprojects detects application from root project into subprojects
func containsGoogleServicesPluginInSubprojects(s string) bool {
	// subprojects { apply plugin: 'com.google.gms.google-services' }
	re := regexp.MustCompile(`(?s)subprojects\s*\{[^}]*apply\s+plugin\s*:\s*['"]com\.google\.gms\.google-services['"][^}]*\}`)
	if re.MatchString(s) {
		return true
	}
	// allprojects { plugins { id("com.google.gms.google-services") } } uncommon but possible
	re2 := regexp.MustCompile(`(?s)(subprojects|allprojects)\s*\{[^}]*plugins\s*\{[^}]*id\s*\(\s*['"]com\.google\.gms\.google-services['"]\s*\)[^}]*\}[^}]*\}`)
	return re2.MatchString(s)
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

	ktPath := filepath.Join(sourceDir, appPath+".kt")
	javaPath := filepath.Join(sourceDir, appPath+".java")

	if _, err := os.Stat(javaPath); err == nil {
		appPath = javaPath
	} else if _, err := os.Stat(ktPath); err == nil {
		appPath = ktPath
	} else {
		logx.Log().Failure().Println(logx.MsgApplicationClassMissing)
		return false, "unknown"
	}

	initializeFound := false

	data, err := os.ReadFile(appPath)
	if err != nil {
		logx.Log().Failure().Println(logx.MsgApplicationFileReadFail)
		return false, "unknown"
	}
	content := string(data)
	if StringContainsClixInitializeInOnCreate(content) {
		initializeFound = true
	}

	if initializeFound {
		logx.Log().Success().Println(logx.MsgClixInitSuccess)
		return true, ""
	} else {
		logx.Log().Failure().Println(logx.MsgClixInitMissing)
		return false, "missing-content"
	}
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
		return false
	}

	logx.Log().Success().Println(logx.MsgGoogleJsonFound)
	return true
}
