import React from 'react';
import { Box, Text } from 'ink';
import { Banner } from '../ui/components/Banner.js';

export const RootCommand: React.FC = () => {
  return (
    <Box flexDirection="column" padding={1}>
      <Banner />
      <Box flexDirection="column" marginTop={1}>
        <Text>
          <Text bold color="cyan">
            Welcome to Clix Mobile SDK Installer!
          </Text>
        </Text>
        <Box marginTop={1} flexDirection="column">
          <Text color="gray">Available commands:</Text>
          <Text>
            {'  '}
            <Text bold color="green">
              clix install
            </Text>
            <Text color="gray"> - Install the SDK using AI assistance</Text>
          </Text>
          <Text>
            {'  '}
            <Text bold color="green">
              clix config
            </Text>
            <Text color="gray"> - Configure your AI CLI tool</Text>
          </Text>
        </Box>
        <Box marginTop={1}>
          <Text color="yellow">
            💡 Tip: Run <Text bold>clix install</Text> to get started!
          </Text>
        </Box>
      </Box>
    </Box>
  );
};
