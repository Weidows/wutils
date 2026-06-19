import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  // Wails uses a custom dev server proxy
  server: {
    port: 34115,
    strictPort: true,
  },
  // Prevent vite from obscuring rust errors
  clearScreen: false,
  // Environment variables passed from Wails
  envPrefix: ['VITE_', 'WAILS_'],
})
