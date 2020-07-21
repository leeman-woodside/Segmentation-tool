module.exports = {
    optimization: {
        minimize: false,
        splitChunks: {
            chunks: 'async',
            minSize: 20000,
            minRemainingSize: 0,
            maxSize: 0,
            minChunks: 1,
            maxAsyncRequests: 6,
            maxInitialRequests: 4,
            automaticNameDelimiter: '~',
            enforceSizeThreshold: 50000,
            cacheGroups: {
                defaultVendors: {
                test: /[\\/]node_modules[\\/]/,
                priority: -10
                },
                default: {
                minChunks: 2,
                priority: -20,
                reuseExistingChunk: true
                }
            }
        }
    }
};