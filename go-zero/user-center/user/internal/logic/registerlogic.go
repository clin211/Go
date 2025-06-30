package logic

import (
	"context"

	"user/internal/svc"
	"user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterRes, err error) {
	// 参数验证
	if req.Username == "" {
		return &types.RegisterRes{
			Code:    400,
			Message: "用户名不能为空",
		}, nil
	}

	if req.Password == "" {
		return &types.RegisterRes{
			Code:    400,
			Message: "密码不能为空",
		}, nil
	}

	if req.Email == "" {
		return &types.RegisterRes{
			Code:    400,
			Message: "邮箱不能为空",
		}, nil
	}

	// 这里可以添加更多验证逻辑
	// 比如检查用户名是否已存在、密码强度等

	// 模拟注册成功
	logx.Infof("用户注册: username=%s, email=%s", req.Username, req.Email)

	return &types.RegisterRes{
		Code:    200,
		Message: "注册成功",
	}, nil
}
