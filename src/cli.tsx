import React from 'react';
import { render } from 'ink';
import meow from 'meow';
import { installCommand } from './commands/install.js';
import { configCommand } from './commands/config.js';
import { RootCommand } from './commands/root.js';

const cli = meow(
  `
  Usage
    $ clix <command> [options]

  Commands
    install     Install Clix Mobile SDK using AI assistance
    config      Configure the AI CLI tool to use

  Options
    --help      Show this help message
    --version   Show version number

  Install Options
    --prompt-url, -p   Custom URL for the installation prompt

  Examples
    $ clix install
    $ clix install --prompt-url https://example.com/prompt.txt
    $ clix config
`,
  {
    importMeta: import.meta,
    flags: {
      promptUrl: {
        type: 'string',
        shortFlag: 'p',
      },
    },
  }
);

async function main() {
  const command = cli.input[0];

  try {
    switch (command) {
      case 'install':
        await installCommand({
          promptUrl: cli.flags.promptUrl,
        });
        break;

      case 'config':
        await configCommand();
        break;

      default:
        // Show welcome message
        const { unmount } = render(<RootCommand />);
        // Auto-unmount after showing the message
        setTimeout(() => {
          unmount();
          process.exit(0);
        }, 100);
        break;
    }
  } catch (error) {
    console.error('Error:', error instanceof Error ? error.message : 'Unknown error');
    process.exit(1);
  }
}

main();
