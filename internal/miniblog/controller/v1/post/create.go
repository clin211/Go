package post

import (
	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/errno"
	"github.com/Forest-211/miniblog/internal/pkg/log"
	v1 "github.com/Forest-211/miniblog/pkg/api/miniblog/v1"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ctrl *PostController) Create(c *gin.Context) {
	// 写入日志
	log.C(c).Infow("Create post function called")

	// 获取参数
	var r v1.CreatePostRequest

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
	if err := ctrl.b.Posts().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, map[string]any{
		"title":   r.Title,
		"content": r.Content,
		"author":  "Forest",
		"id":      uuid.New().String(),
	})
}
