import React from 'react';
import { render } from 'ink';
import { InstallUI } from '../ui/InstallUI.js';
import { configCommand } from './config.js';

interface InstallOptions {
  promptUrl?: string;
}

export async function installCommand(options: InstallOptions = {}): Promise<void> {
  return new Promise((resolve, reject) => {
    let needsConfig = false;

    const { unmount } = render(
      <InstallUI
        promptUrl={options.promptUrl}
        onComplete={() => {
          unmount();
          resolve();
        }}
        onNeedsConfig={() => {
          needsConfig = true;
          unmount();
        }}
      />
    );

    // If config is needed, run config command first
    if (needsConfig) {
      configCommand()
        .then(() => {
          // After config, retry install
          return installCommand(options);
        })
        .then(resolve)
        .catch(reject);
    }
  });
}
