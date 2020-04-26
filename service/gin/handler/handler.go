package handler

import (
	"github.com/game-explorer/animal-chess-server/animal"
	"github.com/game-explorer/animal-chess-server/app/sessionhub"
	"github.com/game-explorer/animal-chess-server/lib/log"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"time"
)

func route(root gin.IRouter) {
	root.GET("/ws", func(ctx *gin.Context) {
		websocket.Handler.ServeHTTP(func(conn *websocket.Conn) {
			stop := make(chan struct{})

			// 使用chan接受客户端消息
			var msgs = make(chan model.Message, 1)
			go func() {
				for {
					var msg model.Message
					err := websocket.Message.Receive(conn, &msg)
					if err != nil {
						log.Errorf("ws Receive err: %v", err)
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

			// 6s内没有登录就关闭
			logined := make(chan struct{})
			go func() {
				select {
				case <-time.After(6 * time.Second):
					close(stop)
				case <-logined:
				}
			}()
			var playerId int64

			defer func() {
				if playerId != 0 {
					sessionhub.Player.Del(playerId)

				}
			}()
			// 处理消息
			for msg := range msgs {
				switch msg.Type {
				case model.Login:
					close(logined)

					var m model.LoginMsg
					err := msg.To(&m)
					if err != nil {
						log.Errorf("msg.To err: %v, raw:% s", err, msg.Raw)
						continue
					}
					playerId = m.PlayerId
					sessionhub.Player.Add(playerId, conn)
				default:
					err:=animal.HandMessage(&msg)
					if err != nil {
						log.Errorf("HandMessage err: %v, raw:% s", err, msg.Raw)
						continue
					}
				}
			}

		}, ctx.Writer, ctx.Request)
	})
}
