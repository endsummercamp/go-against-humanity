const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const ImageminPlugin = require('imagemin-webpack');

const dev = process.env.NODE_ENV !== 'production';

const postCSSPlugins = [];

if (!dev) {
    postCSSPlugins.push(require('postcss-preset-env')({
        browsers: 'last 2 versions'
    }));
    postCSSPlugins.push(require('cssnano')())
}

module.exports = {
    entry: {
        main: './src/scripts/main.js',
        game: './src/scripts/game.js'
    },
    output: {
        path: path.resolve(__dirname, 'public/scripts')
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, 'src/'),
            '~': path.resolve(__dirname, 'node_modules/')
        }
    },
    mode: dev ? 'development' : 'production',
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader"
                }
            },
            {
                test: /\.less$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                        options: {
                            publicPath: './public/styles/',
                            hmr: dev,
                        },
                    },
                    {
                        loader: 'css-loader'
                    },
                    {
                        loader: 'postcss-loader',
                        options: {
                            plugins: postCSSPlugins
                        }
                    },
                    {
                        loader: 'less-loader'
                    }
                ],
            },
            {
                test: /\.css$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                        options: {
                            publicPath: './public/styles/',
                            hmr: dev,
                        },
                    },
                    {
                        loader: 'css-loader'
                    },
                    {
                        loader: 'postcss-loader',
                        options: {
                            plugins: postCSSPlugins
                        }
                    }
                ],
            },
            {
                test: /\.woff2?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
                use: [{
                    loader: 'url-loader?limit=10000',
                    options: {
                        name: '[name].[ext]',
                        outputPath: 'fonts/'
                    }
                }]
            },
            {
                test: /\.(ttf|eot|svg)(\?v=\d+\.\d+\.\d+)?$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: '[name].[ext]',
                            outputPath: 'fonts/'
                        }
                    }
                ]
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "../styles/main.css"
        }),
        new CopyPlugin([
            {from: 'src/images/', to: 'public/images/'}
        ]),
        new ImageminPlugin({
            bail: false,
            cache: true,
            imageminOptions: {
                plugins: [
                    ["gifsicle", {interlaced: true}],
                    ["jpegtran", {progressive: true}],
                    ["optipng", {optimizationLevel: 3}],
                    ["svgo", {plugins: [{removeViewBox: false}]}]
                ]
            }
        })
    ]
};