import { computed } from 'vue'
import { authStore } from '@/stores/auth'

export function useAuth() {
  return {
    token: computed(() => authStore.token),
    userId: computed(() => authStore.userId),
    userInfo: computed(() => authStore.userInfo),
    username: computed(() => authStore.username),
    isLoggedIn: computed(() => authStore.isLoggedIn),
    login: authStore.login.bind(authStore),
    register: authStore.register.bind(authStore),
    fetchUserInfo: authStore.fetchUserInfo.bind(authStore),
    logout: authStore.logout.bind(authStore),
    persistAuth: authStore.persistAuth.bind(authStore),
    initialize: authStore.initialize.bind(authStore)
  }
}
