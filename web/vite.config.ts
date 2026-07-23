import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import basicSsl from '@vitejs/plugin-basic-ssl'

export default defineConfig(({ command }) => ({
  plugins: [
    react(),
    // Self-signed HTTPS so getUserMedia works on mobile LAN during dev.
    // Production HTTPS is handled by the reverse proxy (nginx/Caddy) — this never runs in builds.
    // Disabled locally: this machine's TLS stack rejects the self-signed cert (ERR_SSL_PROTOCOL_ERROR).
    ...(command === 'serve' && false ? [basicSsl()] : []),
  ],
  server: {
    port: 5173,
    host: true, // expose to LAN for mobile testing
    proxy: {
      '/api': {
        target: 'http://localhost:3010',
        changeOrigin: true,
      },
    },
  },
}))
