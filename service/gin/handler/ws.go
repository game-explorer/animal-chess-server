package handler

import (
	"github.com/game-explorer/animal-chess-server/animal"
	"github.com/game-explorer/animal-chess-server/app/sessionhub"
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"github.com/game-explorer/animal-chess-server/internal/pkg/ws"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/game-explorer/animal-chess-server/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"strconv"
	"strings"
	"time"
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

	// 获取玩家状态(是否在房间中), 让前端做对应的动作
	// 如让用户选择是否重新加入房间
	r := repository.NewMysql()
	p, exist, err := r.GetPlayer(playerId)
	if err != nil {
		_ = conn.Close()
		return
	}
	if exist && p.InRoomId != 0 {
		// 发送回到游戏中的询问
		room, _, err := r.GetRoom(p.InRoomId)
		if err != nil {
			log.Errorf("%v", err)
			return
		}

		e := ws.SendMessageRaw(conn, model.GameStatus, model.GameStatusMsgRaw{Status: room.Status, RoomId: room.Id})
		if e != nil {
			log.Errorf("SendMessageRaw %v", e)
		}
	} else {
		e := ws.SendMessageRaw(conn, model.GameStatus, model.GameStatusMsgRaw{Status: 0})
		if e != nil {
			log.Errorf("SendMessageRaw %v", e)
		}
	}

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
		case "ping":
			ws.SendMessage(conn, &model.Message{
				Type: "pong",
				Raw:  nil,
			})
		default:
			rsp, err := animal.HandMessage(r, playerId, &msg)
			if err != nil {
				e := ws.SendMessage(conn, model.NewErrorMsg(err))
				if e != nil {
					log.Errorf("SendMessageRaw %v", e)
				}
				continue
			}
			for _, v := range rsp {
				session, exist := sessionhub.Player.Get(v.ToPlayerId)
				if exist {
					e := ws.SendMessage(session.Conn, &v.Msg)
					if e != nil {
						log.Errorf("SendMessageRaw %v", e)
					}
				} else {
					// 不存在可能是掉线了, 直接跳过发送
					log.Warningf("player %d is offline", v.ToPlayerId)
				}
			}
		}
	}

}
