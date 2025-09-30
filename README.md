# clix

A command-line tool for integrating and managing the Clix SDK in your mobile projects. Clix SDK provides robust support for app push notifications and related features. This CLI helps automate installation, configuration checks, and uninstallation for iOS, Android, React Native Expo, and Flutter projects.

## Getting Started

### Install via Homebrew (Recommended)
```sh
brew tap clix-so/clix-cli
brew install clix-so/clix-cli/clix
```

### Install via Source
```sh
git clone https://github.com/clix-so/homebrew-clix-cli.git
cd homebrew-clix-cli
bun install
bun run build
```

## Requirements
- Bun runtime (for development)
- For iOS: Xcode project in the current directory
- For Android: Android Studio project with Gradle
- For Expo: Expo project with app.json
- For Flutter: Flutter project with pubspec.yaml

## Technology Stack

This CLI is built with:
- **TypeScript** - Type-safe JavaScript
- **React + Ink** - Interactive CLI components
- **Bun** - Fast JavaScript runtime and bundler
- **Commander** - CLI framework

## Features

- **Install Clix SDK**
    - iOS (via Swift Package Manager or CocoaPods)
      - Auto-configures App Groups and NotificationServiceExtension
    - Android (via Gradle)
    - React Native Expo (via npm/yarn)
    - Flutter (via pub)
- **Doctor (Integration Checker)**
    - iOS: Checks your Xcode project for all required Clix SDK and push notification settings, and provides step-by-step guidance for any issues found.
    - Android: Checks your Android project for all required Clix SDK and push notification settings, and provides step-by-step guidance for any issues found.
    - Expo: Validates Expo project configuration and dependencies
    - Flutter: Verifies Flutter project setup and Firebase integration
- **Uninstall Clix SDK**

## Commands

### `install`
Install the Clix SDK into your project.

```
clix install
```

#### Platform-Specific Options

```
clix install --ios           # Install for iOS only
clix install --android       # Install for Android only
clix install --expo          # Install for React Native Expo only
clix install --flutter       # Install for Flutter only
clix install --verbose       # Show verbose output (iOS)
clix install --dry-run       # Preview changes without applying (iOS)
```

During installation, the CLI will automatically:
1. Detect your project type (iOS/Android/Expo/Flutter)
2. Install the Clix SDK using the appropriate package manager
3. Configure platform-specific settings
4. Set up push notifications and required permissions

### `doctor`
Check that your project is correctly set up for push notifications and Clix SDK integration. Provides detailed diagnostics and suggestions for fixing any issues.

```
clix doctor
```

#### Platform-Specific Options

```
clix doctor --ios           # Check iOS only
clix doctor --android       # Check Android only
clix doctor --expo          # Check Expo only
clix doctor --flutter       # Check Flutter only
```

### `uninstall`
Remove the Clix SDK from your project.

```
clix uninstall
```

#### Platform-Specific Options

```
clix uninstall --ios        # Uninstall from iOS
clix uninstall --android    # Uninstall from Android
```

## Development

### Setup
```sh
bun install
```

### Run in Development Mode
```sh
bun run dev
```

### Build
```sh
bun run build
```

### Type Check
```sh
bun run typecheck
```

## Contributing
Pull requests and issues are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details (to be added).

## License
MIT

---

> **Note:**
> - For any issues or feature requests, please open an issue on GitHub.
