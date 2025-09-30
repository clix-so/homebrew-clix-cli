#!/usr/bin/env bun
import React from 'react';
import { render } from 'ink';
import { Command } from 'commander';
import { InstallCommand } from './commands/install.js';
import { DoctorCommand } from './commands/doctor.js';
import { UninstallCommand } from './commands/uninstall.js';

const program = new Command();

program
  .name('clix')
  .description('A CLI tool for integrating and managing the Clix SDK in your mobile projects')
  .version('0.2.8');

program
  .command('install')
  .description('Install Clix SDK into your project')
  .option('--ios', 'Install Clix for iOS')
  .option('--android', 'Install Clix for Android')
  .option('--expo', 'Install Clix for React Native Expo')
  .option('--flutter', 'Install Clix for Flutter')
  .option('--verbose', 'Show verbose output during installation')
  .option('--dry-run', 'Show what would be changed without making changes')
  .action((options) => {
    render(<InstallCommand {...options} />);
  });

program
  .command('doctor')
  .description('Check Clix SDK integration status')
  .option('--ios', 'Check Clix for iOS')
  .option('--android', 'Check Clix for Android')
  .option('--expo', 'Check Clix for React Native Expo')
  .option('--flutter', 'Check Clix for Flutter')
  .action((options) => {
    render(<DoctorCommand {...options} />);
  });

program
  .command('uninstall')
  .description('Uninstall clix from devices')
  .option('--ios', 'Uninstall clix from iOS device')
  .option('--android', 'Uninstall clix from Android device')
  .action((options) => {
    render(<UninstallCommand {...options} />);
  });

program.parse();
