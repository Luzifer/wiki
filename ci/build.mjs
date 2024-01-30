import esbuild from 'esbuild'
import { sassPlugin } from 'esbuild-sass-plugin'
import vuePlugin from 'esbuild-plugin-vue3'

esbuild.build({
  assetNames: '[name]-[hash]',
  bundle: true,
  define: {
    'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'dev'),
  },
  entryPoints: ['src/main.js'],
  legalComments: 'none',
  loader: {
    '.ttf': 'file',
    '.woff2': 'file',
  },
  minify: true,
  outfile: 'frontend/app.js',
  plugins: [
    sassPlugin(),
    vuePlugin(),
  ],
  target: [
    'chrome87',
    'edge87',
    'es2020',
    'firefox84',
    'safari14',
  ],
})
