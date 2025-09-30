import { promises as fs } from 'fs';
import path from 'path';
import { findAppPath } from './locator.js';

/**
 * Uninstalls Clix SDK from the iOS project by removing related code from AppDelegate.swift
 */
export async function uninstallClixIOS(): Promise<void> {
  const appPath = await findAppPath();
  const appDelegatePath = path.join(appPath, 'AppDelegate.swift');

  try {
    await fs.access(appDelegatePath);
  } catch (error) {
    throw new Error(`failed to find AppDelegate.swift: ${error instanceof Error ? error.message : String(error)}`);
  }

  const content = await fs.readFile(appDelegatePath, 'utf-8');
  const lines = content.split('\n');
  const cleaned: string[] = [];

  for (const line of lines) {
    const trimmed = line.trim();

    // Skip lines that contain Clix or Firebase related code
    if (
      trimmed.includes('import Clix') ||
      trimmed.includes('import Firebase') ||
      trimmed.includes('FirebaseApp.configure()') ||
      trimmed.includes('Clix.initialize') ||
      trimmed.includes('UNUserNotificationCenter.current().delegate = self') ||
      trimmed.startsWith('Clix.')
    ) {
      continue;
    }

    cleaned.push(line);
  }

  await fs.writeFile(appDelegatePath, cleaned.join('\n'), 'utf-8');
  console.log('✅ All the Clix SDK code has been removed from the project.');
}
