import { execa } from 'execa';
import { writeFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import type { CLITool } from './llm.js';

export class Executor {
  constructor(private tool: CLITool) {}

  async executeInteractive(initialPrompt: string): Promise<void> {
    // Save prompt to temp file for reference
    const promptFile = await this.savePromptToTemp(initialPrompt);

    try {
      switch (this.tool.name) {
        case 'claude':
          // Claude CLI accepts the prompt as a direct argument
          await execa(this.tool.command, [initialPrompt], {
            stdio: 'inherit',
            cwd: process.cwd(),
          });
          break;

        case 'aider':
        case 'gpt':
        case 'gemini':
        default:
          // Auto-inject prompt via stdin
          await this.runWithAutoInput(initialPrompt);
          break;
      }
    } finally {
      // Clean up temp file
      await rm(promptFile, { force: true });
    }
  }

  private async runWithAutoInput(initialPrompt: string): Promise<void> {
    const subprocess = execa(this.tool.command, [], {
      stdio: ['pipe', 'inherit', 'inherit'],
      cwd: process.cwd(),
    });

    // Auto-inject the initial prompt
    if (subprocess.stdin) {
      subprocess.stdin.write(`${initialPrompt}\n`);
    }

    // Forward stdin to the subprocess
    process.stdin.pipe(subprocess.stdin!);

    await subprocess;
  }

  private async savePromptToTemp(prompt: string): Promise<string> {
    const tmpFile = join(tmpdir(), 'clix-install-prompt.txt');
    await writeFile(tmpFile, prompt, { mode: 0o644 });
    return tmpFile;
  }
}
