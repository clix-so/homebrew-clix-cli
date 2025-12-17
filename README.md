# Clix CLI

A command-line tool for integrating and managing the Clix SDK in your mobile projects. Clix SDK provides robust support for app push notifications and related features. This CLI uses AI-powered assistants to automate installation and configuration.

## ✨ Features

- 🤖 **AI-Powered Installation**: Uses AI assistants (Claude, Gemini, GPT, or Aider) to guide SDK installation
- 🎨 **Beautiful UI**: Built with [Ink](https://github.com/vadimdemedes/ink) for a modern, professional CLI experience
- 📱 **Multi-Platform**: Supports iOS and Android projects
- 🔌 **MCP Integration**: Automatic setup of Clix MCP Server for enhanced AI capabilities
- 💡 **Interactive**: Guided workflows with clear visual feedback

## 📦 Installation

### Install via npm (Recommended)

```bash
npm install -g @clix-so/clix-cli
```

### Install via pnpm

```bash
pnpm add -g @clix-so/clix-cli
```

### Install via yarn

```bash
yarn global add @clix-so/clix-cli
```

## 🚀 Getting Started

### Prerequisites

- Node.js 18 or higher
- One of the supported AI CLI tools:
  - [Claude CLI](https://claude.ai/cli)
  - [Google Gemini CLI](https://ai.google.dev/gemini-api)
  - [OpenAI GPT CLI](https://platform.openai.com/)
  - [Aider](https://github.com/paul-gauthier/aider)

### Quick Start

1. **Configure your AI CLI tool** (first time only):

```bash
clix config
```

The CLI will detect available AI tools and let you choose one.

2. **Install the Clix SDK**:

```bash
clix install
```

The AI assistant will guide you through the installation process!

## 📖 Commands

### `clix install`

Install the Clix SDK into your project using AI assistance.

```bash
clix install

# Use a custom installation prompt
clix install --prompt-url https://example.com/custom-prompt.txt
```

**Options:**
- `-p, --prompt-url <url>` - Custom URL for the installation prompt

### `clix config`

Configure which AI CLI tool to use for installations.

```bash
clix config
```

This will:
- Detect available AI CLI tools on your system
- Show your current selection (if any)
- Let you choose or change your preferred tool

### `clix` (no command)

Show welcome message and available commands.

```bash
clix
```

## 🎯 How It Works

1. **Configuration**: Clix detects and saves your preferred AI CLI tool
2. **MCP Setup**: Automatically configures the Clix MCP Server for enhanced capabilities
3. **Prompt Fetching**: Downloads the latest installation instructions
4. **AI Execution**: Launches your chosen AI assistant with the installation prompt
5. **Guided Installation**: The AI guides you through the entire SDK setup

## 🛠️ Development

### Prerequisites

- Node.js 18+
- npm or pnpm

### Setup

```bash
# Clone the repository
git clone https://github.com/clix-so/homebrew-clix-cli.git
cd homebrew-clix-cli

# Install dependencies
npm install

# Build
npm run build

# Run locally
node dist/cli.js
```

### Project Structure

```
src/
├── cli.tsx              # Main CLI entry point
├── commands/            # Command implementations
│   ├── config.tsx       # Config command
│   ├── install.tsx      # Install command
│   └── root.tsx         # Root command (welcome)
├── lib/                 # Core functionality
│   ├── config.ts        # Configuration management
│   ├── executor.ts      # AI tool execution
│   ├── llm.ts           # AI tool detection
│   ├── mcp.ts           # MCP server management
│   └── prompt.ts        # Prompt fetching
└── ui/                  # Ink UI components
    ├── components/      # Reusable UI components
    │   ├── Banner.tsx
    │   ├── Header.tsx
    │   ├── StatusMessage.tsx
    │   └── ToolSelector.tsx
    ├── ConfigUI.tsx     # Config screen
    └── InstallUI.tsx    # Install screen
```

## 🤝 Contributing

Pull requests and issues are welcome!

### Getting Started

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed development guide.

### Quick Start

```bash
# Install dependencies
npm install

# Build the project
npm run build

# Test locally
npm link
clix --help

# Watch mode for development
npm run dev
```

### Resources

- [DEVELOPMENT.md](DEVELOPMENT.md) - Comprehensive development guide
- [DEPLOYMENT.md](DEPLOYMENT.md) - Release and deployment instructions

## 📄 License

MIT

## 🔗 Links

- [GitHub Repository](https://github.com/clix-so/homebrew-clix-cli)
- [Issue Tracker](https://github.com/clix-so/homebrew-clix-cli/issues)
- [Clix SDK Documentation](https://clix.so)
- [LLMs.txt](llms.txt) - Detailed documentation for AI assistants

---

Made with ❤️ by the Clix team
