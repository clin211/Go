package miniblog

import (
	"github.com/gin-gonic/gin"

	"github.com/Forest-211/miniblog/internal/miniblog/controller/v1/user"
	"github.com/Forest-211/miniblog/internal/miniblog/store"
	"github.com/Forest-211/miniblog/internal/pkg/core"
	"github.com/Forest-211/miniblog/internal/pkg/errno"
	"github.com/Forest-211/miniblog/internal/pkg/log"
)

func installRouters(g *gin.Engine) error {
	// register 404 handler
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// register healthz handler
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	uc := user.NewUserController(store.S)

	v1 := g.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("", uc.Create)
		}
	}

	return nil
}
