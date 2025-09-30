import React, { ReactNode } from 'react';
import { Box, Text } from 'ink';
import Spinner from 'ink-spinner';

interface LoggerProps {
  children: ReactNode;
  indent?: number;
  branch?: boolean;
  gray?: boolean;
  bold?: boolean;
  code?: boolean;
  title?: boolean;
  success?: boolean;
  failure?: boolean;
  spinner?: boolean;
}

export const Logger: React.FC<LoggerProps> = ({
  children,
  indent = 0,
  branch = false,
  gray = false,
  bold = false,
  code = false,
  title = false,
  success = false,
  failure = false,
  spinner = false,
}) => {
  let prefix = '';
  if (branch) prefix += ' └ ';
  if (success) prefix += '✅ ';
  if (failure) prefix += '❌ ';

  const color = gray || code || title ? 'gray' : undefined;
  const fontBold = bold || title;

  return (
    <Box paddingLeft={indent}>
      {spinner && (
        <Text color="gray" bold>
          <Spinner type="dots" /> {prefix}
        </Text>
      )}
      {!spinner && prefix && <Text>{prefix}</Text>}
      <Text color={color} bold={fontBold}>
        {children}
      </Text>
    </Box>
  );
};

export const Separator: React.FC = () => (
  <Text>─────────────────────────────────────────────────────</Text>
);

export const NewLine: React.FC = () => <Text>{''}</Text>;
