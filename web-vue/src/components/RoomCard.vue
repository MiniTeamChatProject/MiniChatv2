<template>
  <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
    <div class="flex justify-between items-start">
      <div class="flex-1">
        <h3 class="font-semibold text-lg">{{ room.name }}</h3>
        <p v-if="room.description" class="text-gray-600 text-sm mt-1">{{ room.description }}</p>
        <div class="flex items-center gap-4 mt-3 text-sm text-gray-500">
          <span>Members: {{ room.current_count }}/{{ room.max_count }}</span>
          <span>Type: {{ getRoomType(room.type) }}</span>
        </div>
      </div>
      <button
        @click="$emit('join')"
        class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 text-sm"
      >
        {{ isJoined ? 'Enter' : 'Join' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Room } from '@/types'

defineProps<{
  room: Room
  isJoined?: boolean
}>()

defineEmits<{
  join: []
}>()

function getRoomType(type: number): string {
  const types: Record<number, string> = {
    1: 'Group Chat',
    2: 'Voice',
    3: 'Video',
    4: 'Live'
  }
  return types[type] || 'Unknown'
}
</script>
