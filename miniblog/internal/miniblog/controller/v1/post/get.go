package post

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/errno"
	"github.com/Forest-211/miniblog/internal/pkg/log"
	v1 "github.com/Forest-211/miniblog/pkg/api/miniblog/v1"
)

func (ctrl *PostController) Get(c *gin.Context) {
	log.C(c).Infow("Get post function called")
	// 获取参数
	var r v1.PostByIDRequest

	r.ID = c.Query("id")

	// 接收 HTTP 请求中的 JSON 数据、并进行参数解析
	if r.ID == "" {
		log.C(c).Errorw("ShouldBindJSON error", "err", "empty id")
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 参数校验
	if _, err := govalidator.ValidateStruct(&r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	if posts, err := ctrl.b.Posts().Get(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	} else {
		core.WriteResponse(c, nil, posts)
	}
}
