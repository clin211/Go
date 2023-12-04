package post

import (
	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

func (ctrl *PostController) Delete(c *gin.Context) {
	log.C(c).Infow("Create post function called")
	core.WriteResponse(c, nil, map[string]any{
		"title":   "delete test",
		"content": "Let's build a blog.",
	})
}
