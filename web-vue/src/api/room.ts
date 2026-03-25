import { roomHttp } from './http'
import type { Room, CreateRoomRequest, Message, SendMessageRequest } from '@/types'
import { normalizeMessages } from '@/utils/messageNormalizer'

export const roomApi = {
  // Get all rooms list
  listRooms: (page = 1, pageSize = 20) =>
    roomHttp.get<{ rooms: Room[], total: number }>('/v1/rooms', {
      params: { page, page_size: pageSize }
    }),

  // Get user's rooms
  getUserRooms: (userId: number, page = 1, pageSize = 20) =>
    roomHttp.get<{ rooms: Room[], total: number }>(`/v1/users/${userId}/rooms`, {
      params: { page, page_size: pageSize }
    }),

  // Create room
  createRoom: (data: CreateRoomRequest) =>
    roomHttp.post<{ room: Room }>('/v1/rooms', data),

  // Get room info
  getRoom: (roomId: number) =>
    roomHttp.get<{ room: Room }>(`/v1/rooms/${roomId}`),

  // Join room
  joinRoom: (roomId: number) =>
    roomHttp.post(`/v1/rooms/${roomId}/join`, {}),

  // Leave room
  leaveRoom: (roomId: number) =>
    roomHttp.post(`/v1/rooms/${roomId}/leave`, {}),

  // Get message history
  getMessages: async (roomId: number, page = 1, pageSize = 50) => {
    const response = await roomHttp.get<{ messages: any[], total: number }>(
      `/v1/rooms/${roomId}/messages`,
      { params: { page, page_size: pageSize } }
    )
    return {
      messages: normalizeMessages(response.messages),
      total: response.total
    }
  },

  // Send message
  sendMessage: (roomId: number, data: SendMessageRequest) =>
    roomHttp.post<{ message: Message }>(`/v1/rooms/${roomId}/messages`, data),
}
