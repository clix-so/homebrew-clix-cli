package ios

import (
	"bytes"
	_ "embed" // Required for go:embed
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed scripts/configure_xcode_project.rb
var configureXcodeProjectRbScript string

// AutomationResult represents the JSON result returned by the Ruby script
type AutomationResult struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// EnsureRuby checks if Ruby is installed and installs it if needed
func EnsureRuby() error {
	_, err := exec.LookPath("ruby")
	if err != nil {
		fmt.Println("⏳ Ruby is required but not found. Installing Ruby...")

		// Check if homebrew is installed
		_, err := exec.LookPath("brew")
		if err != nil {
			return fmt.Errorf("ruby and homebrew are both not installed, please install ruby manually")
		}

		// Install Ruby using Homebrew
		cmd := exec.Command("brew", "install", "ruby")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Ruby: %v", err)
		}

		fmt.Println("✅ Ruby installed successfully")
	}

	return nil
}

// EnsureXcodeproj ensures the xcodeproj gem is installed
func EnsureXcodeproj() error {
	// Check if xcodeproj is already installed
	cmd := exec.Command("gem", "list", "-i", "xcodeproj")
	output, err := cmd.CombinedOutput()
	if err == nil && strings.TrimSpace(string(output)) == "true" {
		return nil // Gem is already installed
	}

	fmt.Println("⏳ Installing xcodeproj gem...")

	// Try to install without sudo first
	cmd = exec.Command("gem", "install", "xcodeproj")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	// If installation fails, try with sudo
	if err != nil {
		fmt.Println("🔐 Regular gem installation failed. Trying with sudo...")
		fmt.Println("You may be prompted for your password.")

		cmd = exec.Command("sudo", "gem", "install", "xcodeproj")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install xcodeproj gem: %v", err)
		}
	}

	fmt.Println("✅ xcodeproj gem installed successfully")
	return nil
}

// FindXcodeProject attempts to find the Xcode project file in the current directory
// and returns its path if found
func FindXcodeProject() (string, error) {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Find .xcodeproj files
	matches, err := filepath.Glob(filepath.Join(cwd, "*.xcodeproj"))
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		// Look one level up
		matches, err = filepath.Glob(filepath.Join(cwd, "..", "*.xcodeproj"))
		if err != nil {
			return "", err
		}
	}

	if len(matches) == 0 {
		// Look in iOS directory if exists
		iosDir := filepath.Join(cwd, "ios")
		if _, err := os.Stat(iosDir); err == nil {
			matches, err = filepath.Glob(filepath.Join(iosDir, "*.xcodeproj"))
			if err != nil {
				return "", err
			}
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no .xcodeproj file found")
	}

	// Return the first .xcodeproj file found
	return matches[0], nil
}

// ConfigureXcodeProject sets up App Groups and framework dependencies in the Xcode project
func ConfigureXcodeProject(projectID string, verbose bool, dryRun bool) error {
	// Ensure Ruby is installed
	if err := EnsureRuby(); err != nil {
		return err
	}

	// Ensure xcodeproj gem is installed
	if err := EnsureXcodeproj(); err != nil {
		return err
	}

	// Find Xcode project
	xcodeProjectPath, err := FindXcodeProject()
	if err != nil {
		return fmt.Errorf("failed to locate Xcode project: %v", err)
	}

	fmt.Printf("📁 Found Xcode project: %s\n", xcodeProjectPath)

	// Create app group ID
	appGroupID := fmt.Sprintf("group.clix.%s", projectID)

	// The Ruby script is embedded in the binary
	// Create a temporary file to hold the script
	tmpfile, err := os.CreateTemp("", "configure_xcode_project*.rb")
	if err != nil {
		return fmt.Errorf("failed to create temporary script file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write the embedded script to the temporary file
	if _, err := tmpfile.Write([]byte(configureXcodeProjectRbScript)); err != nil {
		tmpfile.Close()
		return fmt.Errorf("failed to write script to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %v", err)
	}

	// Set the script path to the temporary file
	scriptPath := tmpfile.Name()

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("ruby script not found at: %s", scriptPath)
	}

	// Build command arguments
	args := []string{
		scriptPath,
		"--project-path", xcodeProjectPath,
		"--app-group-id", appGroupID,
	}

	if verbose {
		args = append(args, "--verbose")
	}

	if dryRun {
		// In dry-run mode, just show what would be done
		fmt.Println("🔍 DRY RUN - Would execute:")
		fmt.Printf("   ruby %s\n", strings.Join(args, " "))
		fmt.Printf("   This would configure App Group '%s'\n", appGroupID)
		fmt.Println("   This would add Clix framework to NotificationServiceExtension")
		return nil
	}

	fmt.Println("🔄 Configuring Xcode project...")

	// Run the Ruby script
	cmd := exec.Command("ruby", args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	stdout := stdoutBuf.Bytes()
	stderr := stderrBuf.String()

	if err != nil {
		// Even if the script fails, it might have printed a JSON error message to stdout.
		var result AutomationResult
		if jsonErr := json.Unmarshal(stdout, &result); jsonErr == nil && !result.Success {
			// We got a structured error message from the script. Use it.
			return fmt.Errorf("%s", result.Message)
		} else if jsonErr != nil {
			// Diagnostic print for JSON parsing failure
			fmt.Fprintf(os.Stderr, "DEBUG: Failed to parse JSON from script stdout during error handling.\nJSON Error: %v\nStdout: %s\n", jsonErr, string(stdout))
		}

		// If we couldn't get a structured error, return a generic one with any stderr output.
		if len(stderr) > 0 {
			return fmt.Errorf("failed to run Ruby script: %v\n\n--- Script Error ---\n%s", err, stderr)
		}
		return fmt.Errorf("failed to run Ruby script: %v", err)
	}

	// Parse JSON result from successful execution
	var result AutomationResult
	if err := json.Unmarshal(stdout, &result); err != nil {
		return fmt.Errorf("failed to parse script output: %v\nOutput: %s", err, string(stdout))
	}

	if !result.Success {
		// This case handles when the script exits 0 but reports failure in JSON.
		fmt.Fprintf(os.Stderr, "Error from configuration script: %s\n", result.Message)
		return fmt.Errorf("script reported failure: %s", result.Message)
	}

	fmt.Println("✅ Xcode project configured successfully!")
	fmt.Println("   - App Groups capability added")
	fmt.Println("   - Clix framework added to NotificationServiceExtension")

	return nil
}
