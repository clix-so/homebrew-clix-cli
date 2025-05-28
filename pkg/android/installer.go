package android

import (
	"fmt"
	"os"
	"path/filepath"
)


// HandleAndroidInstall guides the user through the Android installation checklist.
func HandleAndroidInstall(apiKey, projectID string) {
	fmt.Println("🤖 Installing Clix SDK for Android...")

	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("[ERROR] Could not determine working directory.")
		return
	}

	// 1. Check google-services.json
	gsPath := filepath.Join(projectRoot, "app", "google-services.json")
	if _, err := os.Stat(gsPath); os.IsNotExist(err) {
		fmt.Println("\n❗ google-services.json not found.")
		fmt.Println("Please download google-services.json from your Firebase Console and place it in app/google-services.json.")
		fmt.Println("Firebase Console: https://console.firebase.google.com/")
		return
	}

	// 2. Check Gradle repository and dependency
	repoOK := CheckGradleRepository(projectRoot)
	if !repoOK {
		if AddGradleRepository(projectRoot) {
			fmt.Println("[AUTO-FIX] Added 'repositories { mavenCentral() }' to your Gradle config.")
			repoOK = true
		} else {
			fmt.Println("[FAIL] Could not automatically add mavenCentral() to your Gradle config. Please add it manually.")
		}
	}
	depOK := CheckGradleDependency(projectRoot)
	if !depOK {
		if AddGradleDependency(projectRoot) {
			fmt.Println("[AUTO-FIX] Added 'implementation(\"so.clix:clix-android-sdk:1.0.0\")' to your app/build.gradle(.kts).")
			depOK = true
		} else {
			fmt.Println("[FAIL] Could not automatically add Clix SDK dependency. Please add it manually.")
		}
	}

	// 3. Check Application class for Clix import and initialization
	appOK := CheckAndroidApplicationSetup(projectRoot)
	if !appOK {
		if AddClixInitializationToApplication(projectRoot, apiKey, projectID) {
			fmt.Println("[AUTO-FIX] Added Clix SDK initialization code to your Application class.")
			appOK = true
		} else {
			fmt.Println("[FAIL] Could not automatically add Clix SDK initialization. Please add it manually to your Application class.")
		}
	}
	// 4. Check MainActivity for permission request
	mainActivityOK := CheckAndroidMainActivityPermissions(projectRoot)

	if repoOK && depOK && appOK && mainActivityOK {
		fmt.Println("\n✅ Clix SDK installation checklist complete! Your Android project is ready.")
	} else {
		fmt.Println("\n❗ Please address the above issues and re-run 'clix install --android' or 'clix doctor --android'.")
	}
}

// CheckAndroidApplicationSetup checks Application class for import and initialization, prints instructions if missing
func CheckAndroidApplicationSetup(projectRoot string) bool {
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
	appFiles := []string{}
	findAppFiles := func(root string) {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (filepath.Ext(path) == ".java" || filepath.Ext(path) == ".kt") &&
				len(info.Name()) >= len("Application.java") &&
				info.Name()[len(info.Name())-len("Application.java"):] == "Application.java" ||
				info.Name()[len(info.Name())-len("Application.kt"):] == "Application.kt" {
				appFiles = append(appFiles, path)
			}
			return nil
		})
	}
	findAppFiles(javaDir)
	findAppFiles(kotlinDir)
	if len(appFiles) == 0 {
		fmt.Println("\n❗ No Application class found. Please create one and register it in your AndroidManifest.xml.")
		fmt.Println("Example (Kotlin):\n---------------------")
		fmt.Println(`class MyApp : Application() {\n    override fun onCreate() {\n        super.onCreate()\n        Clix.initialize(this, /* config */)\n    }\n}`)
		return false
	}
	importFound := false
	initFound := false
	for _, file := range appFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := string(data)
		if Contains(content, "import so.clix.core.Clix") {
			importFound = true
		}
		if StringContainsClixInitializeInOnCreate(content) {
			initFound = true
		}
	}
	if !importFound {
		fmt.Println("\n❗ Application class is missing 'import so.clix.core.Clix'. Please add this import.")
	}
	if !initFound {
		fmt.Println("\n❗ Application class is missing Clix.initialize(this, ...) in onCreate(). Please add:")
		fmt.Println("    Clix.initialize(this, /* config */)")
	}
	return importFound && initFound
}

// CheckAndroidMainActivityPermissions checks MainActivity for permission request code, prints instructions if missing
func CheckAndroidMainActivityPermissions(projectRoot string) bool {
	mainActivityFiles := []string{}
	javaDir := filepath.Join(projectRoot, "app", "src", "main", "java")
	kotlinDir := filepath.Join(projectRoot, "app", "src", "main", "kotlin")
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
		fmt.Println("\n❗ No MainActivity.java or MainActivity.kt found. Please ensure you have a MainActivity.")
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
	if !found {
		fmt.Println("\n❗ MainActivity is missing code to request permissions. Example:")
		fmt.Println(`ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.POST_NOTIFICATIONS), 1001)`)
	}
	return found
}

