import { readFile, writeFile, mkdir } from 'node:fs/promises';
import { homedir } from 'node:os';
import { join } from 'node:path';
import { existsSync } from 'node:fs';

export interface Config {
  selectedCLI: string;
}

export class ConfigManager {
  private configPath: string;
  private configDir: string;

  constructor() {
    this.configDir = join(homedir(), '.clix');
    this.configPath = join(this.configDir, 'config.json');
  }

  async getConfigPath(): Promise<string> {
    if (!existsSync(this.configDir)) {
      await mkdir(this.configDir, { recursive: true, mode: 0o755 });
    }
    return this.configPath;
  }

  async load(): Promise<Config> {
    try {
      const path = await this.getConfigPath();
      if (!existsSync(path)) {
        return { selectedCLI: '' };
      }

      const data = await readFile(path, 'utf-8');
      return JSON.parse(data) as Config;
    } catch {
      return { selectedCLI: '' };
    }
  }

  async save(config: Config): Promise<void> {
    const path = await this.getConfigPath();
    const data = JSON.stringify(config, null, 2);
    await writeFile(path, data, { mode: 0o644 });
  }
}
