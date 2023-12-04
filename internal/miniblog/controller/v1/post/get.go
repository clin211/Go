package post

import (
	"github.com/gin-gonic/gin"

	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/log"
)

func (ctrl *PostController) Get(c *gin.Context) {
	log.C(c).Infow("Get post function called")
	core.WriteResponse(c, nil, map[string]any{
		"title":   "get test",
		"content": "Let's build a blog.",
	})
}
