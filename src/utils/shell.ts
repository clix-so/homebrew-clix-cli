import { execa } from 'execa';

export interface ShellResult {
  stdout: string;
  stderr: string;
  exitCode: number;
  success: boolean;
}

export async function runCommand(command: string, args: string[] = [], options: { cwd?: string } = {}): Promise<ShellResult> {
  try {
    const result = await execa(command, args, {
      cwd: options.cwd || process.cwd(),
      shell: true,
      reject: false,
    });

    return {
      stdout: result.stdout,
      stderr: result.stderr,
      exitCode: result.exitCode || 0,
      success: result.exitCode === 0,
    };
  } catch (error: any) {
    return {
      stdout: '',
      stderr: error.message || String(error),
      exitCode: 1,
      success: false,
    };
  }
}

export async function runShellCommand(command: string, options: { cwd?: string } = {}): Promise<ShellResult> {
  try {
    const result = await execa(command, {
      cwd: options.cwd || process.cwd(),
      shell: true,
      reject: false,
    });

    return {
      stdout: result.stdout,
      stderr: result.stderr,
      exitCode: result.exitCode || 0,
      success: result.exitCode === 0,
    };
  } catch (error: any) {
    return {
      stdout: '',
      stderr: error.message || String(error),
      exitCode: 1,
      success: false,
    };
  }
}
