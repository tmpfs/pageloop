// see http://vuejs-templates.github.io/webpack for documentation.
var path = require('path')

module.exports = {
  build: {
    env: require('./prod.env'),
    index: path.resolve(__dirname, '../../public/index.html'),
    assetsRoot: path.resolve(__dirname, '../../public'),
    assetsSubDirectory: '',
    assetsPublicPath: '',
    productionSourceMap: true,
    // Gzip off by default as many popular static hosts such as
    // Surge or Netlify already gzip all static assets for you.
    // Before setting to `true`, make sure to:
    // npm install --save-dev compression-webpack-plugin
    productionGzip: false,
    productionGzipExtensions: ['js', 'css'],
    // Run the build command with an extra argument to
    // View the bundle analyzer report after build finishes:
    // `npm run build --report`
    // Set to `true` or `false` to always turn it on or off
    bundleAnalyzerReport: process.env.npm_config_report
  },
  dev: {
    env: require('./dev.env'),
    port: 9090,
    autoOpenBrowser: true,
    assetsSubDirectory: '',
    assetsPublicPath: '/',
    proxyTable: {
      '/assets/': {
        target: 'http://localhost:3577',
        changeOrigin: true
      },
      '/docs/': {
        target: 'http://localhost:3577',
        changeOrigin: true
      },
      '/tools/': {
        target: 'http://localhost:3577',
        changeOrigin: true
      },
      '/template/': {
        target: 'http://localhost:3577',
        changeOrigin: true
      },
      '/rpc/': {
        target: 'http://localhost:3577',
        changeOrigin: true
      },
      '/api/': {
        target: 'http://localhost:3577',
        changeOrigin: true,
        ws: true
      },
      '/apps/': {
        target: 'http://localhost:3577',
        changeOrigin: true
      }
    },
    // CSS Sourcemaps off by default because relative paths are "buggy"
    // with this option, according to the CSS-Loader README
    // (https://github.com/webpack/css-loader#sourcemaps)
    // In our experience, they generally work as expected,
    // just be aware of this issue when enabling this option.
    cssSourceMap: false
  }
}