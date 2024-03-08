// https://github.com/michael-ciniawsky/postcss-load-config

module.exports = {
  "plugins": {
    "postcss-import": {},
    "postcss-url": {},
    // to edit target browsers: use "browserslist" field in package.json
    "autoprefixer": {
      "browsers": ['last 10 Chrome versions', 'last 5 Firefox versions', 'Safari >= 8']
    },
    'postcss-px2rem':{'remUnit':16}
  }
}
