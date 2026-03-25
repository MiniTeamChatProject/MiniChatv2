import { createApp } from 'vue'
import router from './router'
import App from './App.vue'
import { authStore } from './stores/auth'
import './styles/index.css'

// Initialize auth state from storage before mounting
async function bootstrap() {
  await authStore.initialize()
  createApp(App).use(router).mount('#app')
}

bootstrap()
