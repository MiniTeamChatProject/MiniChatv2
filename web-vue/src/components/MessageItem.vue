<template>
  <div :class="['flex', isCurrentUser ? 'justify-end' : 'justify-start']">
    <div :class="['max-w-xs lg:max-w-md px-4 py-2 rounded-lg', isCurrentUser ? 'bg-blue-500 text-white' : 'bg-white text-gray-800']">
      <div v-if="!isCurrentUser" class="text-xs text-gray-500 mb-1">{{ message.username || 'Unknown' }}</div>
      <div class="break-words">{{ message.content }}</div>
      <div :class="['text-xs mt-1', isCurrentUser ? 'text-blue-100' : 'text-gray-400']">
        {{ formatTime(message.created_at) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Message } from '@/types'

defineProps<{
  message: Message
  isCurrentUser: boolean
}>()

function formatTime(timestamp: number): string {
  const date = new Date(timestamp * 1000)
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
}
</script>
