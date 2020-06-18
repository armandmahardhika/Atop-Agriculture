const HtmlWebPackPlugin = require("html-webpack-plugin");
const path = require("path");
module.exports = {
  entry: "./src/index.tsx",
  devtool: "inline-source-map",
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        exclude: /node_modules/,
        use: {
          loader: "ts-loader",
        },
      },
      {
        test: /\.less$/,
        use: [
          { loader: "style-loader" },
          { loader: "css-loader" },
          {
            loader: "less-loader",
            options: { lessOptions: { javascriptEnabled: true } },
          },
        ],
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
    alias: {
      components: path.resolve(__dirname, "src/components/"),
      asset: path.resolve(__dirname, "src/asset/"),
      apis: path.resolve(__dirname, "src/apis/"),
      store: path.resolve(__dirname, "src/store"),
      routers: path.resolve(__dirname, "src/routers"),
      src: path.resolve(__dirname, "src"),
    },
  },
  devServer: {
    port: 4000,
    historyApiFallback: true,
    proxy: {
      "/api": "http://localhost:3001",
      "/apis": "http://localhost:3001",
      "/webapi": "http://localhost:3001",
    },
  },
  plugins: [
    new HtmlWebPackPlugin({
      template: "./public/index.html",
      filename: "./index.html",
    }),
  ],
};
