import React, { useState, useEffect } from 'react';
import { Box, Text } from 'ink';
import { Header } from './components/Header.js';
import { StatusMessage } from './components/StatusMessage.js';
import { ToolSelector } from './components/ToolSelector.js';
import { detectAvailableTools, getToolByName, SUPPORTED_TOOLS, type CLITool } from '../lib/llm.js';
import { ConfigManager } from '../lib/config.js';

interface ConfigUIProps {
  onComplete: () => void;
}

export const ConfigUI: React.FC<ConfigUIProps> = ({ onComplete }) => {
  const [phase, setPhase] = useState<'detecting' | 'selecting' | 'saving' | 'complete' | 'error'>(
    'detecting'
  );
  const [availableTools, setAvailableTools] = useState<CLITool[]>([]);
  const [currentTool, setCurrentTool] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');

  useEffect(() => {
    const detect = async () => {
      try {
        const tools = await detectAvailableTools();

        if (tools.length === 0) {
          setErrorMessage('No supported AI CLI tools found!');
          setPhase('error');
          return;
        }

        const config = new ConfigManager();
        const cfg = await config.load();

        if (cfg.selectedCLI) {
          const tool = getToolByName(cfg.selectedCLI);
          if (tool) {
            setCurrentTool(tool.displayName);
          }
        }

        setAvailableTools(tools);
        setPhase('selecting');
      } catch (error) {
        setErrorMessage(error instanceof Error ? error.message : 'Unknown error');
        setPhase('error');
      }
    };

    detect();
  }, []);

  const handleSelect = async (tool: CLITool) => {
    setPhase('saving');
    try {
      const config = new ConfigManager();
      await config.save({ selectedCLI: tool.name });
      setCurrentTool(tool.displayName);
      setPhase('complete');
      setTimeout(() => {
        onComplete();
      }, 1500);
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Failed to save config');
      setPhase('error');
    }
  };

  return (
    <Box flexDirection="column" padding={1}>
      <Header title="Configure AI CLI Tool" />

      {phase === 'detecting' && (
        <StatusMessage type="loading" message="Detecting available AI CLI tools..." />
      )}

      {phase === 'selecting' && (
        <Box flexDirection="column">
          <StatusMessage
            type="success"
            message={`Found ${availableTools.length} available AI CLI tool(s)`}
          />
          <Box marginTop={1}>
            <ToolSelector
              tools={availableTools}
              currentTool={currentTool}
              onSelect={handleSelect}
            />
          </Box>
        </Box>
      )}

      {phase === 'saving' && <StatusMessage type="loading" message="Saving configuration..." />}

      {phase === 'complete' && (
        <StatusMessage
          type="success"
          message={`Successfully configured to use ${currentTool}`}
        />
      )}

      {phase === 'error' && (
        <Box flexDirection="column">
          <StatusMessage type="error" message={errorMessage} />
          <Box marginTop={1}>
            <Text color="gray">Supported tools:</Text>
          </Box>
          {SUPPORTED_TOOLS.map((tool) => (
            <Text key={tool.name} color="gray">
              {'  '}• {tool.displayName} ({tool.command})
            </Text>
          ))}
          <Box marginTop={1}>
            <Text color="yellow">
              Please install one of the supported AI CLI tools and try again.
            </Text>
          </Box>
        </Box>
      )}
    </Box>
  );
};
