import { promises as fs } from 'fs';
import path from 'path';
import { execa } from 'execa';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

interface AutomationResult {
  success: boolean;
  message: string;
  data: Record<string, any>;
}

/**
 * Ensures Ruby is installed on the system
 * @throws Error if Ruby cannot be installed
 */
export async function ensureRuby(): Promise<void> {
  try {
    await execa('ruby', ['--version']);
  } catch {
    console.log('⏳ Ruby is required but not found. Installing Ruby...');

    try {
      await execa('brew', ['--version']);
    } catch {
      throw new Error('ruby and homebrew are both not installed, please install ruby manually');
    }

    // Install Ruby using Homebrew
    await execa('brew', ['install', 'ruby'], { stdio: 'inherit' });
    console.log('✅ Ruby installed successfully');
  }
}

/**
 * Ensures the xcodeproj gem is installed
 * @throws Error if gem cannot be installed
 */
export async function ensureXcodeproj(): Promise<void> {
  const minVersion = '>= 1.25.2';

  // Check if a compatible version is installed
  try {
    const { stdout } = await execa('gem', ['query', '-i', '-n', '^xcodeproj$', '-v', minVersion]);
    if (stdout.trim() === 'true') {
      return; // Compatible version already present
    }
  } catch {
    // Not installed or wrong version
  }

  // Print the currently loaded version for diagnostics (best-effort)
  try {
    await execa('ruby', ['-e', 'begin; require \'xcodeproj\'; puts \'Detected xcodeproj \' + Xcodeproj::VERSION; rescue; end']);
  } catch {
    // Ignore
  }

  console.log(`⏳ Installing/upgrading xcodeproj gem to ${minVersion}...`);

  // Try to install or upgrade without sudo first
  const installArgs = ['install', 'xcodeproj', '-v', minVersion];
  try {
    await execa('gem', installArgs, { stdio: 'inherit' });
  } catch {
    console.log('🔐 Regular gem installation failed. Trying with sudo...');
    console.log('You may be prompted for your password.');

    const sudoArgs = ['gem', ...installArgs];
    try {
      await execa('sudo', sudoArgs, { stdio: 'inherit' });
    } catch (error) {
      throw new Error(
        `failed to install required xcodeproj gem version (${minVersion}): ${error instanceof Error ? error.message : String(error)}\n` +
        `This gem version is required to parse newer Xcode projects that use PBXFileSystemSynchronizedRootGroup.\n` +
        `You can also install it manually with: sudo gem install xcodeproj -v '${minVersion}'`
      );
    }
  }

  // Re-check to confirm
  try {
    const { stdout } = await execa('gem', ['query', '-i', '-n', '^xcodeproj$', '-v', minVersion]);
    if (stdout.trim() !== 'true') {
      throw new Error(`xcodeproj gem did not meet version requirement ${minVersion} even after install; please ensure RubyGems is configured and try again`);
    }
  } catch (error) {
    throw new Error(`xcodeproj gem verification failed: ${error instanceof Error ? error.message : String(error)}`);
  }

  console.log('✅ xcodeproj gem is up to date');
}

/**
 * Finds the Xcode project file in the current or nearby directories
 * @returns Path to the .xcodeproj file
 * @throws Error if no project is found
 */
export async function findXcodeProject(): Promise<string> {
  const cwd = process.cwd();

  // Find .xcodeproj files in current directory
  let entries = await fs.readdir(cwd, { withFileTypes: true });
  let matches = entries.filter(e => e.isDirectory() && e.name.endsWith('.xcodeproj')).map(e => path.join(cwd, e.name));

  if (matches.length > 0) {
    return matches[0];
  }

  // Look one level up
  const parentDir = path.dirname(cwd);
  entries = await fs.readdir(parentDir, { withFileTypes: true });
  matches = entries.filter(e => e.isDirectory() && e.name.endsWith('.xcodeproj')).map(e => path.join(parentDir, e.name));

  if (matches.length > 0) {
    return matches[0];
  }

  // Look in iOS directory if exists
  const iosDir = path.join(cwd, 'ios');
  try {
    await fs.access(iosDir);
    entries = await fs.readdir(iosDir, { withFileTypes: true });
    matches = entries.filter(e => e.isDirectory() && e.name.endsWith('.xcodeproj')).map(e => path.join(iosDir, e.name));

    if (matches.length > 0) {
      return matches[0];
    }
  } catch {
    // iOS directory doesn't exist
  }

  throw new Error('no .xcodeproj file found');
}

/**
 * Configures the Xcode project with App Groups and framework dependencies
 * @param projectID The Clix project ID
 * @param verbose Enable verbose output
 * @param dryRun Perform a dry run without making changes
 * @throws Error if configuration fails
 */
export async function configureXcodeProject(
  projectID: string,
  verbose: boolean = false,
  dryRun: boolean = false
): Promise<void> {
  // Ensure Ruby is installed
  await ensureRuby();

  // Ensure xcodeproj gem is installed
  await ensureXcodeproj();

  // Find Xcode project
  const xcodeProjectPath = await findXcodeProject();
  console.log(`📁 Found Xcode project: ${xcodeProjectPath}`);

  // Create app group ID
  const appGroupID = `group.clix.${projectID}`;

  // The Ruby script is in the scripts subdirectory
  const scriptPath = path.join(__dirname, 'scripts', 'configure_xcode_project.rb');

  try {
    await fs.access(scriptPath);
  } catch {
    throw new Error(`ruby script not found at: ${scriptPath}`);
  }

  // Build command arguments
  const args = [
    scriptPath,
    '--project-path', xcodeProjectPath,
    '--app-group-id', appGroupID,
  ];

  if (verbose) {
    args.push('--verbose');
  }

  if (dryRun) {
    // In dry-run mode, just show what would be done
    console.log('🔍 DRY RUN - Would execute:');
    console.log(`   ruby ${args.join(' ')}`);
    console.log(`   This would configure App Group '${appGroupID}'`);
    console.log('   This would add Clix framework to NotificationServiceExtension');
    return;
  }

  console.log('🔄 Configuring Xcode project...');

  // Run the Ruby script
  try {
    const result = await execa('ruby', args);
    const stdout = result.stdout;

    // Parse JSON result from successful execution
    const parsedResult: AutomationResult = JSON.parse(stdout);

    if (!parsedResult.success) {
      console.error(`Error from configuration script: ${parsedResult.message}`);
      throw new Error(`script reported failure: ${parsedResult.message}`);
    }

    console.log('✅ Xcode project configured successfully!');
    console.log('   - App Groups capability added');
    console.log('   - Background Modes enabled: fetch, remote-notification');
    console.log('   - Clix framework added to NotificationServiceExtension (if present)');
  } catch (error: any) {
    // Try to parse JSON error from stdout
    if (error.stdout) {
      try {
        const parsedResult: AutomationResult = JSON.parse(error.stdout);
        if (!parsedResult.success) {
          throw new Error(parsedResult.message);
        }
      } catch (parseError) {
        // Could not parse JSON, show raw error
      }
    }

    // If we couldn't get a structured error, return a generic one
    if (error.stderr) {
      throw new Error(`failed to run Ruby script: ${error.message}\n\n--- Script Error ---\n${error.stderr}`);
    }
    throw new Error(`failed to run Ruby script: ${error.message}`);
  }
}
