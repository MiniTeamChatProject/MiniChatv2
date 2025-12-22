package service

import (
	"context"

	// 1. 引入 Proto 定义，并起别名为 pb
	pb "user/api/helloworld/v1"
	"user/internal/biz"
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
	return &pb.RegisterReply{Msg: "注册成功"}, nil
}

// 5. 实现 Login 接口
func (s *RegistrationService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginReply, error) {
	token, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Msg:   "登录成功",
		Token: token,
	}, nil
}

// 6. 实现 GetProfile 接口
func (s *RegistrationService) GetProfile(ctx context.Context, req *pb.GetProfileReq) (*pb.GetProfileReply, error) {
	user, err := s.uc.GetProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProfileReply{
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}
