import { ref, onUnmounted, computed } from 'vue'
import { WebSocketClient } from '@/api/websocket'
import { getUsername, cacheUsername } from '@/utils/userCache'
import type { Message } from '@/types'

export function useWebSocket() {
  const messages = ref<Message[]>([])
  const isConnected = ref(false)
  let wsClient: WebSocketClient | null = null
  let currentRoomId: number | null = null
  let currentUserId: number | null = null
  let currentUsername: string = ''

  // Set to track message IDs for deduplication
  const messageIds = ref<Set<number>>(new Set())

  function connect(roomId: number, userId: number, username: string) {
    // Validate parameters
    if (!roomId || !userId || !username) {
      console.warn('Cannot connect: missing parameters', { roomId, userId, username })
      return
    }

    // Disconnect existing connection if any
    disconnect()

    currentRoomId = roomId
    currentUserId = userId
    currentUsername = username

    // Cache current user's username
    cacheUsername(userId, username)

    wsClient = new WebSocketClient(roomId, userId, username)

    wsClient.onOpen(() => {
      isConnected.value = true
      console.log('WebSocket connected to room:', roomId)
    })

    wsClient.onClose(() => {
      isConnected.value = false
      console.log('WebSocket disconnected from room:', roomId)
    })

    wsClient.on('message', (data) => {
      // Handle message data - check if nested or flat
      const msgData = data.data || data

      const newMessage: Message = {
        id: msgData.id || Date.now(),
        room_id: msgData.room_id || roomId,
        user_id: msgData.user_id || userId,
        username: msgData.username || getUsername(msgData.user_id || userId),
        content: msgData.content || '',
        type: msgData.type || 1,
        created_at: msgData.created_at || Date.now() / 1000
      }

      // Cache username from message
      if (newMessage.username && newMessage.user_id) {
        cacheUsername(newMessage.user_id, newMessage.username)
      }

      // Deduplicate: only add if message ID not seen
      if (newMessage.id && !messageIds.value.has(newMessage.id)) {
        messageIds.value.add(newMessage.id)
        messages.value.push(newMessage)
      } else if (!newMessage.id) {
        // If no ID, still add the message
        messages.value.push(newMessage)
      }
    })

    wsClient.on('system', (data) => {
      console.log('System message:', data)
    })

    wsClient.on('error', (data) => {
      console.error('WebSocket error:', data)
    })

    wsClient.connect()
  }

  function sendMessage(content: string) {
    if (wsClient?.isConnected && content.trim()) {
      wsClient.send('message', { content })
    }
  }

  function addHistoricalMessages(historicalMsgs: Message[]) {
    // Add historical messages and track their IDs
    for (const msg of historicalMsgs) {
      if (msg.id && !messageIds.value.has(msg.id)) {
        messageIds.value.add(msg.id)

        // Ensure historical messages have username
        const messageWithUsername: Message = {
          ...msg,
          username: msg.username || getUsername(msg.user_id)
        }

        // Cache username from historical messages
        if (messageWithUsername.username && msg.user_id) {
          cacheUsername(msg.user_id, messageWithUsername.username)
        }

        messages.value.push(messageWithUsername)
      }
    }
  }

  function clearMessages() {
    messages.value = []
    messageIds.value.clear()
  }

  function disconnect() {
    if (wsClient) {
      wsClient.disconnect()
      wsClient = null
    }
    isConnected.value = false
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    messages,
    isConnected,
    connect,
    sendMessage,
    addHistoricalMessages,
    clearMessages,
    disconnect
  }
}
