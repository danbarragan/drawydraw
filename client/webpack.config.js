const HtmlWebPackPlugin = require('html-webpack-plugin');

const htmlWebpackPlugin = new HtmlWebPackPlugin({
  template: './src/index.html',
  filename: './index.html',
});
module.exports = {
  entry: [
    '@babel/polyfill',
    './src/index.jsx',
  ],
  output: {
    filename: 'bundled.js',
    path: `${__dirname}/dist`,
  },
  devtool: 'source-map',
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        resolve: { extensions: ['.js', '.jsx'] },
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
        },
      },
      {
        test: /\.css$/,
        use: [
          {
            loader: 'style-loader',
          },
          {
            loader: 'css-loader',
            options: {
              modules: false,
              importLoaders: 1,
              sourceMap: true,
            },
          },
        ],
      },
    ],
  },
  devServer: {
    historyApiFallback: true,
    proxy: { '/api': 'http://localhost:3000' },
  },
  plugins: [
    htmlWebpackPlugin,
  ],
};
