/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./src/**/*.{js,ts,jsx,tsx}",
    "./index.html",
  ],
  theme: {
    extend: {
      colors: {
        wutils: {
          bg: '#1b2636',
          surface: '#243044',
          accent: '#4a90d9',
          success: '#4ade80',
          warning: '#fbbf24',
          error: '#f87171',
        },
      },
    },
  },
  plugins: [],
}
