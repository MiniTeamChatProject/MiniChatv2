import type { WebSocketMessage } from '@/types'

type MessageHandler = (data: any) => void
type ConnectionHandler = () => void

export class WebSocketClient {
  private ws: WebSocket | null = null
  private url: string
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectTimer: number | null = null
  private messageHandlers: Map<string, MessageHandler> = new Map()
  private onOpenHandlers: ConnectionHandler[] = []
  private onCloseHandlers: ConnectionHandler[] = []

  constructor(private roomId: number, private userId: number, private username: string) {
    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8001'
    this.url = `${wsUrl}/ws/room/${roomId}?user_id=${userId}&username=${encodeURIComponent(username)}`
  }

  connect() {
    this.disconnect()

    this.ws = new WebSocket(this.url)

    this.ws.onopen = () => {
      console.log('WebSocket connected to room:', this.roomId)
      this.reconnectAttempts = 0
      this.onOpenHandlers.forEach(handler => handler())
    }

    this.ws.onmessage = (event) => {
      try {
        const rawMessage: WebSocketMessage = JSON.parse(event.data)
        console.log('WebSocket message received:', rawMessage)

        // Get the handler for this message type
        const handler = this.messageHandlers.get(rawMessage.type)
        if (handler) {
          // Pass the full message structure - let the handler decide how to extract data
          handler(rawMessage)
        } else {
          console.warn('No handler for message type:', rawMessage.type)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error, 'Raw data:', event.data)
      }
    }

    this.ws.onclose = () => {
      console.log('WebSocket closed for room:', this.roomId)
      this.onCloseHandlers.forEach(handler => handler())
      this.reconnect()
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error for room:', this.roomId, error)
    }
  }

  on(type: string, handler: MessageHandler) {
    this.messageHandlers.set(type, handler)
  }

  onOpen(handler: ConnectionHandler) {
    this.onOpenHandlers.push(handler)
  }

  onClose(handler: ConnectionHandler) {
    this.onCloseHandlers.push(handler)
  }

  send(type: string, data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      const payload = JSON.stringify({ type, data })
      console.log('WebSocket sending:', payload)
      this.ws.send(payload)
    } else {
      console.warn('Cannot send message: WebSocket not connected, state:', this.ws?.readyState)
    }
  }

  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  private reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      this.reconnectTimer = window.setTimeout(() => {
        console.log(`Reconnecting... Attempt ${this.reconnectAttempts}`)
        this.connect()
      }, 1000 * this.reconnectAttempts)
    } else {
      console.error('Max reconnection attempts reached')
    }
  }

  get isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}
