package data

import (
	"context"
	"fmt"
	"room/internal/biz"
	"room/internal/data/model"
	"strings"

	"gorm.io/gorm"
)

// roomRepo 房间仓储实现
type roomRepo struct {
	data *Data
}

// NewRoomRepo 创建房间仓储
func NewRoomRepo(data *Data) biz.RoomRepository {
	return &roomRepo{data: data}
}

// Create 创建房间
func (r *roomRepo) Create(ctx context.Context, room *biz.Room) (*biz.Room, error) {
	roomModel := &model.Room{
		Name:        room.Name,
		OwnerID:     room.OwnerID,
		Type:        model.RoomType(room.Type),
		MaxCount:    room.MaxCount,
		Avatar:      room.Avatar,
		Description: room.Description,
		Tags:        strings.Join(room.Tags, ","),
		IsPublic:    room.IsPublic,
		Status:      model.RoomStatus(room.Status),
	}

	if err := r.data.db.WithContext(ctx).Create(roomModel).Error; err != nil {
		return nil, err
	}

	return r.toBizRoom(roomModel), nil
}

// Get 获取房间
func (r *roomRepo) Get(ctx context.Context, id int64) (*biz.Room, error) {
	var roomModel model.Room
	err := r.data.db.WithContext(ctx).First(&roomModel, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("room not found")
		}
		return nil, err
	}

	return r.toBizRoom(&roomModel), nil
}

// Update 更新房间
func (r *roomRepo) Update(ctx context.Context, room *biz.Room) (*biz.Room, error) {
	updates := map[string]interface{}{
		"name":        room.Name,
		"avatar":      room.Avatar,
		"description": room.Description,
		"is_public":   room.IsPublic,
	}
	if room.Tags != nil {
		updates["tags"] = strings.Join(room.Tags, ",")
	}

	err := r.data.db.WithContext(ctx).Model(&model.Room{}).
		Where("id = ?", room.ID).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, room.ID)
}

// Delete 删除房间
func (r *roomRepo) Delete(ctx context.Context, id int64) error {
	// 先删除成员关联
	err := r.data.db.WithContext(ctx).Where("room_id = ?", id).Delete(&model.RoomMember{}).Error
	if err != nil {
		return err
	}

	// 删除房间
	return r.data.db.WithContext(ctx).Delete(&model.Room{}, id).Error
}

// ListByUserID 获取用户加入的房间列表
func (r *roomRepo) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*biz.Room, int, error) {
	var rooms []*model.Room
	var total int64

	// 获取房间ID列表
	var roomIDs []int64
	err := r.data.db.WithContext(ctx).Model(&model.RoomMember{}).
		Where("user_id = ? AND status = ?", userID, model.MemberStatusNormal).
		Pluck("room_id", &roomIDs).Error
	if err != nil {
		return nil, 0, err
	}

	if len(roomIDs) == 0 {
		return []*biz.Room{}, 0, nil
	}

	// 获取总数
	err = r.data.db.WithContext(ctx).Model(&model.Room{}).
		Where("id IN ?", roomIDs).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取房间列表
	err = r.data.db.WithContext(ctx).
		Where("id IN ?", roomIDs).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rooms).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*biz.Room, len(rooms))
	for i, room := range rooms {
		result[i] = r.toBizRoom(room)
	}

	return result, int(total), nil
}

// GetMemberCount 获取成员数量
func (r *roomRepo) GetMemberCount(ctx context.Context, roomID int64) (int, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&model.RoomMember{}).
		Where("room_id = ? AND status = ?", roomID, model.MemberStatusNormal).
		Count(&count).Error
	return int(count), err
}

// toBizRoom 转换为领域模型
func (r *roomRepo) toBizRoom(m *model.Room) *biz.Room {
	tags := []string{}
	if m.Tags != "" {
		tags = strings.Split(m.Tags, ",")
	}

	return &biz.Room{
		ID:          m.ID,
		Name:        m.Name,
		OwnerID:     m.OwnerID,
		Type:        biz.RoomType(m.Type),
		MaxCount:    m.MaxCount,
		Avatar:      m.Avatar,
		Description: m.Description,
		Tags:        tags,
		IsPublic:    m.IsPublic,
		Status:      biz.RoomStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
