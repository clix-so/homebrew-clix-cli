import React from 'react';
import { render } from 'ink';
import { ConfigUI } from '../ui/ConfigUI.js';

export async function configCommand(): Promise<void> {
  return new Promise((resolve) => {
    const { unmount } = render(<ConfigUI onComplete={() => {
      unmount();
      resolve();
    }} />);
  });
}
