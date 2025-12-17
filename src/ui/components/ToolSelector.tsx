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
    key: tool.name,
    label: tool.displayName,
    value: tool,
  }));

  const Indicator: React.FC<{ isSelected?: boolean }> = ({ isSelected }) => {
    return (
      <Box marginRight={1}>
        <Text color={isSelected ? 'blue' : undefined}>{isSelected ? '→' : ' '}</Text>
      </Box>
    );
  };

  const Item: React.FC<{ isSelected?: boolean; label: string }> = ({ isSelected, label }) => {
    return <Text color={isSelected ? 'blue' : undefined}>{label}</Text>;
  };

  return (
    <Box flexDirection="column">
      {currentTool && (
        <Box marginBottom={1}>
          <Text dimColor>Currently configured: </Text>
          <Text bold>{currentTool}</Text>
        </Box>
      )}
      <Text dimColor>Select AI tool:</Text>
      <Box marginTop={1}>
        <SelectInput
          items={items}
          onSelect={(item) => onSelect(item.value)}
          indicatorComponent={Indicator}
          itemComponent={Item}
        />
      </Box>
    </Box>
  );
};
