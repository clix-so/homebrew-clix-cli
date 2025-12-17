import fetch from 'node-fetch';

export const DEFAULT_PROMPT_URL =
  'https://raw.githubusercontent.com/clix-so/cli-prompt/refs/heads/main/prompt.txt';

export class PromptFetcher {
  async fetch(url?: string): Promise<string> {
    const targetUrl = url || DEFAULT_PROMPT_URL;

    const response = await fetch(targetUrl);

    if (!response.ok) {
      throw new Error(`Failed to fetch prompt: ${response.statusText}`);
    }

    return await response.text();
  }
}
