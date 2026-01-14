// path: registration/internal/biz/user.go 
//maintainer: Jonas Peng
//logs1: 2025/Dev/17: rebuild whole codes

package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ==========================================
// 1. 定义领域对象 (Domain Objects)
// ==========================================
// 这里只定义数据长什么样，不要加 gorm:"..." 标签
// Data 层会负责把它转换成数据库能存的格式
type User struct {
	ID       int64
	Username string
	Password string
	Nickname string
	Email string 
}

// ==========================================
// 2. 定义接口 (Repository Interface)
// ==========================================
// Biz 层发号施令：Data 层，我需要这几个功能，你去实现！
type UserRepo interface {
	Save(context.Context, *User) (*User, error)//signup user info 
	FindByUsername(context.Context, string) (*User, error)//validate user info
	FindByID(context.Context, int64) (*User, error)// fetch user info
	Delete(context.Context, int64) error // delete user info
	Update(context.Context, *User) error //update user info
}

// ==========================================
// 3. 业务逻辑核心 (UseCase)
// ==========================================
type UserUseCase struct {
	repo   UserRepo
	log    *log.Helper
	jwtKey []byte // 用于签发 Token 的密钥
}

// 初始化函数 (Wire 会调用)
func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo:   repo,
		log:    log.NewHelper(logger),
		jwtKey: []byte("MySecretKey_0721"), // 暂时写死，以后可以从配置文件读
	}
}

// --- 业务功能 1: 注册 ---
func (uc *UserUseCase) Register(ctx context.Context, username, password, nickname string) error {
	uc.log.WithContext(ctx).Infof("开始注册用户: %s", username)

	// 1. 检查用户是否已存在 (这一步也可以在 Data 层做，但这里做逻辑更清晰)
	existing, _ := uc.repo.FindByUsername(ctx, username)
	if existing != nil {
		return errors.New("用户名已存在")
	}

	// 2. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. 构建用户对象
	u := &User{
		Username: username,
		Password: string(hashedPassword),
		Nickname: nickname,
	}

	// 4. 保存到数据库 (调用 Data 层)
	_, err = uc.repo.Save(ctx, u)
	return err
}

// --- 业务功能 2: 登录 ---
func (uc *UserUseCase) Login(ctx context.Context, username, password string) (string, error) {
	uc.log.WithContext(ctx).Infof("用户尝试登录: %s", username)

	// 1. 查找用户
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New("用户不存在或数据库错误") // 也可以返回 nil
	}
	if user == nil {
		return "", errors.New("用户不存在")
	}

	// 2. 比对密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("密码错误")
	}

	// 3. 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	})

	tokenString, err := token.SignedString(uc.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// --- 业务功能 3: 获取信息 ---
func (uc *UserUseCase) GetProfile(ctx context.Context, id int64) (*User, error) {
	return uc.repo.FindByID(ctx, id)
}

//--- 业务功能 4：删除用户 ----
func (uc *UserUseCase) DeleteUser(ctx context.Context, id int64) error {
    uc.log.WithContext(ctx).Infof("正在删除用户 ID: %d", id)
    // 这里以后可以加逻辑：比如检查是不是删自己，或者是不是管理员
    return uc.repo.Delete(ctx, id)
}
//-----业务功能 5：变更用户信息-----
func (uc *UserUseCase) UpdateProfile(ctx context.Context, id int64, nickname, email string) error {
    uc.log.WithContext(ctx).Infof("用户 %d 正在更新资料", id)

    // 构造一个只包含 ID 和要修改数据的对象
    u := &User{
        ID:       id,
        Nickname: nickname,
        Email:    email,
    }

    // 调用 Data 层
    return uc.repo.Update(ctx, u)
}
