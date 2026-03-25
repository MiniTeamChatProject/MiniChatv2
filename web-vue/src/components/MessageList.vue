<template>
  <div class="space-y-3">
    <MessageItem
      v-for="message in messages"
      :key="message.id"
      :message="message"
      :is-current-user="message.user_id === currentUserId"
    />
  </div>
</template>

<script setup lang="ts">
import type { Message } from '@/types'
import MessageItem from './MessageItem.vue'

defineProps<{
  messages: Message[]
  currentUserId: number
}>()
</script>
