import React, { useEffect, useState } from 'react';
import { Box } from 'ink';
import { detectPlatform } from '../utils/detectPlatform.js';
import { Logger } from '../ui/Logger.js';
import { uninstallClixIOS } from '../packages/ios/uninstaller.js';
import { uninstallClixAndroid } from '../packages/android/uninstaller.js';

interface UninstallCommandProps {
  ios?: boolean;
  android?: boolean;
}

export const UninstallCommand: React.FC<UninstallCommandProps> = (props) => {
  const [status, setStatus] = useState<string>('detecting');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        let { ios, android } = props;

        // Auto-detect if no flags provided
        if (!ios && !android) {
          const detected = await detectPlatform();
          ios = detected.isIOS;
          android = detected.isAndroid;

          if (!ios && !android) {
            setError('Could not detect platform. Please specify --ios or --android');
            setStatus('error');
            return;
          }
        }

        setStatus('uninstalling');

        if (ios) {
          await uninstallClixIOS();
        }

        if (android) {
          await uninstallClixAndroid();
        }

        setStatus('complete');
        process.exit(0);
      } catch (err: any) {
        setError(err.message || String(err));
        setStatus('error');
        process.exit(1);
      }
    })();
  }, []);

  if (status === 'detecting') {
    return <Logger spinner>Detecting platform...</Logger>;
  }

  if (status === 'error') {
    return (
      <Box flexDirection="column">
        <Logger failure>{error}</Logger>
      </Box>
    );
  }

  if (status === 'uninstalling') {
    return <Logger spinner>Uninstalling Clix SDK...</Logger>;
  }

  return (
    <Box flexDirection="column">
      <Logger success>Clix SDK uninstalled successfully!</Logger>
    </Box>
  );
};
