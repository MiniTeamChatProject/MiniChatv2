package model

import "time"

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

// Room 房间模型
type Room struct {
	ID          int64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string      `gorm:"column:name;type:varchar(100);not null" json:"name"`
	OwnerID     int64       `gorm:"column:owner_id;not null;index" json:"owner_id"`
	Type        RoomType    `gorm:"column:type;type:smallint;not null;default:1" json:"type"`
	MaxCount    int         `gorm:"column:max_count;type:int;default:500" json:"max_count"`
	Avatar      string      `gorm:"column:avatar;type:varchar(255)" json:"avatar,omitempty"`
	Description string      `gorm:"column:description;type:text" json:"description,omitempty"`
	Tags        string      `gorm:"column:tags;type:varchar(500)" json:"tags,omitempty"`
	IsPublic    bool        `gorm:"column:is_public;type:boolean;default:true" json:"is_public"`
	Status      RoomStatus  `gorm:"column:status;type:smallint;default:1" json:"status"`
	CreatedAt   time.Time   `gorm:"column:created_at;type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;type:timestamptz;default:now()" json:"updated_at"`

	// 关联
	Members     []RoomMember `gorm:"foreignKey:RoomID" json:"members,omitempty"`
}

// TableName 指定表名
func (Room) TableName() string {
	return "rooms"
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

// RoomMember 房间成员模型
type RoomMember struct {
	ID        int64        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoomID    int64        `gorm:"column:room_id;not null;uniqueIndex:idx_room_user" json:"room_id"`
	UserID    int64        `gorm:"column:user_id;not null;uniqueIndex:idx_room_user;index" json:"user_id"`
	Role      MemberRole   `gorm:"column:role;type:smallint;default:1" json:"role"`
	Status    MemberStatus `gorm:"column:status;type:smallint;default:1" json:"status"`
	MuteUntil *time.Time   `gorm:"column:mute_until;type:timestamptz" json:"mute_until,omitempty"`
	JoinedAt  time.Time    `gorm:"column:joined_at;type:timestamptz;default:now()" json:"joined_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at;type:timestamptz;default:now()" json:"updated_at"`

	// 关联
	Room *Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

// TableName 指定表名
func (RoomMember) TableName() string {
	return "room_members"
}
