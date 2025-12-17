import React from 'react';
import { Box, Text } from 'ink';
import Spinner from 'ink-spinner';

export type StatusType = 'info' | 'success' | 'warning' | 'error' | 'loading';

interface StatusMessageProps {
  type: StatusType;
  message: string;
}

const STATUS_ICONS = {
  info: '🔍',
  success: '✅',
  warning: '⚠️',
  error: '❌',
  loading: null,
};

const STATUS_COLORS = {
  info: 'blue',
  success: 'green',
  warning: 'yellow',
  error: 'red',
  loading: 'cyan',
} as const;

export const StatusMessage: React.FC<StatusMessageProps> = ({ type, message }) => {
  const icon = STATUS_ICONS[type];
  const color = STATUS_COLORS[type];

  return (
    <Box>
      {type === 'loading' ? (
        <Text color={color}>
          <Spinner type="dots" />
        </Text>
      ) : (
        <Text>{icon}</Text>
      )}
      <Text color={color}> {message}</Text>
    </Box>
  );
};
