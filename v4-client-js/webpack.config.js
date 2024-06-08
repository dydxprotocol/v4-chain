const path = require('path');

const webpack = require('webpack');

module.exports = {
  mode: 'development',
  target: ['web', 'es5'], // Set the target to 'web' for browser environment
  entry: './src/clients/native.ts',
  devtool: 'source-map',
  output: {
    filename: 'v4-native-client.js',
    path: path.resolve(__dirname, '__native__/__ios__'),
    pathinfo: true,
    libraryTarget: 'umd', // Set the libraryTarget to 'umd' for better compatibility
    globalObject: 'this', // Ensure 'window' is not used in the UMD module for universal compatibility
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
      },
    ],
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    fallback: {
      stream: require.resolve('stream-browserify'),
      zlib: require.resolve('browserify-zlib'),
      https: require.resolve('https-browserify'),
      http: require.resolve('stream-http'),
      path: require.resolve('path-browserify'),
      crypto: require.resolve('crypto-browserify'),
    },
  },
  plugins: [
    new webpack.ProvidePlugin({
      Buffer: ['buffer', 'Buffer'],
    }),
    new webpack.DefinePlugin({
      process: {
        env: {
          NODE_ENV: JSON.stringify('production'),
        },
      },
    }),
  ],
};
