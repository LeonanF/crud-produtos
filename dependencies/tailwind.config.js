

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../**/*.{html,js}"],
  theme: {
    extend: {colors:{
      'main-light': '#fdf3e7',
      'main-brown': '#7a6a5c'
    },},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

