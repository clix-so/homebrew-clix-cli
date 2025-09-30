import { readFile, writeFile, readdir } from 'fs/promises';
import { join } from 'path';

/**
 * Removes Clix SDK initialization from Application.kt and related files
 */
export async function uninstallClixAndroid(): Promise<void> {
  const projectRoot = process.cwd();
  const kotlinDir = join(projectRoot, 'app', 'src', 'main', 'kotlin');

  let removed = false;

  const walk = async (dir: string): Promise<void> => {
    try {
      const files = await readdir(dir, { withFileTypes: true });

      for (const file of files) {
        const fullPath = join(dir, file.name);

        if (file.isDirectory()) {
          await walk(fullPath);
        } else if (file.name.endsWith('Application.kt')) {
          try {
            const data = await readFile(fullPath, 'utf-8');
            const lines = data.split('\n');
            const cleaned: string[] = [];

            for (const line of lines) {
              const trimmed = line.trim();
              if (
                trimmed.includes('Clix.initialize') ||
                trimmed.includes('ClixConfig(') ||
                trimmed.includes('import so.clix.Clix') ||
                trimmed.includes('import so.clix.ClixConfig') ||
                trimmed.includes('import so.clix.ClixLogLevel')
              ) {
                continue;
              }
              cleaned.push(line);
            }

            await writeFile(fullPath, cleaned.join('\n'), 'utf-8');
            removed = true;
          } catch {
            // Failed to process this file, continue
          }
        }
      }
    } catch {
      // Directory doesn't exist or can't be read
    }
  };

  await walk(kotlinDir);

  if (removed) {
    console.log('All the Clix SDK code has been removed from the Android Application class.');
  } else {
    throw new Error('No Application.kt with Clix SDK code found.');
  }
}
