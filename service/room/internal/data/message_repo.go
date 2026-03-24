package data

import (
	"context"
	"room/internal/biz"
	"room/internal/data/model"
)

// messageRepo 消息仓储实现
type messageRepo struct {
	data *Data
}

// NewMessageRepo 创建消息仓储
func NewMessageRepo(data *Data) biz.MessageRepository {
	return &messageRepo{data: data}
}

// Create 创建消息
func (r *messageRepo) Create(ctx context.Context, message *biz.Message) (*biz.Message, error) {
	messageModel := &model.Message{
		RoomID:  message.RoomID,
		UserID:  message.UserID,
		Content: message.Content,
		Type:    model.MessageType(message.Type),
	}

	if err := r.data.db.WithContext(ctx).Create(messageModel).Error; err != nil {
		return nil, err
	}

	return r.toBizMessage(messageModel), nil
}

// ListByRoomID 获取房间消息列表
func (r *messageRepo) ListByRoomID(ctx context.Context, roomID int64, limit, offset int) ([]*biz.Message, int, error) {
	var messages []*model.Message
	var total int64

	// 获取总数
	err := r.data.db.WithContext(ctx).Model(&model.Message{}).
		Where("room_id = ?", roomID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取消息列表（按时间倒序）
	err = r.data.db.WithContext(ctx).
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为领域模型并反转顺序（最新的在前）
	result := make([]*biz.Message, len(messages))
	for i, msg := range messages {
		result[len(messages)-1-i] = r.toBizMessage(msg)
	}

	return result, int(total), nil
}

// toBizMessage 转换为领域模型
func (r *messageRepo) toBizMessage(m *model.Message) *biz.Message {
	return &biz.Message{
		ID:        m.ID,
		RoomID:    m.RoomID,
		UserID:    m.UserID,
		Content:   m.Content,
		Type:      biz.MessageType(m.Type),
		CreatedAt: m.CreatedAt,
	}
}
