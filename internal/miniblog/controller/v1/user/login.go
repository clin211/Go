package user

import (
	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/errno"
	"github.com/Forest-211/miniblog/internal/pkg/log"
	v1 "github.com/Forest-211/miniblog/pkg/api/miniblog/v1"
	"github.com/gin-gonic/gin"
)

// 登录 miniblog 并返回一个 JWT Token.
func (ctrl *UserController) Login(c *gin.Context) {
	log.C(c).Infow("Login function called")

	var r v1.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}
