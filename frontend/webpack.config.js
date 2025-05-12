const path = require('path')
const TerserPlugin = require('terser-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer')

const devFrontendPort = 4114
const devApiPort = 8093

const isProduction = process.env.ENV === 'production'

const chunks = {
  miniapp: path.join(__dirname, './src/app.js')
}

const lightColorScheme = true
const nightColorScheme = true

const defaultLightScheme =
  '--tg-theme-accent-text-color: #168dcd; --tg-color-scheme: light; --tg-theme-bg-color: #ffffff; --tg-theme-bottom-bar-bg-color: #ffffff; --tg-theme-button-color: #40a7e3; --tg-theme-button-text-color: #ffffff; --tg-theme-destructive-text-color: #d14e4e; --tg-theme-header-bg-color: #ffffff; --tg-theme-hint-color: #999999; --tg-theme-link-color: #168dcd; --tg-theme-secondary-bg-color: #f1f1f1; --tg-theme-section-bg-color: #ffffff; --tg-theme-section-header-text-color: #168dcd; --tg-theme-section-separator-color: #e7e7e7; --tg-theme-subtitle-text-color: #999999; --tg-theme-text-color: #000000; --tg-viewport-height: 650px; --tg-viewport-stable-height: 650px; --deprecated-separator-color: #D9D9D9; --second-button-color: rgb(64 167 227 / 10%); --text-confirm-color: #4ECC5E; --button-confirm-color: #4ECC5E; --button-main-disabled-color: #E9E8E8; --text-main-disabled-color: #BABABA; --button-destructive-color: rgb(209 78 78 / 10%); --highlight-default: rgba(0, 0, 0, 0.05); --quick-menu-background: #FFFFFF; --quick-menu-foreground: #222222; --toast-background: rgba(40, 48, 57, 0.92); --text-overlay: #FFFFFF; --toast-link: #85CAFF; --tertiary-fill-background: rgba(122, 122, 122, 0.12); --quaternary-fill-background: rgba(122, 122, 122, 0.08); --separator-non-opaque-color: rgba(153, 153, 153, 0.36); --accent-orange: #F68136; --segmented-control-active-background: #FFFFFF; --tooltip-background: rgba(0, 0, 0, 0.8); --tg-safe-area-inset-top: 0px; --tg-safe-area-inset-bottom: 0px; --tg-safe-area-inset-left: 0px; --tg-safe-area-inset-right: 0px; --tg-content-safe-area-inset-top: 0px; --tg-content-safe-area-inset-bottom: 0px; --tg-content-safe-area-inset-left: 0px; --tg-content-safe-area-inset-right: 0px; background: rgb(241, 241, 241);'
const defaultNightScheme =
  '--tg-theme-subtitle-text-color: #98989e; --tg-theme-section-separator-color: #545458; --tg-theme-accent-text-color: #3e88f7; --tg-theme-destructive-text-color: #eb5545; --tg-theme-header-bg-color: #1a1a1a; --tg-theme-bottom-bar-bg-color: #1d1d1d; --tg-theme-button-color: #3e88f7; --tg-theme-link-color: #3e88f7; --tg-theme-hint-color: #98989e; --tg-theme-text-color: #ffffff; --tg-theme-button-text-color: #ffffff; --tg-theme-secondary-bg-color: #1c1c1d; --tg-theme-section-bg-color: #2c2c2e; --tg-color-scheme: dark; --tg-theme-bg-color: #000000; --tg-theme-section-header-text-color: #8d8e93; --tg-viewport-height: 640px; --tg-viewport-stable-height: 640px; --tg-safe-area-inset-top: 0px; --tg-safe-area-inset-bottom: 0px; --tg-safe-area-inset-left: 0px; --tg-safe-area-inset-right: 0px; --tg-content-safe-area-inset-top: 0px; --tg-content-safe-area-inset-bottom: 0px; --tg-content-safe-area-inset-left: 0px; --tg-content-safe-area-inset-right: 0px; background: #000'
