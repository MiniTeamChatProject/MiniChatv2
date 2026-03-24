package model

import "time"

// MessageType 消息类型
type MessageType int16

const (
	MessageTypeText  MessageType = 1 // 文本消息
	MessageTypeImage MessageType = 2 // 图片消息（预留）
	MessageTypeVoice MessageType = 3 // 语音消息（预留）
	MessageTypeVideo MessageType = 4 // 视频消息（预留）
)

// Message 消息模型
type Message struct {
	ID        int64        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoomID    int64        `gorm:"column:room_id;not null;index" json:"room_id"`
	UserID    int64        `gorm:"column:user_id;not null;index" json:"user_id"`
	Content   string       `gorm:"column:content;type:text;not null" json:"content"`
	Type      MessageType  `gorm:"column:type;type:smallint;not null;default:1" json:"type"`
	CreatedAt time.Time    `gorm:"column:created_at;type:timestamptz;default:now()" json:"created_at"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
