package post

import (
	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/errno"
	"github.com/Forest-211/miniblog/internal/pkg/log"
	v1 "github.com/Forest-211/miniblog/pkg/api/miniblog/v1"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func (ctrl *PostController) List(c *gin.Context) {
	// 写入日志
	log.C(c).Infow("List post function called")

	// 获取参数
	var r v1.ListPostRequest

	// 接收 HTTP 请求中的 JSON 数据、并进行参数解析
	if err := c.ShouldBindJSON(&r); err != nil {
		log.C(c).Errorw("ShouldBindJSON error", "err", err)
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 参数校验
	if _, err := govalidator.ValidateStruct(&r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	// 逻辑处理
	posts, err := ctrl.b.Posts().List(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, posts)
}
