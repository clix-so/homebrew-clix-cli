import { readFile } from 'fs/promises';
import { getAppBuildGradlePath } from './path.js';

/**
 * Extracts the package name from app/build.gradle(.kts)
 * Looks for namespace or applicationId declarations
 */
export async function getPackageName(projectRoot: string): Promise<string> {
  const appBuildGradlePath = await getAppBuildGradlePath(projectRoot);
  if (!appBuildGradlePath) {
    return '';
  }

  try {
    const data = await readFile(appBuildGradlePath, 'utf-8');
    const content = data.toString();

    // namespace = "com.example.app"
    const namespaceRegex = /^\s*namespace\s*=\s*"(.*?)"/m;
    // applicationId = "com.example.app"
    const appIdRegex = /^\s*applicationId\s*=\s*"(.*?)"/m;

    // Priority: namespace > applicationId
    const namespaceMatch = content.match(namespaceRegex);
    if (namespaceMatch && namespaceMatch[1]) {
      return namespaceMatch[1];
    }

    const appIdMatch = content.match(appIdRegex);
    if (appIdMatch && appIdMatch[1]) {
      return appIdMatch[1];
    }

    return '';
  } catch {
    return '';
  }
}
