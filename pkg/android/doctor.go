package android

import (
	"github.com/clix-so/clix-cli/pkg/logx"
)

// StringContainsClixInitializeInOnCreate checks if Clix.initialize(this, ...) is called inside onCreate
func StringContainsClixInitializeInOnCreate(content string) bool {
	// Find signature variants
	onCreateIdx := -1
	if idx := IndexOf(content, "void onCreate"); idx != -1 { // Java
		onCreateIdx = idx
	} else if idx := IndexOf(content, "fun onCreate"); idx != -1 { // Kotlin
		onCreateIdx = idx
	}
	if onCreateIdx == -1 {
		return false
	}

	// Find opening brace after signature (skip annotations / modifiers already inside substring)
	openIdx := -1
	for i := onCreateIdx; i < len(content); i++ {
		c := content[i]
		if c == '{' {
			openIdx = i
			break
		}
		if c == ';' { // Not a method definition
			return false
		}
	}
	if openIdx == -1 {
		return false
	}

	// Extract full block by brace depth tracking
	depth := 0
	endIdx := -1
	for i := openIdx; i < len(content); i++ {
		ch := content[i]
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				endIdx = i + 1
				break
			}
		}
	}
	if endIdx == -1 { // Malformed braces fallback (limit 500 chars)
		endIdx = openIdx + 500
		if endIdx > len(content) {
			endIdx = len(content)
		}
	}

	block := content[openIdx:endIdx]
	norm := removeAllWhitespace(block)

	// Basic pattern (no whitespace) e.g., Clix.initialize(this
	if Contains(norm, "Clix.initialize(this") {
		return true
	}
	// Allow generics: Clix.initialize<...>(this
	if Contains(norm, "Clix.initialize<") && Contains(norm, "(this") {
		return true
	}
	return false
}

// removeAllWhitespace strips common whitespace to allow cross-line pattern detection
func removeAllWhitespace(s string) string {
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			buf = append(buf, s[i])
		}
	}
	return string(buf)
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
