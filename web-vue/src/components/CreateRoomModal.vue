<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-full max-w-md">
      <h2 class="text-xl font-bold mb-4">Create Room</h2>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Room Name</label>
          <input
            v-model="form.name"
            type="text"
            required
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea
            v-model="form.description"
            rows="3"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Max Members</label>
          <input
            v-model.number="form.max_count"
            type="number"
            min="2"
            max="500"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Room Type</label>
          <select
            v-model="form.type"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option :value="1">Group Chat</option>
            <option :value="2">Voice</option>
            <option :value="3">Video</option>
            <option :value="4">Live</option>
          </select>
        </div>

        <div v-if="error" class="text-red-500 text-sm">{{ error }}</div>

        <div class="flex gap-3 pt-2">
          <button
            type="button"
            @click="$emit('close')"
            class="flex-1 px-4 py-2 border rounded-md hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="flex-1 bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 disabled:opacity-50"
          >
            {{ loading ? 'Creating...' : 'Create' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { roomApi } from '@/api/room'
import type { Room, CreateRoomRequest } from '@/types'

const emit = defineEmits<{
  close: []
  created: [room: Room]
}>()

const loading = ref(false)
const error = ref('')

const form = reactive<CreateRoomRequest>({
  name: '',
  description: '',
  max_count: 100,
  type: 1,
  is_public: true
})

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    const response = await roomApi.createRoom(form)
    emit('created', response.room)
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to create room'
  } finally {
    loading.value = false
  }
}
</script>
