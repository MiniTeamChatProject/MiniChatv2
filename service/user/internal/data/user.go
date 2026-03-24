//path: .registration/internal/data/user.go 

package data

import (
	"context"
	"errors"
	"registration/internal/biz"

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
	Email string 
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
// 6. 实现删除用户资料
func (r *userRepo) Delete(ctx context.Context, id int64) error {
    // GORM 的删除操作
    // DELETE FROM users WHERE id = id;
    result := r.data.db.WithContext(ctx).Delete(&User{}, id)
    
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
	    return errors.New("Attention: The user does not exist or has been deleted !")
    }
    return nil
}

//update 
func (r *userRepo) Update(ctx context.Context, u *biz.User) error {
    // 1. 构造 Data 层对象 (PO)
    // 注意：我们只填入 ID 和需要修改的字段
    po := User{
        ID:       u.ID,
        Nickname: u.Nickname,
        Email:    u.Email,
    }

    // 2. 使用 GORM 的 Updates 方法 (复数 s)
    // Model(&po) 会告诉 GORM 去找 id = po.ID 的那一行
    // Updates(po) 会更新 po 里所有非空字段
    // 如果 u.Email 是空字符串，数据库里的 Email 就不会被改，这很安全
    result := r.data.db.WithContext(ctx).Model(&po).Updates(po)

    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
	    return errors.New("Errors: the user does not exist !")
    }
    return nil
}
