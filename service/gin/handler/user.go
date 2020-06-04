package handler

import (
	"github.com/game-explorer/animal-chess-server/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(root gin.IRouter) {
	api := root.Group("/api/v1")

	api.GET("/login", func(ctx *gin.Context) {
		uid, _ := ctx.Cookie("uid")
		if uid == "" {
			uid = repository.GenUid()
		}

		r := repository.NewMysql()
		p, err := r.GetOrCreatePlayer(uid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"msg": err,
			})
			return
		}

		ctx.SetCookie("uid", uid, 24*3600*356, "/", ctx.Request.URL.Hostname(), false, true)
		ctx.JSON(200, p)
	})
}