const defaultDarkScheme =
  '--tg-theme-accent-text-color: #79e8da; --tg-color-scheme: dark; --tg-theme-bg-color: #282e33; --tg-theme-bottom-bar-bg-color: #282e33; --tg-theme-button-color: #3fc1b0; --tg-theme-button-text-color: #ffffff; --tg-theme-destructive-text-color: #f57474; --tg-theme-header-bg-color: #282e33; --tg-theme-hint-color: #82868a; --tg-theme-link-color: #4be1c3; --tg-theme-secondary-bg-color: #313b43; --tg-theme-section-bg-color: #282e33; --tg-theme-section-header-text-color: #4be1c3; --tg-theme-section-separator-color: #242a2e; --tg-theme-subtitle-text-color: #82868a; --tg-theme-text-color: #f5f5f5; --tg-viewport-height: 650px; --tg-viewport-stable-height: 650px; --deprecated-separator-color: #000000; --second-button-color: rgb(63 193 176 / 10%); --text-confirm-color: #61BD67; --button-confirm-color: #61BD67; --button-main-disabled-color: #3C3C3E; --text-main-disabled-color: #606060; --button-destructive-color: rgb(245 116 116 / 10%); --highlight-default: rgba(255, 255, 255, 0.05); --quick-menu-background: #282829; --quick-menu-foreground: #FFFFFF; --toast-background: rgba(40, 40, 40, 0.96); --text-overlay: #FFFFFF; --toast-link: #6CB7FF; --tertiary-fill-background: rgba(123, 123, 123, 0.24); --quaternary-fill-background: rgba(122, 122, 122, 0.18); --separator-non-opaque-color: #000000; --accent-orange: #E58943; --segmented-control-active-background: #636366; --tooltip-background: #3a3a3c; --tg-safe-area-inset-top: 0px; --tg-safe-area-inset-bottom: 0px; --tg-safe-area-inset-left: 0px; --tg-safe-area-inset-right: 0px; --tg-content-safe-area-inset-top: 0px; --tg-content-safe-area-inset-bottom: 0px; --tg-content-safe-area-inset-left: 0px; --tg-content-safe-area-inset-right: 0px; background: #313b43;'

const defaultLightTheme =
  '{"bg_color":"#ffffff","text_color":"#000000","button_text_color":"#ffffff","secondary_bg_color":"#efeff4","section_bg_color":"#ffffff","link_color":"#007aff","button_color":"#007aff","section_separator_color":"#c8c7cc","destructive_text_color":"#ff3b30","bottom_bar_bg_color":"#f2f2f2","accent_text_color":"#007aff","subtitle_text_color":"#8e8e93","header_bg_color":"#f8f8f8","section_header_text_color":"#6d6d72","hint_color":"#8e8e93"}'
const defaultNightTheme =
  '{"bg_color":"#000000","text_color":"#ffffff","button_text_color":"#ffffff","secondary_bg_color":"#1c1c1d","section_bg_color":"#2c2c2e","link_color":"#3e88f7","button_color":"#3e88f7","section_separator_color":"#545458","destructive_text_color":"#eb5545","bottom_bar_bg_color":"#1d1d1d","accent_text_color":"#3e88f7","subtitle_text_color":"#98989e","header_bg_color":"#1a1a1a","section_header_text_color":"#8d8e93","hint_color":"#98989e"}'
const defaultDarkTheme =
  '{"bg_color":"#18222d","button_color":"#2ea6ff","bottom_bar_bg_color":"#21303f","hint_color":"#dbf5ff","subtitle_text_color":"#dbf5ff","header_bg_color":"#21303f","text_color":"#ffffff","destructive_text_color":"#ff6767","button_text_color":"#ffffff","section_bg_color":"#21303f","link_color":"#2ea6ff","section_header_text_color":"#83898e","secondary_bg_color":"#18222d","section_separator_color":"#384656","accent_text_color":"#2ea6ff"}'

const currentColorScheme = lightColorScheme ? 'light' : 'dark'
const currentColorStyle = `style="${
  lightColorScheme
    ? defaultLightScheme
    : nightColorScheme
    ? defaultNightScheme
    : defaultDarkScheme
}"`
const currentColorTheme = lightColorScheme
  ? defaultLightTheme
  : nightColorScheme
  ? defaultNightTheme
  : defaultDarkTheme

const debugAUTH =
  'query_id=AAEdao8PAAAAAB1qjw8BR1jZ&user=%7B%22id%22%3Ar242342423%2C%22first_name%22%3A%22Andrei%22%2C%22last_name%22%3A%22Petkun%22%2C%22username%22%3A%22avpetkun%22%2C%22language_code%22%3A%22nb%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2FJDg-ZPMWDdEqCXkfAyCIvhmCl1MEknt5iGOdYKkJxbM.svg%22%7D&auth_date=1640514145&signature=WJEHHFBWIEFBHHWIEBFIWBHFIBWEFILBJFIBFWELIFBWFI&hash=wefjeowfijwofjewofjewofjewofjewofjweofiew'
const debugTgUser = 260043021

const prodHeadAnalytics = `<script async src="https://tganalytics.xyz/index.js" onload="initAnalytics()" type="text/javascript"></script>`

