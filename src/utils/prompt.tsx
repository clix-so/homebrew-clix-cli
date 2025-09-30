import React, { useState } from 'react';
import { render, Box } from 'ink';
import TextInput from 'ink-text-input';

interface PromptProps {
  question: string;
  onSubmit: (value: string) => void;
}

const PromptComponent: React.FC<PromptProps> = ({ question, onSubmit }) => {
  const [value, setValue] = useState('');

  const handleSubmit = () => {
    onSubmit(value);
  };

  return (
    <Box flexDirection="column">
      <Box>
        {question}: <TextInput value={value} onChange={setValue} onSubmit={handleSubmit} />
      </Box>
    </Box>
  );
};

export async function prompt(question: string): Promise<string> {
  return new Promise((resolve) => {
    const { unmount } = render(
      <PromptComponent
        question={question}
        onSubmit={(value) => {
          unmount();
          resolve(value);
        }}
      />
    );
  });
}
