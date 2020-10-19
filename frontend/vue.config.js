module.exports = {
  devServer: {
    proxy: {
      '^/api': {
        // target: 'http://localhost:9000', // for when running the api locally
        target: 'https://latenightcommits.com/api',
        changeOrigin: true,
        pathRewrite: {
          '^/api/': '',
        },
      },
    },
  },
};
