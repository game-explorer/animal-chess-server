package animal

import (
	"errors"
	"github.com/game-explorer/animal-chess-server/app/sessionhub"
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/game-explorer/animal-chess-server/repository"
)

// 处理在大厅(没有加入房间)的用户消息

var TheHall *Hall

type Hall struct {
	msg chan model.Message
	rep repository.Interface
}

func (h *Hall) Join(playerId int64) {
	s, exist := sessionhub.Players.Get(playerId)
	if !exist {
		return
	}

	s.SetReceiver(h.msg)

	h.Notify(playerId)
}

func (h *Hall) Notify(playerId int64) (err error) {
	// 获取玩家状态(是否在房间中), 让前端做对应的动作
	// 如让用户选择是否重新加入房间
	p, exist, err := h.rep.GetPlayer(playerId)
	if err != nil {
		return
	}

	conn, _ := sessionhub.Players.Get(playerId)

	if exist && p.InRoomId != 0 {
		// 发送回到游戏中的询问
		room, _, e := h.rep.GetRoom(p.InRoomId)
		if e != nil {
			err = e
			return
		}

		e = conn.Send(model.GameStatus, model.GameStatusMsgRaw{Status: room.Status, RoomId: room.Id})
		if e != nil {
			log.Errorf("SendMessageRaw %v", e)
		}
	} else {
		e := conn.Send(model.GameStatus, model.GameStatusMsgRaw{Status: 0})
		if e != nil {
			log.Errorf("SendMessageRaw %v", e)
		}
	}
	return
}

// 处理玩家在房间外发送的消息:
// - CreateRoom
// - JoinRoom
func (h *Hall) Service() {
	for msg := range h.msg {
		var err error
		switch msg.Type {
		case model.CreateRoom:
			// 创建房间
			var roomId int64
			roomId, err := h.rep.CreateRoom(&model.Room{
				PlayerId: msg.PlayerId,
				Status:   model.WaitPeopleStatus,
			})
			if err != nil {
				return
			}
			Response(buildRsp([]int64{msg.PlayerId}, model.Message{
				Type: model.CreateRoom,
				Raw:  buildJson(map[string]interface{}{"room_id": roomId}),
			}))
			return
		case model.JoinRoom:
			// 加入房间
			// 获取room
			var m model.JoinRoomMsgRaw
			msg.UnmarshalRaw(&m)

			room, exist, e := h.rep.GetRoom(m.RoomId)
			if e != nil {
				err = e
				return
			}
			if !exist {
				err = errors.New("not found room")
				return
			}

			_, e = room.PlayerStatus.Join(msg.PlayerId)
			if e != nil {
				err = e
				return
			}

			if room.Status == model.WaitPeopleStatus {
				// 准备中
				if room.PlayerStatus.IsFull() {
					room.Status = model.WaitReadStatus
				}
			}

			// 更新房间
			err = h.rep.SaveRoom(&room)
			if err != nil {
				return
			}

			err = h.rep.UpdatePlayer(&model.Player{
				Id:       msg.PlayerId,
				InRoomId: room.Id,
			})
			if err != nil {
				return
			}

			ro := CreateRoom(room)
			err = ro.Join(msg.PlayerId)
			if err != nil {
				return
			}
		}
	}
}

func init() {
	TheHall = &Hall{
		msg: make(chan model.Message, 10),
		rep: repository.NewMysql(),
	}

	go TheHall.Service()
}
