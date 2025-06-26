package user

import (
	"github.com/gin-gonic/gin"

	"github.com/Forest-211/miniblog/internal/pkg/core"
)

func (*UserController) Update(c *gin.Context) {
	core.WriteResponse(c, nil, map[string]any{
		"username": "test update user",
	})
}
