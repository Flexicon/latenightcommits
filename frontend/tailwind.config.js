module.exports = {
  future: {
    // removeDeprecatedGapUtilities: true,
    // purgeLayersByDefault: true,
  },
  purge: [
    './**/*.html',
    './src/**/*.vue'
  ],
  theme: {
    extend: {
      fontSize: {
        '3xs': '.5rem',
        '2xs': '.625rem',
      },
      spacing: {
        36: '9rem',
      },
    },
  },
  variants: {},
  plugins: [],
};
