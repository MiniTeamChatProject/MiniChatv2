<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Header -->
    <div class="bg-white shadow-sm">
      <div class="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
        <h1 class="text-xl font-bold">MiniChat</h1>
        <div class="flex items-center gap-4">
          <span class="text-gray-600">{{ authStore.userInfo?.nickname || authStore.userInfo?.username }}</span>
          <button
            @click="logout"
            class="text-red-500 hover:text-red-700"
          >
            Logout
          </button>
        </div>
      </div>
    </div>

    <!-- Room List -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-lg font-semibold">Rooms</h2>
        <button
          @click="showCreateModal = true"
          class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600"
        >
          Create Room
        </button>
      </div>

      <div v-if="loading" class="text-center py-8">Loading...</div>

      <div v-else-if="rooms.length === 0" class="text-center py-8 text-gray-500">
        No rooms available. Create one to get started!
      </div>

      <div v-else class="grid gap-4">
        <RoomCard
          v-for="room in rooms"
          :key="room.id"
          :room="room"
          :is-joined="joinedRoomIds.has(String(room.id))"
          @join="handleJoinRoom(room)"
        />
      </div>
    </div>

    <!-- Create Room Modal -->
    <CreateRoomModal
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @created="handleRoomCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { authStore } from '@/stores/auth'
import { roomApi } from '@/api/room'
import type { Room } from '@/types'
import RoomCard from '@/components/RoomCard.vue'
import CreateRoomModal from '@/components/CreateRoomModal.vue'

const router = useRouter()

const rooms = ref<Room[]>([])
const joinedRoomIds = ref<Set<string>>(new Set())
const loading = ref(false)
const showCreateModal = ref(false)

async function loadRooms() {
  loading.value = true
  try {
    const response = await roomApi.listRooms()
    rooms.value = response.rooms
    // Load user's joined rooms
    await loadJoinedRooms()
  } catch (error) {
    console.error('Failed to load rooms:', error)
  } finally {
    loading.value = false
  }
}

async function loadJoinedRooms() {
  try {
    const response = await roomApi.getUserRooms(Number(authStore.userId))
    joinedRoomIds.value = new Set(response.rooms.map(r => String(r.id)))
  } catch (error) {
    console.error('Failed to load joined rooms:', error)
  }
}

async function handleJoinRoom(room: Room) {
  // Just navigate to the room, ChatRoom component will handle joining
  router.push(`/room/${room.id}`)
}

function handleRoomCreated(room: Room) {
  showCreateModal.value = false
  rooms.value.unshift(room)
  joinedRoomIds.value.add(String(room.id))
  // Auto enter the newly created room
  router.push(`/room/${room.id}`)
}

function logout() {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadRooms()
})

// Reload data when returning from other pages
onActivated(() => {
  loadRooms()
})
</script>
