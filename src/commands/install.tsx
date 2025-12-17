import React from 'react';
import { render } from 'ink';
import { InstallUI } from '../ui/InstallUI.js';
import { configCommand } from './config.js';

interface InstallOptions {
  promptUrl?: string;
}

export async function installCommand(options: InstallOptions = {}): Promise<void> {
  return new Promise((resolve, reject) => {
    const { unmount } = render(
      <InstallUI
        promptUrl={options.promptUrl}
        onComplete={() => {
          unmount();
          resolve();
        }}
        onNeedsConfig={async () => {
          unmount();
          try {
            // Run config command first
            await configCommand();
            // After config, retry install
            await installCommand(options);
            resolve();
          } catch (error) {
            reject(error);
          }
        }}
      />
    );
  });
}
