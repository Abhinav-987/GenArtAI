/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/**/*.templ}", "./**/*.templ"],
  theme: {
    extend: {
      fontSize: {
        '5xl': '3rem',
      },
      colors: {
        pink: {
          light: '#ffafcc',  // Light Pink
          DEFAULT: '#ff4d6d', // Medium Pink
          dark: '#ff1e56',   // Dark Pink
        },
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    theme: ["dark"]
  }
}

