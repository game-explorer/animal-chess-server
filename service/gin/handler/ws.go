package handler

import (
	"github.com/game-explorer/animal-chess-server/animal"
	"github.com/game-explorer/animal-chess-server/app/sessionhub"
	"github.com/game-explorer/animal-chess-server/lib/log"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"strconv"
)

func Ws(root gin.IRouter) {
	root.GET("/ws", func(ctx *gin.Context) {
		websocket.Handler.ServeHTTP(func(conn *websocket.Conn) {
			// 如果头部没有playerId就关闭连接
			playerIdStr, _ := ctx.GetQuery("player_id")
			if playerIdStr == "" {
				playerIdStr = ctx.GetHeader("PLAYER_ID")
			}

			playerId, _ := strconv.ParseInt(playerIdStr, 10, 54)

			if playerId == 0 {
				_ = conn.Close()
				return
			}

			stop := make(chan struct{})

			// 使用chan接受客户端消息
			var msgs = make(chan model.Message, 1)
			go func() {
				for {
					var msg model.Message
					err := websocket.Message.Receive(conn, &msg)
					if err != nil {
						log.Errorf("ws Receive err: %v", err)
						// TODO 根据错误类型判断
						//time.Sleep(1 * time.Second)
						break
					}

					select {
					case msgs <- msg:
					case <-stop:
						break
					}
				}

				close(msgs)
			}()

			sessionhub.Player.Add(playerId, conn)

			defer func() {
				sessionhub.Player.Del(playerId)
			}()

			// 处理消息
			for msg := range msgs {
				switch msg.Type {
				default:
					rsp, err := animal.HandMessage(playerId, &msg)
					if err != nil {
						log.Errorf("HandMessage err: %v, raw:% s", err, msg.Raw)
						continue
					}
					for _, v := range rsp {
						v.ToPlayerId
					}
				}
			}

		}, ctx.Writer, ctx.Request)
	})
}
