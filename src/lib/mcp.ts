import { readFile, writeFile, mkdir } from 'node:fs/promises';
import { homedir } from 'node:os';
import { join } from 'node:path';
import { existsSync } from 'node:fs';
import { execaCommand } from 'execa';

const MCP_SERVER_REPO = 'https://github.com/clix-so/clix-mcp-server';
const MCP_SERVER_NAME = 'clix-mcp-server';

interface MCPServer {
  command: string;
  args?: string[];
  env?: Record<string, string>;
}

interface MCPConfig {
  mcpServers: Record<string, MCPServer>;
}

export class MCPInstaller {
  async getMCPConfigPath(toolName: string): Promise<string> {
    const home = homedir();
    let configDir: string;
    let configFile: string;

    switch (toolName) {
      case 'claude':
        configDir = join(home, '.config', 'claude');
        configFile = 'claude_desktop_config.json';
        break;

      case 'aider':
        configDir = join(home, '.aider');
        configFile = 'mcp_config.json';
        break;

      case 'gpt':
      case 'openai':
        configDir = join(home, '.config', 'openai');
        configFile = 'mcp_config.json';
        break;

      case 'gemini':
        configDir = join(home, '.config', 'gemini');
        configFile = 'mcp_config.json';
        break;

      default:
        configDir = join(home, '.config', 'mcp');
        configFile = `${toolName}_config.json`;
    }

    if (!existsSync(configDir)) {
      await mkdir(configDir, { recursive: true, mode: 0o755 });
    }

    return join(configDir, configFile);
  }

  async isServerInstalled(toolName: string): Promise<boolean> {
    try {
      const configPath = await this.getMCPConfigPath(toolName);

      if (!existsSync(configPath)) {
        return false;
      }

      const data = await readFile(configPath, 'utf-8');
      const config: MCPConfig = JSON.parse(data);

      return MCP_SERVER_NAME in config.mcpServers;
    } catch {
      return false;
    }
  }

  async installServer(toolName: string): Promise<void> {
    // Check if npx is available
    try {
      await execaCommand('command -v npx', { shell: true });
    } catch {
      throw new Error('npx not found. Please install Node.js and npm first');
    }

    const configPath = await this.getMCPConfigPath(toolName);

    // Load existing config or create new one
    let config: MCPConfig = { mcpServers: {} };

    if (existsSync(configPath)) {
      try {
        const data = await readFile(configPath, 'utf-8');
        config = JSON.parse(data);
      } catch {
        // Invalid JSON, start fresh
      }
    }

    // Add our MCP server
    config.mcpServers[MCP_SERVER_NAME] = {
      command: 'npx',
      args: ['-y', MCP_SERVER_REPO],
    };

    // Write config back
    const data = JSON.stringify(config, null, 2);
    await writeFile(configPath, data, { mode: 0o644 });
  }

  async ensureServerInstalled(toolName: string): Promise<boolean> {
    const installed = await this.isServerInstalled(toolName);

    if (installed) {
      return true;
    }

    await this.installServer(toolName);
    return false;
  }
}
