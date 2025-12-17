# Development Guide

This guide covers local development, testing, and contributing to Clix CLI.

## Prerequisites

- Node.js 18 or higher
- npm or pnpm
- Git

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/clix-so/homebrew-clix-cli.git
cd homebrew-clix-cli

# 2. Install dependencies
npm install

# 3. Build the project
npm run build

# 4. Test locally
node dist/cli.js --help
```

## Development Workflow

### 1. Install Dependencies

```bash
npm install
```

This installs all required dependencies including:
- `ink` - React for CLIs
- `meow` - CLI argument parsing
- `chalk` - Terminal styling
- `execa` - Process execution
- TypeScript and build tools

### 2. Build the Project

```bash
# One-time build
npm run build

# Watch mode (rebuilds on file changes)
npm run dev
```

**Build output:** `dist/` directory
- `dist/cli.js` - Main bundled CLI
- `dist/cli.d.ts` - TypeScript definitions
- `dist/chunk-*.js` - Code-split chunks
- `dist/devtools-*.js` - React devtools (external)
- `dist/multipart-parser-*.js` - Multipart parser chunk

### 3. Test Locally

#### Method A: Direct Execution (Quick testing)

```bash
# Run directly with Node.js
node dist/cli.js --help
node dist/cli.js --version
node dist/cli.js config
node dist/cli.js install
```

#### Method B: npm link (Recommended for full testing)

```bash
# Link the CLI globally
npm link

# Now use 'clix' command anywhere
clix --help
clix config
clix install

# When done, unlink
npm unlink -g @clix-so/clix-cli
```

#### Method C: Watch Mode Development

```bash
# Terminal 1: Start watch mode
npm run dev

# Terminal 2: Test changes
node dist/cli.js --help
# or if linked:
clix --help
```

### 4. Type Checking

```bash
# Check types without building
npm run typecheck
```

## Project Structure

```
.
├── src/
│   ├── cli.tsx              # Main entry point
│   ├── commands/            # Command implementations
│   │   ├── config.tsx       # Config command
│   │   ├── install.tsx      # Install command
│   │   └── root.tsx         # Root/welcome command
│   ├── lib/                 # Core functionality
│   │   ├── config.ts        # Configuration management
│   │   ├── executor.ts      # AI tool execution
│   │   ├── llm.ts          # AI tool detection
│   │   ├── mcp.ts          # MCP server management
│   │   └── prompt.ts       # Prompt fetching
│   └── ui/                  # Ink UI components
│       ├── components/      # Reusable components
│       │   ├── Banner.tsx
│       │   ├── Header.tsx
│       │   ├── StatusMessage.tsx
│       │   └── ToolSelector.tsx
│       ├── ConfigUI.tsx     # Config screen
│       └── InstallUI.tsx    # Install screen
├── dist/                    # Build output (gitignored)
├── scripts/                 # Build scripts
│   └── add-shebang.js      # Adds shebang to CLI
├── package.json
├── tsconfig.json
└── tsup.config.ts          # Build configuration
```

## Making Changes

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Edit files in `src/`:
- Add new commands in `src/commands/`
- Add new UI components in `src/ui/components/`
- Add new utilities in `src/lib/`

### 3. Test Your Changes

```bash
# Build
npm run build

# Test
node dist/cli.js --help
# or
npm link && clix --help
```

### 4. Type Check

```bash
npm run typecheck
```

### 5. Commit Your Changes

```bash
git add .
git commit -m "feat: add your feature description"
```

Use [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks

## Testing in Real Projects

### Test with iOS Project

```bash
# 1. Link the CLI
cd /Users/jace-yoo/Documents/workspace/homebrew-clix-cli
npm link

# 2. Go to an iOS project
cd /path/to/ios/project

# 3. Test the install command
clix install

# 4. Verify the AI assistant launches and works correctly
```

### Test with Android Project

```bash
# 1. Link the CLI (if not already)
npm link

# 2. Go to an Android project
cd /path/to/android/project

# 3. Test the install command
clix install

