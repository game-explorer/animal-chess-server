package handler

import (
	"github.com/game-explorer/animal-chess-server/animal"
	"github.com/game-explorer/animal-chess-server/app/sessionhub"
	"github.com/game-explorer/animal-chess-server/lib/log"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"strconv"
	"strings"
	"time"
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

			// 获取玩家状态(是否在房间中), 让前端做对应的动作
			// 如让用户选择是否重新加入房间

			stop := make(chan struct{})

			// 使用chan接受客户端消息
			var msgs = make(chan model.Message, 1)
			go func() {
				for {
					var bs string
					err := websocket.Message.Receive(conn, &bs)
					if err != nil {
						es := err.Error()
						if strings.Contains(es, "EOF") || strings.Contains(es, "use of closed network connection") {
							break
						}
						log.Errorf("ws Receive err: %v", err)
						time.Sleep(1 * time.Second)
						continue
					}

					var msg model.Message
					err = msg.Unmarshal([]byte(bs))
					if err != nil {
						log.Errorf("Unmarshal err: %v", err)
						time.Sleep(1 * time.Second)
						continue
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
						e := websocket.Message.Send(conn, string(model.NewErrorMsg(err).Marshal()))
						if e != nil {
							log.Errorf("websocket.Message.Send %v", e)
						}
						continue
					}
					for _, v := range rsp {
						session, exist := sessionhub.Player.Get(v.ToPlayerId)
						if exist {
							e := websocket.Message.Send(session.Conn, string(v.Msg.Marshal()))
							if e != nil {
								log.Errorf("websocket.Message.Send %v", e)
							}
						} else {
							// 不存在可能是掉线了, 直接跳过发送
							log.Warningf("player %s is offline", v.ToPlayerId)
						}
					}
				}
			}

		}, ctx.Writer, ctx.Request)
	})
}
