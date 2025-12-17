import { execaCommand } from 'execa';

export interface CLITool {
  name: string;
  command: string;
  displayName: string;
}

export const SUPPORTED_TOOLS: CLITool[] = [
  { name: 'claude', command: 'claude', displayName: 'Claude CLI' },
  { name: 'gemini', command: 'gemini', displayName: 'Google Gemini CLI' },
  { name: 'gpt', command: 'gpt', displayName: 'OpenAI GPT CLI' },
  { name: 'aider', command: 'aider', displayName: 'Aider' },
];

async function isCommandAvailable(command: string): Promise<boolean> {
  try {
    await execaCommand(`command -v ${command}`, { shell: true });
    return true;
  } catch {
    return false;
  }
}

export async function detectAvailableTools(): Promise<CLITool[]> {
  const available: CLITool[] = [];

  for (const tool of SUPPORTED_TOOLS) {
    if (await isCommandAvailable(tool.command)) {
      available.push(tool);
    }
  }

  return available;
}

export function getToolByName(name: string): CLITool | undefined {
  return SUPPORTED_TOOLS.find((tool) => tool.name === name);
}
