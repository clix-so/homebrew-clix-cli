#!/usr/bin/env node
import { readFileSync, writeFileSync, chmodSync } from 'fs';
import { join } from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const cliPath = join(__dirname, '..', 'dist', 'cli.js');
let content = readFileSync(cliPath, 'utf-8');

// Remove any existing shebangs
const lines = content.split('\n');
const filteredLines = lines.filter(line => !line.startsWith('#!'));

// Add single shebang at the top
content = `#!/usr/bin/env node\n${filteredLines.join('\n')}`;

writeFileSync(cliPath, content);
chmodSync(cliPath, 0o755);
console.log('✅ Added shebang to dist/cli.js and made it executable');
