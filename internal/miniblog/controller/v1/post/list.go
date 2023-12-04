package post

import (
	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

func (ctrl *PostController) List(c *gin.Context) {
	log.C(c).Infow("List post function called")
	core.WriteResponse(c, nil, map[string]any{
		"title":   "list test",
		"content": "Let's build a blog.",
	})
}
