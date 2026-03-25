package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"room/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/websocket"
)

// WebSocketMessage WebSocket 消息格式
type WebSocketMessage struct {
	Type    string      `json:"type"`    // 消息类型: message, system, error
	Data    interface{} `json:"data"`    // 消息数据
	Error   string      `json:"error"`   // 错误信息（可选）
}

// WebSocketClient WebSocket 客户端
type WebSocketClient struct {
	id       string
	roomID   int64
	userID   int64
	username string
	conn     *websocket.Conn
	send     chan *WebSocketMessage
	hub      *WebSocketHub
	ctx      context.Context
	cancel   context.CancelFunc
}

// WebSocketHub WebSocket 连接管理中心
type WebSocketHub struct {
	// 房间ID -> 客户端列表
	rooms map[int64]map[*WebSocketClient]bool

	// 注册客户端
	register chan *WebSocketClient

	// 注销客户端
	unregister chan *WebSocketClient

	// 广播消息到房间
	broadcast chan *RoomMessage

	// 互斥锁
	mu sync.RWMutex

	log *log.Helper
}

// RoomMessage 房间消息
type RoomMessage struct {
	RoomID  int64
	Message *WebSocketMessage
}

// NewWebSocketHub 创建 WebSocket Hub
func NewWebSocketHub(logger log.Logger) *WebSocketHub {
	hub := &WebSocketHub{
		rooms:      make(map[int64]map[*WebSocketClient]bool),
		register:   make(chan *WebSocketClient),
		unregister: make(chan *WebSocketClient),
		broadcast:  make(chan *RoomMessage, 256),
		log:        log.NewHelper(logger),
	}

	go hub.run()

	return hub
}

// run 运行 Hub
func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case roomMsg := <-h.broadcast:
			h.broadcastToRoom(roomMsg)
		}
	}
}

// registerClient 注册客户端
func (h *WebSocketHub) registerClient(client *WebSocketClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.roomID] == nil {
		h.rooms[client.roomID] = make(map[*WebSocketClient]bool)
	}
	h.rooms[client.roomID][client] = true

	h.log.Infof("Client %s joined room %d", client.id, client.roomID)
}

// unregisterClient 注销客户端
func (h *WebSocketHub) unregisterClient(client *WebSocketClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.roomID]; ok {
		if clients[client] {
			delete(clients, client)
			close(client.send)
			// 取消 context
			if client.cancel != nil {
				client.cancel()
			}
			h.log.Infof("Client %s left room %d", client.id, client.roomID)
		}
		if len(clients) == 0 {
			delete(h.rooms, client.roomID)
		}
	}
}

// broadcastToRoom 广播消息到房间
func (h *WebSocketHub) broadcastToRoom(roomMsg *RoomMessage) {
	h.mu.RLock()
	clients := h.rooms[roomMsg.RoomID]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.send <- roomMsg.Message:
		default:
			// 客户端发送缓冲区满，关闭连接
			h.unregister <- client
		}
	}
}

// GetRoomCount 获取房间在线人数
func (h *WebSocketHub) GetRoomCount(roomID int64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, ok := h.rooms[roomID]; ok {
		return len(clients)
	}
	return 0
}

// WebSocketService WebSocket 服务
type WebSocketService struct {
	hub *WebSocketHub
	uc  *biz.MessageUsecase
	muc *biz.RoomMemberUsecase
	log *log.Helper

	// WebSocket 升级器
	upgrader websocket.Upgrader
}

