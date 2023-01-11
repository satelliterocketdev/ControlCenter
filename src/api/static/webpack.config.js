var path = require('path');
var webpack = require('webpack');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

const ENVIRONMENT = process.env.NODE_ENV;
const SRC_DIR = path.resolve(__dirname, 'src');
const BUILD_DIR = __dirname;

const getEntry = () => {
    return SRC_DIR + '/app/index.js';
};

const getOutput = () => {
    return {
        path: BUILD_DIR,
        filename: 'js/app.bundle.js',
        publicPath: "/static"
    }
};

const getModule = () => {

    const BASE_RULES = [
        {
            test: /\.(png|gif|jpe?g|svg)$/,
            loader: 'file-loader?name=/images/[name].[ext]'
        },
        {
            test: /\.(eot|ttf|woff2|woff?)(\?.*$|$)/,
            loader: 'file-loader?name=/fonts/[name].[ext]'
        }
    ];

    const DEV_RULES = [
        {
            test: /\.js?$/,
            exclude: /node_modules/,
            loader: 'babel-loader?presets[]=latest&presets[]=react&presets[]=stage-2'
        },
        {
            test: /\.css$/,
            loader: ExtractTextPlugin.extract({
                fallback: 'style-loader',
                use: [
                    'css-loader?sourceMap'
                ],
            })
        },
        {
            test: /\.(sass|scss)$/,
            loader: ExtractTextPlugin.extract(['css-loader', 'sass-loader'])
        }
    ];
    const PROD_RULES = [
        {
            test: /\.js?$/,
            exclude: /node_modules/,
            loader: 'babel-loader?presets[]=latest&presets[]=react&presets[]=stage-2'
        },
        {
            test: /\.css$/,
            loader: ExtractTextPlugin.extract({
                fallback: 'style-loader',
                use: [
                    'css-loader?sourceMap'
                ],
            }),
        },
        {
            test: /\.(sass|scss)$/,
            loader: ExtractTextPlugin.extract(['css-loader', 'sass-loader'])
        }
    ];

    return {
        rules: ENVIRONMENT === 'production'
            ? PROD_RULES.concat(BASE_RULES)
            : DEV_RULES.concat(BASE_RULES)
    }

};

const watchOptions = () => {
    return {
        ignored: /node_modules/,
        poll: true
    };
};

const getResolve = () => {
    return {
        extensions: [".js", ".json", ".css"],
        alias: {
            app: SRC_DIR + '/app',
            fonts: SRC_DIR + '/fonts',
            img: SRC_DIR + '/images',
            libs: SRC_DIR + '/libs',
            svg: SRC_DIR + '/svg'
        }
    }
};

const getDevTools = () => {
    return 'source-map'
};

const getPlugins = () => {

    const BASE_PLUGINS = [
        new webpack.EnvironmentPlugin({NODE_ENV: ENVIRONMENT || 'development'}),
        new webpack.ProvidePlugin({
            'React': 'react',
            'ReactDOM': 'react-dom'
        }),
        new webpack.NoEmitOnErrorsPlugin(),
        new ExtractTextPlugin({
            filename: 'css/app.css',
            allChunks: true
        }),
    ];

    const DEV_PLUGINS = [
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: `${SRC_DIR}/index.html`
        })
    ];

    const PROD_PLUGINS = [
        new HtmlWebpackPlugin({
            filename: 'index.html',
            hash: true,
            template: `${SRC_DIR}/index.html`
        })

    ];
    return BASE_PLUGINS.concat(ENVIRONMENT === 'production' ? PROD_PLUGINS : DEV_PLUGINS);
};


module.exports = {
    entry: getEntry(),
    output: getOutput(),
    module: getModule(),
    watchOptions: watchOptions(),
    devtool: getDevTools(),
    resolve: getResolve(),
    plugins: getPlugins()
};
