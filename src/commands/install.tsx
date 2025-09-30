import React, { useEffect, useState } from 'react';
import { Box, Text } from 'ink';
import { detectAllPlatforms } from '../utils/detectPlatform.js';
import { prompt } from '../utils/prompt.js';
import { Logger, Separator, NewLine } from '../ui/Logger.js';
import { handleIOSInstall } from '../packages/ios/installer.js';
import { handleAndroidInstall } from '../packages/android/installer.js';
import { handleExpoInstall } from '../packages/expo/installer.js';
import { handleFlutterInstall } from '../packages/flutter/installer.js';

interface InstallCommandProps {
  ios?: boolean;
  android?: boolean;
  expo?: boolean;
  flutter?: boolean;
  verbose?: boolean;
  dryRun?: boolean;
}

export const InstallCommand: React.FC<InstallCommandProps> = (props) => {
  const [status, setStatus] = useState<string>('detecting');
  const [error, setError] = useState<string | null>(null);
  const [platforms, setPlatforms] = useState({ ios: false, android: false, expo: false, flutter: false });

  useEffect(() => {
    (async () => {
      try {
        let { ios, android, expo, flutter } = props;

        // Auto-detect if no flags provided
        if (!ios && !android && !expo && !flutter) {
          const detected = await detectAllPlatforms();
          ios = detected.isIOS;
          android = detected.isAndroid;
          expo = detected.isExpo;
          flutter = detected.isFlutter;

          if (!ios && !android && !expo && !flutter) {
            setError('Could not detect platform. Please specify --ios, --android, --expo, or --flutter');
            setStatus('error');
            return;
          }
        }

        setPlatforms({ ios: !!ios, android: !!android, expo: !!expo, flutter: !!flutter });
        setStatus('installing');

        // Get credentials
        const projectID = await prompt('Enter your Project ID');
        const apiKey = await prompt('Enter your Public API Key');

        // Install for each platform
        if (ios) {
          await handleIOSInstall(projectID, apiKey);
        }

        if (android) {
          await handleAndroidInstall(apiKey, projectID);
        }

        if (expo) {
          await handleExpoInstall(apiKey, projectID);
        }

        if (flutter) {
          await handleFlutterInstall(apiKey, projectID);
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

  if (status === 'installing') {
    return <Logger spinner>Installing Clix SDK...</Logger>;
  }

  return (
    <Box flexDirection="column">
      <Logger success>Clix SDK installation complete!</Logger>
    </Box>
  );
};