// NewWebSocketService 创建 WebSocket 服务
func NewWebSocketService(hub *WebSocketHub, uc *biz.MessageUsecase, muc *biz.RoomMemberUsecase, logger log.Logger) *WebSocketService {
	return &WebSocketService{
		hub: hub,
		uc:  uc,
		muc: muc,
		log: log.NewHelper(logger),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 生产环境需要更严格的检查
			},
		},
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (s *WebSocketService) HandleWebSocket(ctx http.Context) error {
	// 从路径获取 room_id: /ws/room/{room_id}
	// 使用 URL 路径解析
	path := ctx.Request().URL.Path
	// 路径格式: /ws/room/123
	var roomID int64
	_, err := fmt.Sscanf(path, "/ws/room/%d", &roomID)
	if err != nil || roomID <= 0 {
		return ctx.JSON(400, map[string]string{"error": "Invalid room ID"})
	}

	// 获取用户 ID - 支持多种方式
	// 1. 从请求头获取（API 网关转发）
	userIDStr := ctx.Request().Header.Get("X-User-ID")
	// 2. 从 URL 参数获取（用于测试）
	if userIDStr == "" {
		userIDStr = ctx.Query().Get("user_id")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		// 开发环境使用默认值
		userID = int64(1)
		s.log.Warnf("No valid user ID found, using default: %d", userID)
	}

	// 获取用户名（从 URL 参数）
	username := ctx.Query().Get("username")
	if username == "" {
		username = fmt.Sprintf("User%d", userID) // 默认用户名
	}

	// 验证用户是房间成员
	member, err := s.muc.Get(ctx, roomID, userID)
	if err != nil {
		return ctx.JSON(403, map[string]string{"error": "Not a member of this room"})
	}
	if member.Status != 1 { // MemberStatusNormal
		return ctx.JSON(403, map[string]string{"error": "Member status not normal"})
	}

	// 升级 HTTP 连接为 WebSocket
	conn, err := s.upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		s.log.Errorf("WebSocket upgrade failed: %v", err)
		return err
	}

	// 创建客户端（使用独立的 context，不依赖 HTTP 请求的 context）
	clientCtx, cancel := context.WithCancel(context.Background())
	client := &WebSocketClient{
		id:       fmt.Sprintf("%d-%d", userID, time.Now().UnixNano()),
		roomID:   roomID,
		userID:   userID,
		username: username,
		conn:     conn,
		send:     make(chan *WebSocketMessage, 256),
		hub:      s.hub,
		ctx:      clientCtx,
		cancel:   cancel,
	}

	// 注册客户端
	s.hub.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump(s)

	return nil
}

// readPump 读取客户端消息
func (c *WebSocketClient) readPump(ws *WebSocketService) {
	defer func() {
		ws.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ws.log.Errorf("WebSocket error: %v", err)
			}
			break
		}

		// 解析消息
		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			ws.log.Errorf("Failed to parse message: %v", err)
			continue
		}

		// 处理不同类型的消息
		switch msg.Type {
		case "message":
			ws.handleChatMessage(c, &msg)
		case "ping":
			// 响应心跳
			c.send <- &WebSocketMessage{Type: "pong"}
		default:
			ws.log.Warnf("Unknown message type: %s", msg.Type)
		}
	}
}

// handleChatMessage 处理聊天消息
func (ws *WebSocketService) handleChatMessage(client *WebSocketClient, msg *WebSocketMessage) {
	// 解析消息数据
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		client.send <- &WebSocketMessage{
			Type:  "error",
			Error: "Invalid message format",
		}
		return
	}

	content, ok := data["content"].(string)
	if !ok || content == "" {
		client.send <- &WebSocketMessage{
			Type:  "error",
			Error: "Message content is required",
		}
		return
	}

	// 发送消息
	message, err := ws.uc.Send(client.ctx, client.roomID, client.userID, content, 1)
	if err != nil {
		client.send <- &WebSocketMessage{
			Type:  "error",
			Error: err.Error(),
		}
		return
	}

	// 广播消息到房间
	ws.hub.broadcast <- &RoomMessage{
		RoomID: client.roomID,
		Message: &WebSocketMessage{
			Type: "message",
			Data: map[string]interface{}{
				"id":         message.ID,
				"room_id":    message.RoomID,
				"user_id":    message.UserID,
				"username":   client.username,
				"content":    message.Content,
				"type":       message.Type,
				"created_at": message.CreatedAt.Unix(),
			},
		},
	}
}

// writePump 写入消息到客户端
func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				c.hub.log.Errorf("Failed to marshal message: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.hub.log.Errorf("Failed to write message: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// GetHub 获取 Hub 实例（用于其他服务广播消息）
func (s *WebSocketService) GetHub() *WebSocketHub {
	return s.hub
}
