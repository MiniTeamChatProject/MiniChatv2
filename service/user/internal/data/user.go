//path: .user/internal/data/user.go 

package data

import (
	"context"
	"errors"
	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// 1. 定义数据库模型 (PO - Persistent Object)
// 这才是真正跟数据库表 users 对应及绑定的结构体
type User struct {
	ID       int64  `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Nickname string
}

// 2. 定义仓库结构体
type userRepo struct {
	data *Data
	log  *log.Helper
}

// 初始化函数
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// ======================================
// 实现 Biz 层定义的接口
// ======================================

// 3. 实现 Save (注册用)
func (r *userRepo) Save(ctx context.Context, u *biz.User) (*biz.User, error) {
	// A. 模型转换：Biz对象 -> Data对象
	po := User{
		Username: u.Username,
		Password: u.Password,
		Nickname: u.Nickname,
	}

	// B. 存入数据库
	result := r.data.db.WithContext(ctx).Create(&po)
	if result.Error != nil {
		return nil, result.Error
	}

	// C. 把生成的 ID 赋回去
	u.ID = po.ID
	return u, nil
}

// 4. 实现 FindByUsername (登录用)
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	var po User
	// 去数据库查
	result := r.data.db.WithContext(ctx).Where("username = ?", username).First(&po)
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // 没找到
	}
	if result.Error != nil {
		return nil, result.Error
	}

	// 模型转换：Data对象 -> Biz对象
	return &biz.User{
		ID:       po.ID,
		Username: po.Username,
		Password: po.Password,
		Nickname: po.Nickname,
	}, nil
}

// 5. 实现 FindByID (获取信息用)
func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	var po User
	result := r.data.db.WithContext(ctx).First(&po, id)
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &biz.User{
		ID:       po.ID,
		Username: po.Username,
		Password: po.Password,
		Nickname: po.Nickname,
	}, nil
}
