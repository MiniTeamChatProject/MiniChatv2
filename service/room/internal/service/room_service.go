package service

import (
	"context"
	"time"
	"room/api/room/v1"
	"room/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// RoomService 房间服务
type RoomService struct {
	v1.UnimplementedRoomServiceServer

	uc  *biz.RoomUsecase
	muc *biz.RoomMemberUsecase
	log *log.Helper
}

// NewRoomService 创建房间服务
func NewRoomService(uc *biz.RoomUsecase, muc *biz.RoomMemberUsecase, logger log.Logger) *RoomService {
	return &RoomService{
		uc:  uc,
		muc: muc,
		log: log.NewHelper(logger),
	}
}

// CreateRoom 创建房间
func (s *RoomService) CreateRoom(ctx context.Context, req *v1.CreateRoomRequest) (*v1.CreateRoomReply, error) {
	// TODO: 从 context 中获取 user_id
	userID := int64(1) // 临时硬编码

	room := &biz.Room{
		Name:        req.Name,
		OwnerID:     userID,
		Type:        biz.RoomType(req.Type),
		MaxCount:    int(req.MaxCount),
		Avatar:      req.Avatar,
		Description: req.Description,
		Tags:        req.Tags,
		IsPublic:    req.IsPublic,
	}

	if room.MaxCount == 0 {
		room.MaxCount = 500
	}

	created, err := s.uc.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	// 房主自动加入房间
	if _, err := s.muc.Join(ctx, created.ID, userID); err != nil {
		s.log.Errorf("Failed to add owner as member: %v", err)
	}

	return &v1.CreateRoomReply{
		Room: s.toProtoRoom(created),
	}, nil
}

// GetRoom 获取房间信息
func (s *RoomService) GetRoom(ctx context.Context, req *v1.GetRoomRequest) (*v1.GetRoomReply, error) {
	room, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.GetRoomReply{
		Room: s.toProtoRoom(room),
	}, nil
}

// UpdateRoom 更新房间信息
func (s *RoomService) UpdateRoom(ctx context.Context, req *v1.UpdateRoomRequest) (*v1.UpdateRoomReply, error) {
	room := &biz.Room{
		ID:          req.Id,
		Name:        req.Name,
		Avatar:      req.Avatar,
		Description: req.Description,
		Tags:        req.Tags,
		IsPublic:    req.IsPublic,
	}

	updated, err := s.uc.Update(ctx, room)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateRoomReply{
		Room: s.toProtoRoom(updated),
	}, nil
}

// DeleteRoom 删除房间
func (s *RoomService) DeleteRoom(ctx context.Context, req *v1.DeleteRoomRequest) (*v1.DeleteRoomReply, error) {
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteRoomReply{
		Success: true,
	}, nil
}

// JoinRoom 加入房间
func (s *RoomService) JoinRoom(ctx context.Context, req *v1.JoinRoomRequest) (*v1.JoinRoomReply, error) {
	// TODO: 从 context 中获取 user_id
	userID := int64(1) // 临时硬编码

	member, err := s.muc.Join(ctx, req.RoomId, userID)
	if err != nil {
		return nil, err
	}

	room, _ := s.uc.Get(ctx, req.RoomId)

	return &v1.JoinRoomReply{
		Member: s.toProtoMember(member),
		Room:   s.toProtoRoom(room),
	}, nil
}

// LeaveRoom 退出房间
func (s *RoomService) LeaveRoom(ctx context.Context, req *v1.LeaveRoomRequest) (*v1.LeaveRoomReply, error) {
	// TODO: 从 context 中获取 user_id
	userID := int64(1) // 临时硬编码

	err := s.muc.Leave(ctx, req.RoomId, userID)
	if err != nil {
		return nil, err
	}

	return &v1.LeaveRoomReply{
		Success: true,
	}, nil
}

// ListMembers 获取成员列表
func (s *RoomService) ListMembers(ctx context.Context, req *v1.ListMembersRequest) (*v1.ListMembersReply, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	members, total, err := s.muc.ListMembers(ctx, req.RoomId, page, pageSize)
	if err != nil {
		return nil, err
	}

	protoMembers := make([]*v1.RoomMember, len(members))
	for i, m := range members {
		protoMembers[i] = s.toProtoMember(m)
	}

	return &v1.ListMembersReply{
		Members: protoMembers,
		Total:   int32(total),
	}, nil
}

// KickMember 踢出成员
func (s *RoomService) KickMember(ctx context.Context, req *v1.KickMemberRequest) (*v1.KickMemberReply, error) {
	// TODO: 从 context 中获取 operator_id
	operatorID := int64(1) // 临时硬编码

	err := s.muc.Kick(ctx, req.RoomId, operatorID, req.UserId)
	if err != nil {
		return nil, err
	}

	return &v1.KickMemberReply{
		Success: true,
	}, nil
}

// UpdateMemberRole 更新成员角色
func (s *RoomService) UpdateMemberRole(ctx context.Context, req *v1.UpdateMemberRoleRequest) (*v1.UpdateMemberRoleReply, error) {
	// TODO: 从 context 中获取 operator_id
	operatorID := int64(1) // 临时硬编码

	err := s.muc.UpdateRole(ctx, req.RoomId, operatorID, req.UserId, biz.MemberRole(req.Role))
	if err != nil {
		return nil, err
	}

	return &v1.UpdateMemberRoleReply{
		Success: true,
	}, nil
}

// MuteMember 禁言/解禁成员
func (s *RoomService) MuteMember(ctx context.Context, req *v1.MuteMemberRequest) (*v1.MuteMemberReply, error) {
	// TODO: 从 context 中获取 operator_id
	operatorID := int64(1) // 临时硬编码

	var muteUntil *time.Time
	if req.MuteUntil > 0 {
		t := time.Unix(req.MuteUntil, 0)
		muteUntil = &t
	}

	err := s.muc.Mute(ctx, req.RoomId, operatorID, req.UserId, muteUntil)
	if err != nil {
		return nil, err
	}

	return &v1.MuteMemberReply{
		Success: true,
	}, nil
}

// ListUserRooms 获取用户房间列表
func (s *RoomService) ListUserRooms(ctx context.Context, req *v1.ListUserRoomsRequest) (*v1.ListUserRoomsReply, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	rooms, total, err := s.uc.ListByUserID(ctx, req.UserId, page, pageSize)
	if err != nil {
		return nil, err
	}

	protoRooms := make([]*v1.Room, len(rooms))
	for i, r := range rooms {
		protoRooms[i] = s.toProtoRoom(r)
	}

	return &v1.ListUserRoomsReply{
		Rooms: protoRooms,
		Total: int32(total),
	}, nil
}

// toProtoRoom 转换为 Proto 模型
func (s *RoomService) toProtoRoom(r *biz.Room) *v1.Room {
	if r == nil {
		return nil
	}
	return &v1.Room{
		Id:           r.ID,
		Name:         r.Name,
		OwnerId:      r.OwnerID,
		Type:         v1.RoomType(r.Type),
		MaxCount:     int32(r.MaxCount),
		CurrentCount: int32(r.CurrentCount),
		Avatar:       r.Avatar,
		Description:  r.Description,
		Tags:         r.Tags,
		IsPublic:     r.IsPublic,
		Status:       v1.RoomStatus(r.Status),
		CreatedAt:    r.CreatedAt.Unix(),
		UpdatedAt:    r.UpdatedAt.Unix(),
	}
}

// toProtoMember 转换为 Proto 模型
func (s *RoomService) toProtoMember(m *biz.RoomMember) *v1.RoomMember {
	if m == nil {
		return nil
	}
	var muteUntil int64
	if m.MuteUntil != nil {
		muteUntil = m.MuteUntil.Unix()
	}
	return &v1.RoomMember{
		Id:        m.ID,
		RoomId:    m.RoomID,
		UserId:    m.UserID,
		Role:      v1.MemberRole(m.Role),
		Status:    v1.MemberStatus(m.Status),
		MuteUntil: muteUntil,
		JoinedAt:  m.JoinedAt.Unix(),
		UpdatedAt: m.UpdatedAt.Unix(),
	}
}
