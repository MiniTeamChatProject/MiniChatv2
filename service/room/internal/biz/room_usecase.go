package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrRoomNotFound      = errors.New("room not found")
	ErrRoomAlreadyJoined = errors.New("user already in room")
	ErrRoomFull          = errors.New("room is full")
	ErrMemberNotFound    = errors.New("member not found")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrInvalidRole       = errors.New("invalid role")
)

// RoomUsecase 房间用例
type RoomUsecase struct {
	repo RoomRepository
	log  *log.Helper
}

// NewRoomUsecase 创建房间用例
func NewRoomUsecase(repo RoomRepository, logger log.Logger) *RoomUsecase {
	return &RoomUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Create 创建房间
func (uc *RoomUsecase) Create(ctx context.Context, room *Room) (*Room, error) {
	if room.Name == "" {
		return nil, errors.New("room name is required")
	}
	if room.OwnerID == 0 {
		return nil, errors.New("owner ID is required")
	}
	if room.Type == 0 {
		room.Type = RoomTypeGroup
	}
	if room.MaxCount == 0 {
		room.MaxCount = 500
	}
	if room.Status == 0 {
		room.Status = RoomStatusNormal
	}

	return uc.repo.Create(ctx, room)
}

// Get 获取房间
func (uc *RoomUsecase) Get(ctx context.Context, id int64) (*Room, error) {
	room, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, ErrRoomNotFound
	}

	// 获取当前成员数
	count, err := uc.repo.GetMemberCount(ctx, id)
	if err == nil {
		room.CurrentCount = count
	}

	return room, nil
}

// Update 更新房间
func (uc *RoomUsecase) Update(ctx context.Context, room *Room) (*Room, error) {
	existing, err := uc.repo.Get(ctx, room.ID)
	if err != nil {
		return nil, ErrRoomNotFound
	}

	// 只允许更新部分字段
	existing.Name = room.Name
	existing.Avatar = room.Avatar
	existing.Description = room.Description
	if room.Tags != nil {
		existing.Tags = room.Tags
	}
	existing.IsPublic = room.IsPublic

	return uc.repo.Update(ctx, existing)
}

// Delete 删除房间
func (uc *RoomUsecase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.Get(ctx, id)
	if err != nil {
		return ErrRoomNotFound
	}
	return uc.repo.Delete(ctx, id)
}

// ListByUserID 获取用户加入的房间列表
func (uc *RoomUsecase) ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Room, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return uc.repo.ListByUserID(ctx, userID, pageSize, offset)
}

// CanJoin 检查是否可以加入房间
func (uc *RoomUsecase) CanJoin(ctx context.Context, roomID int64) error {
	room, err := uc.repo.Get(ctx, roomID)
	if err != nil {
		return ErrRoomNotFound
	}

	if room.Status != RoomStatusNormal {
		return errors.New("room is not available")
	}

	count, err := uc.repo.GetMemberCount(ctx, roomID)
	if err == nil && count >= room.MaxCount {
		return ErrRoomFull
	}

	return nil
}

// RoomMemberUsecase 房间成员用例
type RoomMemberUsecase struct {
	memberRepo RoomMemberRepository
	roomRepo   RoomRepository
	log        *log.Helper
}

// NewRoomMemberUsecase 创建房间成员用例
func NewRoomMemberUsecase(memberRepo RoomMemberRepository, roomRepo RoomRepository, logger log.Logger) *RoomMemberUsecase {
	return &RoomMemberUsecase{
		memberRepo: memberRepo,
		roomRepo:   roomRepo,
		log:        log.NewHelper(logger),
	}
}

