import React, { useState, useEffect } from 'react';
import { Box, Text } from 'ink';
import { Header } from './components/Header.js';
import { StatusMessage } from './components/StatusMessage.js';
import { ConfigManager } from '../lib/config.js';
import { detectAvailableTools, getToolByName } from '../lib/llm.js';
import { PromptFetcher } from '../lib/prompt.js';
import { Executor } from '../lib/executor.js';
import { MCPInstaller } from '../lib/mcp.js';

interface InstallUIProps {
  promptUrl?: string;
  onComplete: () => void;
  onNeedsConfig: () => void;
}

export const InstallUI: React.FC<InstallUIProps> = ({ promptUrl, onComplete, onNeedsConfig }) => {
  const [phase, setPhase] = useState<
    | 'loading_config'
    | 'checking_tool'
    | 'checking_mcp'
    | 'fetching_prompt'
    | 'executing'
    | 'complete'
    | 'error'
  >('loading_config');
  const [toolName, setToolName] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');
  const [promptSize, setPromptSize] = useState<number>(0);

  useEffect(() => {
    const run = async () => {
      try {
        // Load configuration
        setPhase('loading_config');
        const config = new ConfigManager();
        const cfg = await config.load();

        if (!cfg.selectedCLI) {
          onNeedsConfig();
          return;
        }

        // Check if tool is available
        setPhase('checking_tool');
        const tool = getToolByName(cfg.selectedCLI);
        if (!tool) {
          setErrorMessage(`Configured tool '${cfg.selectedCLI}' not found`);
          setPhase('error');
          return;
        }

        const available = await detectAvailableTools();
        const found = available.find((t) => t.name === tool.name);
        if (!found) {
          setErrorMessage(
            `Configured tool '${tool.displayName}' is not available. Please run 'clix config' to select a different tool.`
          );
          setPhase('error');
          return;
        }

        setToolName(tool.displayName);

        // Ensure MCP server is installed
        setPhase('checking_mcp');
        const mcpInstaller = new MCPInstaller();
        const wasInstalled = await mcpInstaller.ensureServerInstalled(tool.name);

        // Fetch installation prompt
        setPhase('fetching_prompt');
        const fetcher = new PromptFetcher();
        const prompt = await fetcher.fetch(promptUrl);
        setPromptSize(prompt.length);

        // Execute installation
        setPhase('executing');
        const executor = new Executor(tool);
        await executor.executeInteractive(prompt);

        setPhase('complete');
        setTimeout(() => {
          onComplete();
        }, 1000);
      } catch (error) {
        setErrorMessage(error instanceof Error ? error.message : 'Unknown error');
        setPhase('error');
      }
    };

    run();
  }, [promptUrl, onComplete, onNeedsConfig]);

  return (
    <Box flexDirection="column" padding={1}>
      <Header title="Install Clix Mobile SDK" />

      {phase === 'loading_config' && (
        <StatusMessage type="loading" message="Loading configuration..." />
      )}

      {phase === 'checking_tool' && (
        <StatusMessage type="loading" message="Checking AI CLI tool availability..." />
      )}

      {phase === 'checking_mcp' && (
        <Box flexDirection="column">
          <StatusMessage type="success" message={`Using ${toolName} for installation`} />
          <StatusMessage type="loading" message="Checking Clix MCP Server..." />
        </Box>
      )}

      {phase === 'fetching_prompt' && (
        <Box flexDirection="column">
          <StatusMessage type="success" message={`Using ${toolName} for installation`} />
          <StatusMessage type="success" message="MCP Server configured" />
          <StatusMessage type="loading" message="Fetching installation instructions..." />
        </Box>
      )}

      {phase === 'executing' && (
        <Box flexDirection="column">
          <StatusMessage type="success" message={`Using ${toolName} for installation`} />
          <StatusMessage type="success" message="MCP Server configured" />
          <StatusMessage
            type="success"
            message={`Installation instructions loaded (${promptSize} bytes)`}
          />
          <Box marginTop={1} marginBottom={1}>
            <Text bold color="cyan">
              🚀 Starting AI-assisted installation...
            </Text>
          </Box>
          <Text color="gray">{'━'.repeat(50)}</Text>
          <Box marginTop={1}>
            <Text color="yellow">
              Note: The AI assistant is now taking over. Follow its instructions to complete the
              installation.
            </Text>
          </Box>
        </Box>
      )}

      {phase === 'complete' && (
        <Box flexDirection="column">
          <StatusMessage type="success" message="Installation process completed!" />
          <Box marginTop={1}>
            <Text>The AI assistant has finished processing the installation instructions.</Text>
          </Box>
          <Box>
            <Text>Please verify that the Clix Mobile SDK has been installed correctly.</Text>
          </Box>
        </Box>
      )}

      {phase === 'error' && (
        <Box flexDirection="column">
          <StatusMessage type="error" message={errorMessage} />
        </Box>
      )}
    </Box>
  );
};
