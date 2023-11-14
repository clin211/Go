package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/errno"
	"github.com/Forest-211/miniblog/internal/pkg/log"
	v1 "github.com/Forest-211/miniblog/pkg/api/miniblog/v1"
)

func (ctrl *UserController) Create(c *gin.Context) {
	// 写入日志
	log.C(c).Infow("Create user function called")

	// 创建用户时，发送的参数
	var r v1.CreateUserRequest

	// 接收 HTTP 请求中的 JSON 数据、并进行参数解析
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 参数校验
	if _, err := govalidator.ValidateStruct(&r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	// 逻辑处理
	// if err := ctrl.b.Users().Create(c, &r); err != nil {
	// 	core.WriteResponse(c, err, nil)
	// 	return
	// }

	// 返回结果
	core.WriteResponse(c, nil, nil)
}
