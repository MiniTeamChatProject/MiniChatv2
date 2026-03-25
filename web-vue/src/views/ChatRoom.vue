<template>
  <div class="h-screen flex flex-col bg-gray-100">
    <!-- Header -->
    <div class="bg-white shadow-sm px-4 py-3 flex justify-between items-center">
      <div>
        <h1 class="font-semibold">{{ room?.name || 'Loading...' }}</h1>
        <span v-if="room" class="text-sm text-gray-500">
          {{ room.current_count }}/{{ room.max_count }} members
        </span>
      </div>
      <button
        @click="handleLeave"
        class="text-red-500 hover:text-red-700"
      >
        Leave
      </button>
    </div>

    <!-- Messages -->
    <div class="flex-1 overflow-y-auto p-4">
      <div v-if="loading" class="text-center py-8">Loading messages...</div>

      <MessageList
        v-else
        :messages="messages"
        :current-user-id="Number(authStore.userId)"
      />
    </div>

    <!-- Input -->
    <div class="bg-white border-t p-4">
      <form @submit.prevent="handleSend" class="flex gap-2">
        <input
          v-model="inputMessage"
          type="text"
          placeholder="Type a message..."
          :disabled="!isConnected"
          class="flex-1 px-4 py-2 border rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
        />
        <button
          type="submit"
          :disabled="!inputMessage.trim() || !isConnected"
          class="bg-blue-500 text-white px-6 py-2 rounded-full hover:bg-blue-600 disabled:opacity-50"
        >
          Send
        </button>
      </form>
      <div v-if="!isConnected" class="text-sm text-red-500 mt-2">Connecting...</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authStore } from '@/stores/auth'
import { useWebSocket } from '@/composables/useWebSocket'
import { roomApi } from '@/api/room'
import type { Room } from '@/types'
import MessageList from '@/components/MessageList.vue'

const router = useRouter()
const route = useRoute()

const roomId = computed(() => Number(route.params.id))
const room = ref<Room | null>(null)
const loading = ref(false)
const inputMessage = ref('')

// WebSocket connection - initialize without connecting
const { messages, isConnected, connect, sendMessage, addHistoricalMessages, disconnect } = useWebSocket()

async function loadRoom() {
  loading.value = true
  const userId = authStore.userId
  const username = authStore.username

  // Validate auth state
  if (!userId || !username) {
    console.error('User not authenticated')
    router.push('/login')
    return
  }

  console.log('Loading room:', roomId.value, 'userId:', userId, 'username:', username)

  try {
    // Try to join the room (will handle rejoining if user was in Left status)
    await roomApi.joinRoom(roomId.value)
    console.log('Joined room successfully')
  } catch (error: any) {
    // If already in room, that's fine - continue to load room data
    const errorMessage = error?.message || ''
    if (!errorMessage.includes('user already in room')) {
      console.error('Failed to join room:', error)
      alert(errorMessage || 'Failed to load room')
      router.push('/rooms')
      return
    }
    console.log('User already in room, continuing...')
  }

  try {
    // Get room info
    const response = await roomApi.getRoom(roomId.value)
    room.value = response.room
    console.log('Room loaded:', room.value)

    // Load message history
    const messagesResponse = await roomApi.getMessages(roomId.value)
    addHistoricalMessages(messagesResponse.messages)
    console.log('Messages loaded:', messagesResponse.messages.length)

    // Connect to WebSocket
    connect(roomId.value, Number(userId), username)
  } catch (error: any) {
    console.error('Failed to load room data:', error)
    alert(error?.message || 'Failed to load room data')
    router.push('/rooms')
  } finally {
    loading.value = false
  }
}

function handleSend() {
  if (inputMessage.value.trim() && isConnected.value) {
    sendMessage(inputMessage.value)
    inputMessage.value = ''
  }
}

async function handleLeave() {
  try {
    await roomApi.leaveRoom(roomId.value)
  } catch (error) {
    console.error('Failed to leave room:', error)
  } finally {
    disconnect()
    router.push('/rooms')
  }
}

onMounted(() => {
  loadRoom()
})
</script>
