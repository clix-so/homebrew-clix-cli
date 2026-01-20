# Homebrew Tap for Clix CLI

This is the official Homebrew tap for [Clix CLI](https://github.com/clix-so/clix-cli).

## Installation

```bash
# Add the tap
brew tap clix-so/clix-cli

# Install clix
brew install --cask clix
```

Or install directly in one command:

```bash
brew install --cask clix-so/clix-cli/clix
```

## Upgrade

```bash
brew update
brew upgrade --cask clix
```

## Uninstall

### Basic Uninstall (Binary Only)
```bash
brew uninstall --cask clix
```

### Complete Uninstall (Binary + User Data)
Remove the app and all configuration/session files:
```bash
brew uninstall --zap --cask clix
```

This will remove:
- Binary: `/opt/homebrew/bin/clix` (or `/usr/local/bin/clix`)
- Configuration: `~/.config/clix/`
- Session files: `~/.local/state/clix/`

### Remove Tap (Optional)
```bash
brew untap clix-so/clix-cli
```

## About Clix CLI

Clix CLI is an interactive AI-powered assistant for Clix SDK development. Built with React/Ink for terminal UI, it supports multiple AI agents (Claude, Codex, Gemini, OpenCode, Cursor, Copilot) with streaming responses, slash commands, and pre-built skills for SDK workflows.

For more information, see the [main repository](https://github.com/clix-so/clix-cli).

## Alternative Installation Methods

### Via npm (requires Node.js 20+)
```bash
npm install -g @clix-so/clix-cli
```

### Via shell script (standalone binary)
```bash
curl -fsSL https://cli.clix.so/install.sh | bash
```

## Links

- [Clix CLI Repository](https://github.com/clix-so/clix-cli)
- [npm Package](https://www.npmjs.com/package/@clix-so/clix-cli)
- [Clix Documentation](https://docs.clix.so)

## License

[MIT with Custom Restrictions](LICENSE)
