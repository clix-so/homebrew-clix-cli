package android

import (
	"github.com/clix-so/clix-cli/pkg/logx"
)

// stringContainsImportClix checks if the given file content contains the import statement for so.clix.core.Clix
func stringContainsImportClix(content string) bool {
	return len(content) > 0 && Contains(content, "import so.clix.core.Clix")
}

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
		logx.Log().Indent(3).Println(logx.FixGradleRepo)
		logx.NewLine()
		logx.Log().Indent(3).Code().Println(logx.CodeGradleRepo)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleClixDependencyCheck)
	if !CheckGradleDependency(projectRoot) {
		logx.Log().Indent(3).Println(logx.FixClixDependency)
		logx.NewLine()
		logx.Log().Indent(3).Code().Println(logx.CodeClixDependency)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleGmsPluginCheck)
	if !CheckGradlePlugin(projectRoot) {
		logx.Log().Indent(3).Println(logx.FixGmsPlugin)
		logx.NewLine()
		logx.Log().Indent(3).Code().Println(logx.CodeGmsPlugin)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleClixInitializationCheck)
	CheckClixCoreImport(projectRoot)
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitlePermissionCheck)
	if !CheckAndroidMainActivityPermissions(projectRoot) {
		logx.Log().Indent(3).Println(logx.FixPermissionRequest)
		logx.NewLine()
		logx.Log().Indent(3).Code().Println(logx.CodePermissionRequest)
	}
	logx.NewLine()

	logx.Log().WithSpinner().Title().Println(logx.TitleGoogleServicesJsonCheck)
	CheckGoogleServicesJSON(projectRoot)
	logx.NewLine()
}
