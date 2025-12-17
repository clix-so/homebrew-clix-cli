import React, { useState } from 'react';
import { Box, Text } from 'ink';
import SelectInput from 'ink-select-input';
import type { CLITool } from '../../lib/llm.js';

interface ToolSelectorProps {
  tools: CLITool[];
  currentTool?: string;
  onSelect: (tool: CLITool) => void;
}

export const ToolSelector: React.FC<ToolSelectorProps> = ({ tools, currentTool, onSelect }) => {
  const items = tools.map((tool) => ({
    label: tool.displayName,
    value: tool,
  }));

  return (
    <Box flexDirection="column">
      {currentTool && (
        <Box marginBottom={1}>
          <Text color="gray">Current selection: </Text>
          <Text bold color="cyan">
            {currentTool}
          </Text>
        </Box>
      )}
      <Text>
        <Text bold>Select an AI CLI tool:</Text>
      </Text>
      <Box marginTop={1}>
        <SelectInput items={items} onSelect={(item) => onSelect(item.value)} />
      </Box>
    </Box>
  );
};
