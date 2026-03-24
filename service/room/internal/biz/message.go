package biz

import (
	"context"
	"time"
)

// MessageType 消息类型
type MessageType int16

const (
	MessageTypeText  MessageType = 1 // 文本消息
	MessageTypeImage MessageType = 2 // 图片消息（预留）
	MessageTypeVoice MessageType = 3 // 语音消息（预留）
	MessageTypeVideo MessageType = 4 // 视频消息（预留）
)

// Message 消息领域模型
type Message struct {
	ID        int64
	RoomID    int64
	UserID    int64
	Content   string
	Type      MessageType
	CreatedAt time.Time
}

// MessageRepository 消息仓储接口
type MessageRepository interface {
	Create(ctx context.Context, message *Message) (*Message, error)
	ListByRoomID(ctx context.Context, roomID int64, limit, offset int) ([]*Message, int, error)
}