# 4. Verify the AI assistant launches and works correctly
```

## Debugging

### Enable Debug Mode

```bash
export CLIX_DEBUG=1
clix install
```

### View Detailed Logs

```bash
# Add console.log statements in your code
console.log('Debug info:', someVariable);

# Rebuild and test
npm run build
node dist/cli.js install
```

### Check Build Output

```bash
# View bundled code
cat dist/cli.js | head -100

# Check file sizes
ls -lh dist/

# Verify shebang
head -n 1 dist/cli.js
# Should show: #!/usr/bin/env node
```

## Common Tasks

### Add a New Command

1. Create a new file in `src/commands/`:
```tsx
// src/commands/mycommand.tsx
import React from 'react';
import { render } from 'ink';

export function myCommand() {
  const ui = (
    <Box>
      <Text>My Command!</Text>
    </Box>
  );
  render(ui);
}
```

2. Register in `src/cli.tsx`:
```tsx
if (cli.input[0] === 'mycommand') {
  myCommand();
}
```

3. Build and test:
```bash
npm run build
node dist/cli.js mycommand
```

### Add a New UI Component

1. Create component in `src/ui/components/`:
```tsx
// src/ui/components/MyComponent.tsx
import React from 'react';
import { Box, Text } from 'ink';

interface Props {
  message: string;
}

export function MyComponent({ message }: Props) {
  return (
    <Box>
      <Text color="green">{message}</Text>
    </Box>
  );
}
```

2. Use in a command:
```tsx
import { MyComponent } from '../ui/components/MyComponent.js';

const ui = (
  <MyComponent message="Hello!" />
);
```

### Modify Build Configuration

Edit `tsup.config.ts`:
```typescript
export default defineConfig({
  entry: ['src/cli.tsx'],
  format: ['esm'],
  // ... other options
});
```

## CI/CD Testing

The repository has automated workflows:

### CI Workflow (`.github/workflows/ci.yml`)

Runs on every push and PR:
- Type checking
- Build verification
- Tests on Node 18, 20, 22

### Release Workflow (`.github/workflows/release.yml`)

Runs when version changes in `package.json`:
- Builds the project
- Publishes to npm
- Updates Homebrew formula
- Creates GitHub release

## Troubleshooting

### "Cannot find module" errors

**Cause:** Missing dependencies or incorrect imports

**Solution:**
```bash
npm install
npm run build
```

### "Permission denied" when running CLI

**Cause:** Shebang not added or file not executable

**Solution:**
```bash
npm run build  # Runs add-shebang.js automatically
# or manually:
chmod +x dist/cli.js
```

### TypeScript errors

**Cause:** Type mismatches or missing types

**Solution:**
```bash
npm run typecheck
# Fix errors in src/ files
```

### Build output is huge

**Cause:** All dependencies are bundled (expected)

**Current size:** ~2.3 MB (bundled with Ink, React, etc.)

This is normal for a CLI with a UI framework.

### npm link doesn't work

**Cause:** npm global bin not in PATH

**Solution:**
```bash
# Find npm global bin
npm config get prefix

# Add to PATH in ~/.bashrc or ~/.zshrc
export PATH="$PATH:$(npm config get prefix)/bin"

# Reload shell
source ~/.bashrc  # or source ~/.zshrc
```

## Code Style

- Use TypeScript for type safety
- Use React/Ink for UI components
- Follow existing code structure
- Add JSDoc comments for public APIs
- Keep components small and focused

## Before Submitting PR

- [ ] Code builds without errors (`npm run build`)
- [ ] Types check correctly (`npm run typecheck`)
- [ ] Tested locally with `npm link`
- [ ] Tested in a real project (if applicable)
- [ ] Commit messages follow Conventional Commits
- [ ] Updated documentation if needed

## Resources

- [Ink Documentation](https://github.com/vadimdemedes/ink) - React for CLIs
- [meow Documentation](https://github.com/sindresorhus/meow) - CLI framework
- [tsup Documentation](https://tsup.egoist.dev/) - Build tool
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)

## Getting Help

- [GitHub Issues](https://github.com/clix-so/homebrew-clix-cli/issues)
- [Discussions](https://github.com/clix-so/homebrew-clix-cli/discussions)

---

Happy coding! 🚀
