package handler

import (
	util "github.com/game-explorer/animal-chess-server/internal/pkg/rand"
	"github.com/game-explorer/animal-chess-server/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(root gin.IRouter) {
	api := root.Group("/api/v1")

	api.GET("/login", func(ctx *gin.Context) {
		uid, _ := ctx.Cookie("uid")
		if uid == "" {
			uid = genUid()
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

func genUid() string {
	return util.RandomString(32)
}
