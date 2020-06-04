package handler

import (
	"github.com/game-explorer/animal-chess-server/animal"
	"github.com/game-explorer/animal-chess-server/app/sessionhub"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"strconv"
)

func Ws(root gin.IRouter) {
	root.GET("/ws", func(ctx *gin.Context) {
		websocket.Handler(Handle).ServeHTTP(ctx.Writer, ctx.Request)
	})
}

func Handle(conn *websocket.Conn) {
	req := conn.Request()
	// 如果头部没有playerId就关闭连接
	playerIdStr := req.URL.Query().Get("player_id")
	if playerIdStr == "" {
		playerIdStr = req.Header.Get("PLAYER_ID")
	}

	playerId, _ := strconv.ParseInt(playerIdStr, 10, 54)

	if playerId == 0 {
		_ = conn.Close()
		return
	}
	sessionhub.Players.Add(playerId, conn)
	animal.TheHall.Join(playerId)

	defer sessionhub.Players.Del(playerId)

	sessionhub.Players.Read(playerId)

}
