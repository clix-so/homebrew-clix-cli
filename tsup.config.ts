import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['src/cli.tsx'],
  format: ['esm'],
  dts: true,
  clean: true,
  shims: true,
  outDir: 'dist',
  target: 'node18',
  noExternal: [/.*/], // Bundle all dependencies
  banner: {
    js: `#!/usr/bin/env node
import { createRequire } from 'module';
const require = createRequire(import.meta.url);
const __filename = new URL('', import.meta.url).pathname;
const __dirname = new URL('.', import.meta.url).pathname;
`,
  },
  esbuildOptions(options) {
    options.mainFields = ['module', 'main'];
    options.external = ['react-devtools-core'];
    options.platform = 'node';
  },
});
