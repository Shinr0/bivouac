import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// En production, l'Ingress route /api vers le BFF sur la même origine.
// En développement, ce proxy reproduit ce comportement.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: process.env.FORGE_BFF_URL ?? 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