const prodAnalytics = `<script>
    function initAnalytics() {
      window.telegramAnalytics.init({
        token: 'wefewfwefewfewfewfwefwefwefweewf',
        appName: 'ThePrimeBot'
      });
    }
</script>`

module.exports = {
  mode: isProduction ? 'production' : 'development',
  entry: chunks,
  output: {
    path: path.resolve(__dirname, 'static'),
    filename: isProduction ? '[name].[chunkhash].js' : '[name].js',
    chunkFilename: isProduction ? '[name].[chunkhash].js' : '[name].js',
    publicPath: '/static/',
    clean: true
  },
  optimization: {
    // splitChunks: {
    //   chunks: 'all',
    //   name: 'runtime'
    // }
    // splitChunks: {
    //   cacheGroups: {
    //     common: {
    //       name: 'common',
    //       chunks: (chunk) => true,
    //       reuseExistingChunk: true,
    //       priority: 1,
    //       minChunks: 1,
    //       minSize: 0,
    //       test: /[\\/]node_modules[\\/]/
    //     }
    //   }
    // }
  },
  module: {
    rules: [
      {
        test: /\.js?$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
        options: {
          babelrc: false,
          cacheDirectory: !isProduction,
          presets: [['@babel/preset-env', { targets: { node: 'current' } }]]
        }
      },
      {
        test: /\.svelte$/,
        use: {
          loader: 'svelte-loader',
          options: {
            emitCss: true
          }
        }
      },
      {
        test: /\.(css|sass|scss)$/,
        use: [
          MiniCssExtractPlugin.loader,
          ...(!isProduction ? ['cache-loader'] : []),
          {
            loader: 'css-loader',
            options: {
              importLoaders: 1
            }
          },
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  require('postcss-import'),
                  require('postcss-for'),
                  require('postcss-nested'),
                  require('postcss-css-variables')({ preserve: true }),
                  require('postcss-simple-vars')({ silent: true, keep: true }),
                  require('postcss-preset-env')({
                    stage: false,
                    features: {
                      'custom-properties': false
                    }
                  }),
                  require('postcss-csso')
                ]
              }
            }
          }
        ]
      },
      {
        test: /\.(png|jpg|jpeg|gif|svg|ico|woff|woff2?|ttf|eot|ogg|wav)$/,
        type: 'asset/resource'
      }
    ]
  },
  resolve: {
    modules: [
      path.join(__dirname, 'src'),
      path.join(__dirname, 'node_modules')
    ],
    alias: {
      '~': path.join(__dirname, 'src'),
      '@': path.join(__dirname, 'assets'),
      svelte: path.resolve('node_modules', 'svelte')
    },
    extensions: ['.js', '.vue', '.css', '.scss', '.json', '.mjs', '.svelte'],
    conditionNames: ['svelte', 'import', 'node', 'default']
  },
  devServer: {
    contentBase: path.resolve(__dirname, 'static'),
    historyApiFallback: true,
    hot: true,
    overlay: true,
    noInfo: false,
    port: devFrontendPort,
    host: process.env.PUBLIC ? '0.0.0.0' : 'localhost',
    proxy: [{ path: '/api', target: `http://localhost:${devApiPort}/` }],
    stats: {
      assets: false,
      colors: true,
      timings: true,
      version: false,
      hash: false,
      chunks: true,
      chunkModules: false,
      modules: false
    }
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: isProduction ? '[name].[chunkhash].css' : '[name].css'
    }),
    ...Object.keys(chunks).map(
      (chunk) =>
        new HtmlWebpackPlugin({
          chunks: [chunk],
          filename: `${chunk}.html`,
          template: `${chunk}.html`,
          templateParameters: {
            HTML: isProduction ? '' : currentColorStyle,
            HEAD: isProduction
              ? prodHeadAnalytics
              : `<script>
                  window.AUTH = '${debugAUTH}';
                  window.TG_USER = ${debugTgUser};
                  window.SCHEME = '${currentColorScheme}';
                  window.THEME = ${currentColorTheme};
                  window.DEBUG = true;
                  window.LANG = 'en';
                </script>`,
            BODY: isProduction ? prodAnalytics : ''
          }
        })
    )
  ]
}

if (isProduction) {
  module.exports.optimization.minimize = true
  module.exports.optimization.minimizer = [
    new TerserPlugin({
      cache: true,
      parallel: true,
      sourceMap: false,
      extractComments: true,
      terserOptions: {
        compress: {
          drop_console: true
        }
      }
    })
  ]
}

if (process.env.ANALYZE) {
  module.exports.plugins.push(
    new BundleAnalyzerPlugin({
      analyzerMode: 'static'
    })
  )
}
