<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-2xl font-bold text-center mb-6">MiniChat</h1>

      <div class="flex mb-6 border-b">
        <button
          @click="isLogin = true"
          :class="['flex-1 pb-2 text-center', isLogin ? 'border-b-2 border-blue-500 text-blue-500' : 'text-gray-500']"
        >
          Login
        </button>
        <button
          @click="isLogin = false"
          :class="['flex-1 pb-2 text-center', !isLogin ? 'border-b-2 border-blue-500 text-blue-500' : 'text-gray-500']"
        >
          Register
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input
            v-model="form.username"
            type="text"
            required
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
            v-model="form.password"
            type="password"
            required
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div v-if="!isLogin">
          <label class="block text-sm font-medium text-gray-700 mb-1">Nickname</label>
          <input
            v-model="form.nickname"
            type="text"
            required
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div v-if="error" class="text-red-500 text-sm">{{ error }}</div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 disabled:opacity-50"
        >
          {{ loading ? 'Loading...' : (isLogin ? 'Login' : 'Register') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authStore } from '@/stores/auth'

const router = useRouter()

const isLogin = ref(true)
const loading = ref(false)
const error = ref('')

const form = reactive({
  username: '',
  password: '',
  nickname: ''
})

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    if (isLogin.value) {
      await authStore.login(form.username, form.password)
      await authStore.fetchUserInfo()
      router.push('/rooms')
    } else {
      await authStore.register(form.username, form.password, form.nickname)
      // Auto login after registration
      await authStore.login(form.username, form.password)
      await authStore.fetchUserInfo()
      router.push('/rooms')
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || err.message || 'Operation failed'
  } finally {
    loading.value = false
  }
}
</script>
