package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrMessageTooLong    = errors.New("message content is too long")
	ErrNotRoomMember     = errors.New("user is not a member of this room")
	ErrUserMuted         = errors.New("user is muted")
	ErrMessageNotFound   = errors.New("message not found")
)

// MessageUsecase 消息用例
type MessageUsecase struct {
	messageRepo MessageRepository
	memberRepo  RoomMemberRepository
	log         *log.Helper
}

// NewMessageUsecase 创建消息用例
func NewMessageUsecase(messageRepo MessageRepository, memberRepo RoomMemberRepository, logger log.Logger) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
		memberRepo:  memberRepo,
		log:         log.NewHelper(logger),
	}
}

// Send 发送消息
func (uc *MessageUsecase) Send(ctx context.Context, roomID, userID int64, content string, msgType MessageType) (*Message, error) {
	// 验证消息内容
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}
	if len(content) > 5000 {
		return nil, ErrMessageTooLong
	}

	// 验证用户是房间成员
	member, err := uc.memberRepo.Get(ctx, roomID, userID)
	if err != nil {
		return nil, ErrNotRoomMember
	}

	// 验证用户状态
	if member.Status != MemberStatusNormal {
		return nil, ErrNotRoomMember
	}

	// 检查用户是否被禁言
	if member.MuteUntil != nil && member.MuteUntil.After(time.Now()) {
		return nil, ErrUserMuted
	}

	// 创建消息
	message := &Message{
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		Type:      msgType,
		CreatedAt: time.Now(),
	}

	// 默认为文本消息
	if message.Type == 0 {
		message.Type = MessageTypeText
	}

	return uc.messageRepo.Create(ctx, message)
}

// ListMessages 获取消息历史
func (uc *MessageUsecase) ListMessages(ctx context.Context, roomID int64, page, pageSize int) ([]*Message, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	return uc.messageRepo.ListByRoomID(ctx, roomID, pageSize, offset)
}
