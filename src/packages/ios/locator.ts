import { promises as fs } from 'fs';
import path from 'path';

/**
 * Finds the app path by locating the .xcodeproj directory
 * @returns The path to the app directory (without .xcodeproj extension)
 * @throws Error if no .xcodeproj is found
 */
export async function findAppPath(): Promise<string> {
  // Check if there is a .xcodeproj folder in the current directory
  const entries = await fs.readdir('.', { withFileTypes: true });

  let projectName: string | null = null;

  for (const entry of entries) {
    if (entry.isDirectory() && entry.name.endsWith('.xcodeproj')) {
      projectName = entry.name.replace('.xcodeproj', '');
      break;
    }
  }

  if (!projectName) {
    throw new Error('❌ No .xcodeproj found. Please run this command from the root of your Xcode project');
  }

  return path.join('.', projectName);
}
