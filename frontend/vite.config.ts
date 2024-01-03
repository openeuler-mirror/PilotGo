import { fileURLToPath, URL } from 'node:url';

import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 8080,
    https: false,
    cors: true,
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8888',
        secure: false, // 如果是https接口，需要配置这个参数
        changeOrigin: true, // 如果接口跨域，需要进行这个参数配置
        // rewrite: path => path.replace(/^\/demo/, '/demo')
      },
      '/plugin/*': {
        target: 'http://localhost:8888',
        secure: false, // 如果是https接口，需要配置这个参数
        changeOrigin: true, // 如果接口跨域，需要进行这个参数配置
        // rewrite: path => path.replace(/^\/demo/, '/demo')
      },
    },
  },
});
