package android

import (
	"github.com/clix-so/clix-cli/pkg/logx"
)

// StringContainsClixInitializeInOnCreate checks if Clix.initialize(this, ...) is called inside onCreate
func StringContainsClixInitializeInOnCreate(content string) bool {
	// Simple heuristic: check for 'void onCreate' or 'fun onCreate', then 'Clix.initialize(this'
	onCreateIdx := -1
	if idx := IndexOf(content, "void onCreate"); idx != -1 {
		onCreateIdx = idx
	} else if idx := IndexOf(content, "fun onCreate"); idx != -1 {
		onCreateIdx = idx
	}
	if onCreateIdx == -1 {
		return false
	}
	// Check 200 chars after onCreate for Clix.initialize(this
	endIdx := onCreateIdx + 200
	if endIdx > len(content) {
		endIdx = len(content)
	}
	return Contains(content[onCreateIdx:endIdx], "Clix.initialize(this")
}

// Contains is a helper for substring check (to avoid strings package import)
func Contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && (func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})()
}

// IndexOf returns the index of the first instance of substr in s, or -1 if not present
func IndexOf(s, substr string) int {
	if len(substr) == 0 || len(s) < len(substr) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}



// RunDoctor runs all Android doctor checks.
func RunDoctor(projectRoot string) {
	logx.Log().WithSpinner().Title().Println(logx.TitleGradleRepoCheck)
	if !CheckGradleRepository(projectRoot) {
		logx.Log().Branch().Println(logx.FixGradleRepo)
		logx.NewLine()
		logx.Log().Indent(3).Code().Println(logx.CodeGradleRepo)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleClixDependencyCheck)
	if !CheckGradleDependency(projectRoot) {
		logx.Log().Branch().Println(logx.FixClixDependency)
		logx.NewLine()
		logx.Log().Indent(3).Code().Println(logx.CodeClixDependency)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleGmsPluginCheck)
	if !CheckGradlePlugin(projectRoot) {
		logx.Log().Branch().Println(logx.FixGmsPlugin)
		logx.NewLine()
		logx.Log().Indent(3).Code().Println(logx.CodeGmsPlugin)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleClixInitializationCheck)
	ok, _ := CheckClixCoreImport(projectRoot)
	if !ok {
		logx.Log().Branch().Println(logx.FixClixInitialization)
		logx.Log().Indent(3).Code().Println(logx.ClixInitializationLink)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitlePermissionCheck)
	if !CheckAndroidMainActivityPermissions(projectRoot) {
		logx.Log().Branch().Println(logx.FixPermissionRequest)
		logx.Log().Indent(3).Code().Println(logx.PermissionRequestLink)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleGoogleServicesJsonCheck)
	if !CheckGoogleServicesJSON(projectRoot) {
		logx.Log().Branch().Println(logx.FixGoogleServicesJson)
		logx.Log().Indent(3).Code().Println(logx.GoogleServicesJsonLink)
	}
	logx.NewLine()
}
