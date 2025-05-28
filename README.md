# clix

A command-line tool for integrating and managing the Clix SDK in your mobile projects. Clix SDK provides robust support for app push notifications and related features. This CLI helps automate installation, configuration checks, and (soon) uninstallation for both iOS and Android projects.

## Getting Started

### Install via Homebrew (Recommended)
```sh
brew tap clix-so/clix-cli
brew install clix-so/clix-cli/clix
```

### Install via Source
```sh
git clone https://github.com/clix-so/clix-cli.git
cd clix-cli
make install
```

## Requirements
- For iOS: Xcode project in the current directory


## Features

- **Install Clix SDK**
    - iOS (via Swift Package Manager or CocoaPods)
    - Android (via Gradle)
- **Doctor (Integration Checker)**
    - iOS: Checks your Xcode project for all required Clix SDK and push notification settings, and provides step-by-step guidance for any issues found.
    - Android: Checks your Android project for all required Clix SDK and push notification settings, and provides step-by-step guidance for any issues found.
- **Uninstall Clix SDK**

## Commands

### `install`
Install the Clix SDK into your project.

```
clix install
```

### `doctor`
Check that your project is correctly set up for push notifications and Clix SDK integration. Provides detailed diagnostics and suggestions for fixing any issues.

```
clix doctor
```

### `uninstall`
Remove the Clix SDK from your project.

```
clix uninstall
```

## Contributing
Pull requests and issues are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details (to be added).

## License
MIT

---

> **Note:**
> - For any issues or feature requests, please open an issue on GitHub.
