package logic

import (
	"context"

	"user/internal/svc"
	"user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	// 参数验证
	if req.Username == "" || req.Password == "" {
		return &types.LoginRes{
			Code:    400,
			Message: "用户名和密码不能为空",
		}, nil
	}

	// 简单的用户验证（实际项目中应该查询数据库）
	if req.Username == "admin" && req.Password == "123456" {
		logx.Infof("用户登录成功: username=%s", req.Username)
		return &types.LoginRes{
			Code:    200,
			Message: "登录成功",
			Token:   "mock-jwt-token-" + req.Username,
		}, nil
	}

	// 登录失败
	logx.Infof("用户登录失败: username=%s", req.Username)
	return &types.LoginRes{
		Code:    401,
		Message: "用户名或密码错误",
	}, nil
}
