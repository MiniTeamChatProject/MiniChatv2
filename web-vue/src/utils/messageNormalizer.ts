import type { Message } from '@/types'

/**
 * Normalize message from backend format (camelCase) to frontend format (snake_case)
 * Backend returns: { id, roomId, userId, type, content, createdAt }
 * Frontend expects: { id, room_id, user_id, type, content, created_at }
 */
export function normalizeMessage(msg: any): Message {
  return {
    id: msg.id,
    room_id: msg.room_id ?? msg.roomId,
    user_id: msg.user_id ?? msg.userId,
    username: msg.username,
    content: msg.content,
    type: msg.type,
    created_at: msg.created_at ?? msg.createdAt,
  }
}

/**
 * Normalize array of messages
 */
export function normalizeMessages(messages: any[]): Message[] {
  return messages.map(normalizeMessage)
}
