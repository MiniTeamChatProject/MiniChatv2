// User types
export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  nickname: string
}

export interface LoginResponse {
  msg: string
  token: string
}

export interface UserProfile {
  username: string
  nickname: string
}

// Room types
export interface Room {
  id: number
  name: string
  owner_id: number
  type: number
  current_count: number
  max_count: number
  avatar?: string
  description?: string
  tags?: string[]
  is_public: boolean
  status: number
  created_at: number
  updated_at: number
}

export interface CreateRoomRequest {
  name: string
  type?: number
  max_count?: number
  avatar?: string
  description?: string
  tags?: string[]
  is_public?: boolean
}

// Message types
export interface Message {
  id: number
  room_id: number
  user_id: number
  username?: string
  content: string
  type: number
  created_at: number
  // Backend may return camelCase fields
  roomId?: number
  userId?: number
  createdAt?: number
}

export interface SendMessageRequest {
  content: string
  type?: number
}

// WebSocket types
export interface WebSocketMessage {
  type: 'message' | 'system' | 'error' | 'ping' | 'pong'
  data?: {
    id?: number
    room_id?: number
    user_id?: number
    username?: string
    content?: string
    created_at?: number
    message?: string
  }
  error?: string
}
