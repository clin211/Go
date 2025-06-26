package user

import (
	"github.com/gin-gonic/gin"

	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/log"
)

func (*UserController) List(c *gin.Context) {
	log.C(c).Infow("List user function called")
	core.WriteResponse(c, nil, map[string]interface{}{
		"username": "test list user",
	})
}
