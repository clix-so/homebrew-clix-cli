import React from 'react';
import { Box, Text } from 'ink';
import chalk from 'chalk';

interface HeaderProps {
  title: string;
}

export const Header: React.FC<HeaderProps> = ({ title }) => {
  return (
    <Box flexDirection="column" marginBottom={1}>
      <Text bold color="cyan">
        {title}
      </Text>
      <Text color="gray">
        {'━'.repeat(Math.min(title.length, 50))}
      </Text>
    </Box>
  );
};
