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
    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const detect = async () => {
      try {
        await delay(800);
        const tools = await detectAvailableTools();

        if (tools.length === 0) {
          setErrorMessage('No supported AI CLI tools found');
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
        await delay(400);
        setPhase('selecting');
      } catch (error) {
        setErrorMessage(error instanceof Error ? error.message : 'Unknown error');
        setPhase('error');
      }
    };

    detect();
  }, []);

  const handleSelect = async (tool: CLITool) => {
    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    setPhase('saving');
    try {
      await delay(600);
      const config = new ConfigManager();
      await config.save({ selectedCLI: tool.name });
      setCurrentTool(tool.displayName);
      setPhase('complete');
      setTimeout(() => {
        onComplete();
      }, 1500);
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Failed to save configuration');
      setPhase('error');
    }
  };

  return (
    <Box flexDirection="column" padding={1}>
      <Header title="Configure AI CLI Tool" />

      {phase === 'detecting' && (
        <StatusMessage type="loading" message="Detecting AI CLI tools..." />
      )}

      {phase === 'selecting' && (
        <Box flexDirection="column">
          <StatusMessage
            type="success"
            message={`Found ${availableTools.length} available tool${availableTools.length > 1 ? 's' : ''}`}
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
          message={`Configured to use ${currentTool}`}
        />
      )}

      {phase === 'error' && (
        <Box flexDirection="column">
          <StatusMessage type="error" message={errorMessage} />
          <Box marginTop={1} marginBottom={1}>
            <Text dimColor>Supported AI CLI tools:</Text>
          </Box>
          {SUPPORTED_TOOLS.map((tool) => (
            <Box key={tool.name} marginLeft={2}>
              <Text dimColor>  </Text>
              <Text>{tool.displayName}</Text>
              <Text dimColor> - {tool.command}</Text>
            </Box>
          ))}
          <Box marginTop={1}>
            <Text dimColor>Install one of these tools and try again.</Text>
          </Box>
        </Box>
      )}
    </Box>
  );
};
