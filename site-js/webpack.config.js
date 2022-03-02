const path = require("path");
const TerserPlugin = require("terser-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = (env, argv) => {
  // const name = argv.name;
  // const mode = argv.mode || "development";
  return {
    entry: {
      main: "./src/index.tsx",
    },
    module: {
      rules: [
        {
          oneOf: [
            {
              test: /\.css$/i,
              use: [MiniCssExtractPlugin.loader, "css-loader"],
            },
            {
              test: /\.tsx?$/,
              use: [
                {
                  loader: "babel-loader",
                  options: {
                    presets: ["solid"],
                    plugins: [],
                  },
                },
                {
                  loader: "ts-loader",
                  options: {
                    configFile: "tsconfig.json",
                    // onlyCompileBundledFiles: mode === "development",
                    transpileOnly: true,
                  },
                },
              ],
              exclude: /node_modules/,
            },
          ],
        },
      ],
    },
    plugins: [new MiniCssExtractPlugin()],
    resolve: {
      extensions: [".tsx", ".ts", ".jsx", ".js"],
    },
    output: {
      filename: "[name].js",
      path: path.resolve(__dirname, "..", "dist"),
    },
    mode: "production",
    optimization: {
      minimizer: [
        new TerserPlugin({
          extractComments: false,
        }),
      ],
    },
  };
};
