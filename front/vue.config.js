module.exports = {
  css: {
    // Enable CSS source maps.
    sourceMap: process.env.NODE_ENV !== 'production'
  },
  publicPath: process.env.NODE_ENV === 'production' ? '/' : '/',
};