// Join 加入房间
func (uc *RoomMemberUsecase) Join(ctx context.Context, roomID, userID int64) (*RoomMember, error) {
	// 检查房间是否可以加入
	if err := uc.checkRoomAvailability(ctx, roomID); err != nil {
		return nil, err
	}

	// 检查是否已加入
	_, err := uc.memberRepo.Get(ctx, roomID, userID)
	if err == nil {
		return nil, ErrRoomAlreadyJoined
	}

	member := &RoomMember{
		RoomID:    roomID,
		UserID:    userID,
		Role:      RoleMember,
		Status:    MemberStatusNormal,
		JoinedAt:  time.Now(),
		UpdatedAt: time.Now(),
	}

	return uc.memberRepo.Create(ctx, member)
}

// Leave 退出房间
func (uc *RoomMemberUsecase) Leave(ctx context.Context, roomID, userID int64) error {
	_, err := uc.memberRepo.Get(ctx, roomID, userID)
	if err != nil {
		return ErrMemberNotFound
	}
	return uc.memberRepo.Delete(ctx, roomID, userID)
}

// Kick 踢出成员
func (uc *RoomMemberUsecase) Kick(ctx context.Context, roomID, operatorID, targetUserID int64) error {
	// 检查操作者权限
	operator, err := uc.memberRepo.Get(ctx, roomID, operatorID)
	if err != nil {
		return ErrPermissionDenied
	}

	// 检查目标成员
	target, err := uc.memberRepo.Get(ctx, roomID, targetUserID)
	if err != nil {
		return ErrMemberNotFound
	}

	// 权限检查：房主可以踢任何人，管理员不能踢房主和管理员
	if operator.Role != RoleOwner {
		if operator.Role != RoleAdmin {
			return ErrPermissionDenied
		}
		if target.Role == RoleAdmin || target.Role == RoleOwner {
			return ErrPermissionDenied
		}
	}

	return uc.memberRepo.Delete(ctx, roomID, targetUserID)
}

// UpdateRole 更新成员角色
func (uc *RoomMemberUsecase) UpdateRole(ctx context.Context, roomID, operatorID, targetUserID int64, newRole MemberRole) error {
	// 只有房主可以更新角色
	operator, err := uc.memberRepo.Get(ctx, roomID, operatorID)
	if err != nil || operator.Role != RoleOwner {
		return ErrPermissionDenied
	}

	// 不能设置为房主（房主需要转移房间）
	if newRole == RoleOwner {
		return ErrInvalidRole
	}

	return uc.memberRepo.UpdateRole(ctx, roomID, targetUserID, newRole)
}

// Mute 禁言/解禁成员
func (uc *RoomMemberUsecase) Mute(ctx context.Context, roomID, operatorID, targetUserID int64, muteUntil *time.Time) error {
	// 检查操作者权限
	operator, err := uc.memberRepo.Get(ctx, roomID, operatorID)
	if err != nil || (operator.Role != RoleOwner && operator.Role != RoleAdmin) {
		return ErrPermissionDenied
	}

	// 不能禁言房主
	target, err := uc.memberRepo.Get(ctx, roomID, targetUserID)
	if err != nil {
		return ErrMemberNotFound
	}
	if target.Role == RoleOwner {
		return ErrPermissionDenied
	}

	return uc.memberRepo.UpdateMute(ctx, roomID, targetUserID, muteUntil)
}

// ListMembers 获取成员列表
func (uc *RoomMemberUsecase) ListMembers(ctx context.Context, roomID int64, page, pageSize int) ([]*RoomMember, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	return uc.memberRepo.List(ctx, roomID, pageSize, offset)
}

// checkRoomAvailability 检查房间是否可用
func (uc *RoomMemberUsecase) checkRoomAvailability(ctx context.Context, roomID int64) error {
	room, err := uc.roomRepo.Get(ctx, roomID)
	if err != nil {
		return ErrRoomNotFound
	}

	if room.Status != RoomStatusNormal {
		return errors.New("room is not available")
	}

	count, err := uc.roomRepo.GetMemberCount(ctx, roomID)
	if err == nil && count >= room.MaxCount {
		return ErrRoomFull
	}

	return nil
}
