package service

import (
	"context"
	"errors"

	// 1. 引入 Proto 定义，并起别名为 pb
	pb "registration/api/helloworld/v1"
	"registration/internal/biz"
)

// 2. 定义 Service 结构体
type RegistrationService struct {
	pb.UnimplementedRegistrationServer
	
	// 这里一定要用大写 C 的 UseCase，和 Biz 层保持一致
	uc *biz.UserUseCase
}

// 3. 初始化函数
// 注意：这里返回的是 *RegistrationService，不能是 GreeterService
func NewRegistrationService(uc *biz.UserUseCase) *RegistrationService {
	return &RegistrationService{uc: uc}
}

// 4. 实现 Register 接口
func (s *RegistrationService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterReply, error) {
	// 调用 Biz 层
	err := s.uc.Register(ctx, req.Username, req.Password, req.Nickname)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{Msg: "Registered successfully !"}, nil
}

// 5. 实现 Login 接口
func (s *RegistrationService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {
	token, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Msg:   "Loged in !",
		Token: token,
	}, nil
}

// 6. 实现 GetProfile 接口
func (s *RegistrationService) GetProfile(ctx context.Context, req *pb.GetProfileReq) (*pb.GetProfileReply, error) {
	// 【关键修改】从 Context 中提取 Token 解析出来的 user_id
	// 而不是使用 req.Id (前端传的参数)
	userIdVal := ctx.Value("user_id")
	if userIdVal == nil {
		return nil, errors.New("Unanthenticated：unable to access user ID !")
	}

	// 类型断言：确保它是 int64
	var userId int64
	switch v := userIdVal.(type) {
	case int64:
		userId = v
	case float64:
		userId = int64(v)
	default:
		return nil, errors.New("Invalid type of userID !")
	}

	// 使用从 Token 拿到的 userId 去查询
	user, err := s.uc.GetProfile(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileReply{
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

// 7. 实现 DeleteUser 接口
func (s *RegistrationService) DeleteUser(ctx context.Context, req *pb.DeleteUserReq) (*pb.DeleteUserReply, error) {
    // req.Id 会自动从 URL 的 /user/{id} 中提取出来
    err := s.uc.DeleteUser(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    
    return &pb.DeleteUserReply{
        Msg: "The user has been successfully !",
    }, nil
}

//Delete User Info api
// 实现 UpdateProfile
func (s *RegistrationService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileReq) (*pb.UpdateProfileReply, error) {
    // 【关键安全步骤】
    // 从 Context 中提取 User ID，而不是让前端传 ID
    // 假设你的中间件把解析出来的 ID 存为了 "user_id"
    // 如果使用 krtos/v2/middleware/auth/jwt，获取方式可能不同
    
    // 这里演示通用的从 Value 获取 (你需要确保 Server 层中间件塞进去了)
    userIdVal := ctx.Value("user_id") 
    if userIdVal == nil {
        return nil, errors.New("Unauthenticated：unable to fetch user ID !")
    }
    
    // 类型断言：确保它是 int64 或 float64 (JWT解析出来经常是 float64)
    var userId int64
    switch v := userIdVal.(type) {
    case int64:
        userId = v
    case float64:
        userId = int64(v)
    default:
        return nil, errors.New("Invalid type of user ID !")
    }

    // 调用 Biz 层
    err := s.uc.UpdateProfile(ctx, userId, req.Nickname, req.Email)
    if err != nil {
        return nil, err
    }

    return &pb.UpdateProfileReply{Msg: "The information has been updated successfully !"}, nil
}

// 8. VerifyUser - 服务间通信接口，验证用户是否存在
func (s *RegistrationService) VerifyUser(ctx context.Context, req *pb.VerifyUserReq) (*pb.VerifyUserReply, error) {
	user, err := s.uc.GetProfile(ctx, req.UserId)
	if err != nil {
		return &pb.VerifyUserReply{
			Valid:  false,
			UserId: req.UserId,
		}, nil
	}

	return &pb.VerifyUserReply{
		Valid:    true,
		UserId:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

// 9. GetUser - 服务间通信接口，获取用户信息
func (s *RegistrationService) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserReply, error) {
	user, err := s.uc.GetProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserReply{
		Id:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
	}, nil
}

