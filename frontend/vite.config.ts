import { VitePWA } from 'vite-plugin-pwa';
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import autoprefixer from 'autoprefixer'
import tailwind from 'tailwindcss'

// https://vitejs.dev/config/
export default defineConfig({
	css: {
		postcss: {
			plugins: [
				tailwind(), autoprefixer(),
			],
		},
	},
	plugins:
		[
			vue(),
			VitePWA(
				{
					registerType: 'prompt',
					injectRegister: false,
					pwaAssets: {
						disabled: false,
						config: true,
					},
					manifest: {
						name: 'frontend',
						short_name: 'frontend',
						description: 'frontend',
						theme_color: '#ffffff',
					},
					workbox: {
						globPatterns: ['**/*.{js,css,html,svg,png,ico}'],
						cleanupOutdatedCaches: true,
						clientsClaim: true,
					},
					devOptions: {
						enabled: false,
						navigateFallback: 'index.html',
						suppressWarnings: true,
						type: 'module',
					},
				}
			),

		],
	server: {
		port: 3000,
	},
})
