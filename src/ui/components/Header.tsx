import React from 'react';
import { Box, Text } from 'ink';

interface HeaderProps {
  title: string;
}

export const Header: React.FC<HeaderProps> = ({ title }) => {
  return (
    <Box flexDirection="column" marginBottom={2}>
      <Text bold>{title}</Text>
    </Box>
  );
};
