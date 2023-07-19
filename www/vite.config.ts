import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  // enable acces through Docker container
  // TODO: make this conditional if /.dockerenv file exists
  // server: {
  //   host: "0.0.0.0"
  // }
})
