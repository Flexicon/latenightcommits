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
        21: '5.25rem',
        22: '5.5rem',
        23: '5.75rem',
        36: '9rem',
      },
    },
  },
  variants: {},
  plugins: [],
};
