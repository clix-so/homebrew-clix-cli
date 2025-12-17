import React from 'react';
import { Box, Text } from 'ink';

export const Banner: React.FC = () => {
  return (
    <Box flexDirection="column" marginBottom={2}>
      <Text bold>Clix CLI</Text>
      <Text dimColor>AI-powered mobile SDK installer</Text>
    </Box>
  );
};
