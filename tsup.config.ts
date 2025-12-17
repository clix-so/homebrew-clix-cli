import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['src/cli.tsx'],
  format: ['esm'],
  dts: true,
  clean: true,
  shims: true,
  outDir: 'dist',
  target: 'node18',
  external: ['react'],
});
