import { readFile } from 'fs/promises';
import { parseStringPromise } from 'xml2js';

interface ManifestApplication {
  $?: {
    'android:name'?: string;
    [key: string]: any;
  };
}

interface Manifest {
  manifest?: {
    $?: {
      package?: string;
    };
    application?: ManifestApplication[];
  };
}

/**
 * Extracts the application class name from AndroidManifest.xml
 * Returns the value of android:name attribute from the <application> tag
 */
export async function extractApplicationClassName(manifestPath: string): Promise<string> {
  try {
    const data = await readFile(manifestPath, 'utf-8');
    const result: Manifest = await parseStringPromise(data);

    if (result.manifest?.application && result.manifest.application.length > 0) {
      const app = result.manifest.application[0];
      const appName = app.$?.[' android:name'] || app.$?.['android:name'] || '';
      return appName;
    }

    return '';
  } catch {
    return '';
  }
}
