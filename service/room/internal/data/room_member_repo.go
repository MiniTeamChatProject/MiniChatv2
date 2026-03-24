package data

import (
	"context"
	"fmt"
	"room/internal/biz"
	"room/internal/data/model"
	"time"

	"gorm.io/gorm"
)

// roomMemberRepo 房间成员仓储实现
type roomMemberRepo struct {
	data *Data
}

// NewRoomMemberRepo 创建房间成员仓储
func NewRoomMemberRepo(data *Data) biz.RoomMemberRepository {
	return &roomMemberRepo{data: data}
}

// Create 创建成员
func (r *roomMemberRepo) Create(ctx context.Context, member *biz.RoomMember) (*biz.RoomMember, error) {
	memberModel := &model.RoomMember{
		RoomID:    member.RoomID,
		UserID:    member.UserID,
		Role:      model.MemberRole(member.Role),
		Status:    model.MemberStatus(member.Status),
		MuteUntil: member.MuteUntil,
	}

	if err := r.data.db.WithContext(ctx).Create(memberModel).Error; err != nil {
		return nil, err
	}

	return r.toBizMember(memberModel), nil
}

// Get 获取成员
func (r *roomMemberRepo) Get(ctx context.Context, roomID, userID int64) (*biz.RoomMember, error) {
	var memberModel model.RoomMember
	err := r.data.db.WithContext(ctx).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		First(&memberModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("member not found")
		}
		return nil, err
	}

	return r.toBizMember(&memberModel), nil
}

// Update 更新成员
func (r *roomMemberRepo) Update(ctx context.Context, member *biz.RoomMember) (*biz.RoomMember, error) {
	updates := map[string]interface{}{
		"role":   model.MemberRole(member.Role),
		"status": model.MemberStatus(member.Status),
	}
	if member.MuteUntil != nil {
		updates["mute_until"] = member.MuteUntil
	}

	err := r.data.db.WithContext(ctx).Model(&model.RoomMember{}).
		Where("room_id = ? AND user_id = ?", member.RoomID, member.UserID).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, member.RoomID, member.UserID)
}

// Delete 删除成员
func (r *roomMemberRepo) Delete(ctx context.Context, roomID, userID int64) error {
	return r.data.db.WithContext(ctx).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.RoomMember{}).Error
}

// List 获取成员列表
func (r *roomMemberRepo) List(ctx context.Context, roomID int64, limit, offset int) ([]*biz.RoomMember, int, error) {
	var members []*model.RoomMember
	var total int64

	query := r.data.db.WithContext(ctx).
		Where("room_id = ? AND status = ?", roomID, model.MemberStatusNormal)

	// 获取总数
	err := query.Model(&model.RoomMember{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	err = query.Order("role DESC, joined_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*biz.RoomMember, len(members))
	for i, member := range members {
		result[i] = r.toBizMember(member)
	}

	return result, int(total), nil
}

// UpdateRole 更新成员角色
func (r *roomMemberRepo) UpdateRole(ctx context.Context, roomID, userID int64, role biz.MemberRole) error {
	return r.data.db.WithContext(ctx).Model(&model.RoomMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("role", model.MemberRole(role)).Error
}

// UpdateMute 更新禁言状态
func (r *roomMemberRepo) UpdateMute(ctx context.Context, roomID, userID int64, muteUntil *time.Time) error {
	updates := map[string]interface{}{
		"mute_until": muteUntil,
	}
	// 如果设置禁言时间，更新状态为禁言中
	if muteUntil != nil && muteUntil.After(time.Now()) {
		updates["status"] = model.MemberStatusMuted
	} else {
		updates["status"] = model.MemberStatusNormal
	}

	return r.data.db.WithContext(ctx).Model(&model.RoomMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Updates(updates).Error
}

// toBizMember 转换为领域模型
func (r *roomMemberRepo) toBizMember(m *model.RoomMember) *biz.RoomMember {
	return &biz.RoomMember{
		ID:        m.ID,
		RoomID:    m.RoomID,
		UserID:    m.UserID,
		Role:      biz.MemberRole(m.Role),
		Status:    biz.MemberStatus(m.Status),
		MuteUntil: m.MuteUntil,
		JoinedAt:  m.JoinedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
