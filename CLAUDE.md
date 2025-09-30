# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Clix CLI is a command-line tool for integrating and managing the Clix SDK (push notifications) in mobile projects. It supports iOS, Android, React Native Expo, and Flutter platforms. **Migrated from Go to TypeScript + React (Ink) + Bun for better developer experience and maintainability.**

## Technology Stack

- **TypeScript** - Type-safe JavaScript
- **React + Ink** - Terminal UI components for interactive CLI
- **Bun** - Fast JavaScript runtime and bundler
- **Commander** - CLI argument parsing and command framework
- **Execa** - Process execution for shell commands
- **xml2js** - XML parsing for Android manifests

## Key Commands

### Building and Testing
```bash
# Install dependencies
bun install

# Build the project
bun run build

# Run in development mode
bun run dev

# Type check
bun run typecheck

# Test specific commands
bun run src/cli.tsx install --help
bun run src/cli.tsx doctor --android
```

### Release Process
The project uses GoReleaser for building and releasing (to be updated for TypeScript):
```bash
# Create a release (handled by GitHub Actions on git tag)
# GoReleaser automatically updates clix.rb Homebrew formula
```

### Development Workflow
```bash
# Test the install command
bun run dev install

# Test the doctor command (checks SDK integration)
bun run dev doctor

# Test the uninstall command
bun run dev uninstall

# Test with specific platforms
bun run dev install --ios
bun run dev install --android
bun run dev doctor --expo
bun run dev doctor --flutter
```

## Architecture

### Command Structure
- **Entry point**: `src/cli.tsx` - Sets up Commander CLI with React Ink rendering
- **Command components**: `src/commands/` directory contains React components:
  - `install.tsx`: SDK installation command with platform flags
  - `doctor.tsx`: Integration verification command
  - `uninstall.tsx`: SDK removal command

### Package Organization

**`src/packages/ios/`** - iOS-specific logic:
- `installer.ts`: Handles SPM/CocoaPods installation
- `xcode-project.ts`: Xcode project file manipulation via Ruby scripts
- `notification-service.ts`: Notification service extension setup
- `doctor.ts`: iOS integration verification
- `uninstaller.ts`: iOS SDK removal
- `firebase-checks.ts`: Firebase integration validation
- `locator.ts`: Find Xcode project files
- `constants.ts`: iOS-specific constants
- `scripts/configure_xcode_project.rb`: Ruby script for Xcode automation

**`src/packages/android/`** - Android-specific logic:
- `installer.ts`: Gradle dependency management
- `manifest-parser.ts`: AndroidManifest.xml parsing (uses xml2js)
- `doctor.ts`: Android integration verification
- `check.ts`: Firebase/Google Services JSON validation
- `uninstaller.ts`: Android SDK removal
- `path.ts`: Path utilities for Android project files
- `package-name.ts`: Extract package name from build.gradle
- `index.ts`: Public API exports

**`src/packages/expo/`** - React Native Expo support:
- `installer.ts`: Expo plugin installation with version-aware MMKV selection
- `doctor.ts`: Expo integration checks

**`src/packages/flutter/`** - Flutter support:
- `installer.ts`: Flutter plugin installation via pubspec.yaml
- `doctor.ts`: Flutter integration checks

**`src/utils/`** - Shared utilities:
- `detectPlatform.ts`: Auto-detects project type (iOS/Android/Expo/Flutter) by checking for characteristic files
- `prompt.ts`: User input handling using Ink TextInput
- `shell.ts`: Shell command execution wrapper around execa

**`src/ui/`** - UI components and styling:
- `Logger.tsx`: React component for structured console output with spinners, colors, and formatting
- `messages.ts`: Centralized message constants

**`src/packages/versions.ts`** - SDK version management

### Platform Detection Logic

The CLI auto-detects platforms using file markers:
- **iOS**: `.xcodeproj`, `.xcworkspace`, `Podfile`, `Package.swift`
- **Android**: `build.gradle`, `settings.gradle`, `AndroidManifest.xml`
- **Expo**: `app.json` + `expo` in `package.json`
- **Flutter**: `pubspec.yaml` with `flutter:` dependency
- Priority order: Flutter > Expo > Native iOS/Android

### Installation Flow

1. Platform detection (if not specified via flags)
2. Prompt for Project ID and API Key (using Ink's interactive TextInput)
3. Platform-specific installation:
   - **iOS**: Guides through SPM/CocoaPods setup, configures App Groups via Ruby script, updates NotificationServiceExtension
   - **Android**: Checks/adds gradle repositories, adds Clix dependency (via version catalog or direct), updates AndroidManifest.xml, adds initialization code
   - **Expo**: Installs dependencies with version-aware MMKV selection, updates app.json, creates Clix initialization file
   - **Flutter**: Installs Firebase CLI, runs flutterfire configure, adds dependencies, updates main.dart
4. Runs `doctor` command post-install to verify setup

### Doctor Command Logic

Runs comprehensive checks for each platform:
- SDK dependency presence
- Required configurations (App Groups, notification permissions, etc.)
- Firebase setup (Android, Expo, Flutter)
- Provides actionable fix instructions when issues found

### React + Ink UI Pattern

Commands are React components that:
1. Use `useEffect` to run async logic on mount
2. Manage state with `useState` for status tracking
3. Render UI based on current state (detecting, installing, error, complete)
4. Use `<Logger>` component for consistent output styling
5. Call `process.exit()` when complete

Example pattern:
```tsx
export const MyCommand: React.FC<Props> = (props) => {
  const [status, setStatus] = useState('detecting');

  useEffect(() => {
    (async () => {
      // Async logic here
      setStatus('complete');
      process.exit(0);
    })();
  }, []);

  if (status === 'detecting') {
    return <Logger spinner>Detecting...</Logger>;
  }

  return <Logger success>Complete!</Logger>;
};
```

## Distribution

**Homebrew**: Primary distribution method via `clix.rb` formula (to be updated for TypeScript build)
- Tap: `clix-so/clix-cli`
- Formula lives in this repo at `clix.rb`

**GitHub Releases**: Binaries published for macOS (amd64/arm64) and Linux (amd64/arm64)

## Important Notes

- The CLI modifies user project files (Xcode projects, Gradle files, manifests). Changes should be reviewable and reversible.
- Installation process is interactive with user prompts for credentials
- Supports `--dry-run` and `--verbose` flags for iOS installation
- All platform-specific code is isolated in respective `src/packages/` subdirectories
- Commands gracefully handle missing platforms (warn user to specify flags)
- Uses ES modules with `.js` extensions in imports (TypeScript convention for Node.js ESM)
- All file operations use `fs/promises` for async operations
- Shell commands executed via `execa` for better cross-platform support

## Migration Notes

This project was migrated from Go to TypeScript. Key changes:
- Go Cobra commands → Commander + React Ink components
- `pkg/` directory → `src/packages/` directory
- Go's `os` package → Node.js `fs/promises`
- Go's `exec.Command` → `execa`
- Custom logging → React Ink `<Logger>` component
- All functionality preserved with equivalent TypeScript implementations
