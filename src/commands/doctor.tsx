import React, { useEffect, useState } from 'react';
import { Box } from 'ink';
import { detectAllPlatforms } from '../utils/detectPlatform.js';
import { Logger, Separator } from '../ui/Logger.js';
import { runIOSDoctor } from '../packages/ios/doctor.js';
import { runAndroidDoctor } from '../packages/android/doctor.js';
import { runExpoDoctor } from '../packages/expo/doctor.js';
import { runFlutterDoctor } from '../packages/flutter/doctor.js';

interface DoctorCommandProps {
  ios?: boolean;
  android?: boolean;
  expo?: boolean;
  flutter?: boolean;
}

export const DoctorCommand: React.FC<DoctorCommandProps> = (props) => {
  const [status, setStatus] = useState<string>('detecting');
  const [error, setError] = useState<string | null>(null);

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

        setStatus('checking');

        // Run doctor for each platform
        if (ios) {
          console.log('🔍 Checking Clix SDK integration for iOS...');
          await runIOSDoctor();
        }

        if (android) {
          console.log('🔍 Checking Clix SDK integration for Android...');
          console.log('─────────────────────────────────────────────────────');
          await runAndroidDoctor('');
        }

        if (expo) {
          await runExpoDoctor();
        }

        if (flutter) {
          console.log('🔍 Checking Clix SDK integration for Flutter...');
          await runFlutterDoctor();
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

  if (status === 'checking') {
    return <Logger spinner>Checking integration...</Logger>;
  }

  return (
    <Box flexDirection="column">
      <Logger success>Doctor check complete!</Logger>
    </Box>
  );
};
