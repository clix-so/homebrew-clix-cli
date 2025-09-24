package android

import (
	"fmt"
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
	logx.Separatorln()
	logx.Log().Title().Println("Starting Clix SDK doctor for Android…")
	logx.Separatorln()

	// Suppress internal logs during checks to control output style
	logx.Mute()

	// 1/6 Gradle repository
	fmt.Print("[1/6] Checking Gradle repositories... ")
	repoOK := CheckGradleRepository(projectRoot)
	if repoOK {
		fmt.Println("OK")
	} else {
		fmt.Println("❌")
		logx.Unmute()
		logx.Log().Branch().Println(logx.FixGradleRepo)
		logx.Log().Indent(3).Code().Println(logx.CodeGradleRepo)
		logx.Mute()
	}

	// 2/6 Clix SDK dependency
	fmt.Print("[2/6] Checking Clix SDK dependency... ")
	depOK := CheckGradleDependency(projectRoot)
	if depOK {
		fmt.Println("OK")
	} else {
		fmt.Println("❌")
		logx.Unmute()
		logx.Log().Branch().Println(logx.FixClixDependency)
		logx.Log().Indent(3).Code().Println(logx.CodeClixDependency)
		logx.Mute()
	}

	// 3/6 Google Services plugin
	fmt.Print("[3/6] Checking Google Services plugin... ")
	gmsOK := CheckGradlePlugin(projectRoot)
	if gmsOK {
		fmt.Println("OK")
	} else {
		fmt.Println("❌")
		logx.Unmute()
		logx.Log().Branch().Println(logx.FixGmsPlugin)
		logx.Log().Indent(3).Code().Println(logx.CodeGmsPlugin)
		logx.Mute()
	}

	// 4/6 Clix initialization
	fmt.Print("[4/6] Checking Clix SDK initialization... ")
	initOK, _ := CheckClixCoreImport(projectRoot)
	if initOK {
		fmt.Println("OK")
	} else {
		fmt.Println()
		logx.Unmute()
		logx.Log().Branch().Println(logx.FixClixInitialization)
		logx.Log().Indent(3).Code().Println(logx.ClixInitializationLink)
		logx.Mute()
	}

	// 5/6 Permission request in MainActivity
	fmt.Print("[5/6] Checking permission request implementation... ")
	permOK := CheckAndroidMainActivityPermissions(projectRoot)
	if permOK {
		fmt.Println("OK")
	} else {
		fmt.Println()
		logx.Unmute()
		logx.Log().Branch().Println(logx.FixPermissionRequest)
		logx.Log().Indent(3).Code().Println(logx.PermissionRequestLink)
		logx.Mute()
	}

	// 6/6 google-services.json
	fmt.Print("[6/6] Checking google-services.json... ")
	gsOK := CheckGoogleServicesJSON(projectRoot)
	if gsOK {
		fmt.Println("OK")
	} else {
		fmt.Println()
		logx.Unmute()
		logx.Log().Branch().Println(logx.FixGoogleServicesJson)
		logx.Log().Indent(3).Code().Println(logx.GoogleServicesJsonLink)
		logx.Mute()
	}

	// Restore logging
	logx.Unmute()

	// Summary
	logx.Separatorln()
	if repoOK && depOK && gmsOK && initOK && permOK && gsOK {
		fmt.Println("🎉 Your Android project is properly configured for Clix SDK!")
		fmt.Println("  └ Push notifications should be working correctly.")
	} else {
		fmt.Println("⚠️ Some issues were found with your Clix SDK integration.")
		fmt.Println("  └ Please fix the issues mentioned above.")
		fmt.Println("  └ Run 'clix-cli install --android' to fix most issues automatically.")
	}
	logx.Separatorln()
}
