package biz

import (
	"context"
	"time"
)

// RoomType 房间类型
type RoomType int16

const (
	RoomTypeGroup  RoomType = 1 // 普通群聊
	RoomTypeVoice  RoomType = 2 // 语音房
	RoomTypeVideo  RoomType = 3 // 视频房
	RoomTypeLive   RoomType = 4 // 直播间
)

// RoomStatus 房间状态
type RoomStatus int16

const (
	RoomStatusNormal RoomStatus = 1 // 正常
	RoomStatusMuted  RoomStatus = 2 // 禁言
	RoomStatusBanned RoomStatus = 3 // 封禁
)

// Room 房间领域模型
type Room struct {
	ID          int64
	Name        string
	OwnerID     int64
	Type        RoomType
	MaxCount    int
	Avatar      string
	Description string
	Tags        []string
	IsPublic    bool
	Status      RoomStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// 非持久化字段
	CurrentCount int
}

// RoomRepository 房间仓储接口
type RoomRepository interface {
	Create(ctx context.Context, room *Room) (*Room, error)
	Get(ctx context.Context, id int64) (*Room, error)
	Update(ctx context.Context, room *Room) (*Room, error)
	Delete(ctx context.Context, id int64) error
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*Room, int, error)
	GetMemberCount(ctx context.Context, roomID int64) (int, error)
}

// MemberRole 成员角色
type MemberRole int16

const (
	RoleMember MemberRole = 1 // 普通成员
	RoleAdmin  MemberRole = 2 // 管理员
	RoleOwner  MemberRole = 3 // 房主
)

// MemberStatus 成员状态
type MemberStatus int16

const (
	MemberStatusNormal MemberStatus = 1 // 正常
	MemberStatusMuted  MemberStatus = 2 // 禁言中
	MemberStatusLeft   MemberStatus = 3 // 已离开
)

// RoomMember 房间成员领域模型
type RoomMember struct {
	ID        int64
	RoomID    int64
	UserID    int64
	Role      MemberRole
	Status    MemberStatus
	MuteUntil *time.Time
	JoinedAt  time.Time
	UpdatedAt time.Time
}

// RoomMemberRepository 房间成员仓储接口
type RoomMemberRepository interface {
	Create(ctx context.Context, member *RoomMember) (*RoomMember, error)
	Get(ctx context.Context, roomID, userID int64) (*RoomMember, error)
	Update(ctx context.Context, member *RoomMember) (*RoomMember, error)
	Delete(ctx context.Context, roomID, userID int64) error
	List(ctx context.Context, roomID int64, limit, offset int) ([]*RoomMember, int, error)
	UpdateRole(ctx context.Context, roomID, userID int64, role MemberRole) error
	UpdateMute(ctx context.Context, roomID, userID int64, muteUntil *time.Time) error
}
