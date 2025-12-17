import React from 'react';
import { Box, Text } from 'ink';

export const Banner: React.FC = () => {
  return (
    <Box flexDirection="column" marginBottom={1}>
      <Text bold color="magenta">
        ╔═══════════════════════════════════════════╗
      </Text>
      <Text bold color="magenta">
        ║     🚀 Clix Mobile SDK Installer 🚀      ║
      </Text>
      <Text bold color="magenta">
        ║   AI-Powered Installation Assistant      ║
      </Text>
      <Text bold color="magenta">
        ╚═══════════════════════════════════════════╝
      </Text>
    </Box>
  );
};
