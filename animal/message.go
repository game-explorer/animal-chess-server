package animal

import (
	"encoding/json"
	"errors"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/game-explorer/animal-chess-server/repository"
)

type MessageRsp struct {
	ToPlayerId int64
	Msg        model.Message
}

func HandMessage(playerId int64, msg *model.Message) (rsp []MessageRsp, err error) {
	r := repository.NewMysql()

	switch msg.Type {
	case model.CreateRoom:
		// 创建房间
		var roomId int64
		roomId, err = r.CreateRoom(&model.Room{
			PlayerId: playerId,
		})
		if err != nil {
			return
		}
		rsp = buildRsp([]int64{playerId}, model.Message{
			Type: model.CreateRoom,
			Raw:  buildJson(map[string]interface{}{"room_id": roomId}),
		})
		return
	case model.JoinRoom:
		// 加入房间
		// 获取room
		var m model.JoinRoomMsgRaw
		msg.UnmarshalRaw(&m)

		room, exist, e := r.GetRoom(m.RoomId)
		if e != nil {
			err = e
			return
		}
		if !exist {
			err = errors.New("not found room")
			return
		}

		status, e := room.PlayerStatus.Join(playerId)
		if e != nil {
			err = e
			return
		}
		// 当再次进入房间的时候, 需要下发全部的游戏信息: 人员, 游戏状态, 棋子

		// 更新房间
		err = r.SaveRoom(&room)
		if err != nil {
			return
		}

		err = r.UpdatePlayer(&model.Player{
			Id:       playerId,
			InRoomId: room.Id,
		})
		if err != nil {
			return
		}

		// 发送消息给房间内所有人
		ids, e := getPlayerIdsInRoom(r, room.Id)
		if e != nil {
			err = e
			return
		}

		rsp = buildRsp(ids, model.Message{
			Type: model.JoinRoom,
			Raw: buildJson(model.JoinRoomMsgRaw{
				RoomId:   room.Id,
				PlayerId: playerId,
				Camp:     status.Camp,
			}),
		})
	case model.LeaveRoom:
		// 离开房间
		p, exist, e := r.GetPlayer(playerId)
		if e != nil {
			err = e
			return
		}
		if !exist || p.InRoomId == 0 {
			rsp = buildRsp([]int64{playerId}, model.Message{
				Type: model.LeaveRoom,
				Raw: buildJson(model.LeaveRoomMsgRaw{
					PlayerId: playerId,
				}),
			})
		} else {
			// 发送消息给房间内所有人
			ids, e := getPlayerIdsInRoom(r, p.InRoomId)
			if e != nil {
				err = e
				return
			}

			rsp = buildRsp(ids, model.Message{
				Type: model.LeaveRoom,
				Raw: buildJson(model.LeaveRoomMsgRaw{
					PlayerId: playerId,
				}),
			})
		}

		// 更新用户
		if exist {
			err = r.UpdatePlayer(&model.Player{
				Id:       playerId,
				InRoomId: 0,
			})
			if err != nil {
				return
			}
		}
	case model.SetPiece:
		var m model.SetPieceMsgRaw
		msg.UnmarshalRaw(&m)

		p, _, e := r.GetPlayer(playerId)
		if e != nil {
			err = e
			return
		}

		room, exist, e := r.GetRoom(p.InRoomId)
		if e != nil {
			err = e
			return
		}
		if !exist {
			err = errors.New("not found room")
			return
		}

		err = room.PlayerStatus.Ready(playerId, m.Pieces)
		if err != nil {
			return
		}

		err = r.SaveRoom(&room)
		if err != nil {
			return
		}

		// 发送消息给房间内所有人
		ids, e := getPlayerIdsInRoom(r, room.Id)
		if e != nil {
			err = e
			return
		}
		rsp = buildRsp(ids, model.Message{
			Type: model.SetPiece,
			Raw: buildJson(model.SetPieceMsgRaw{
				Pieces:   nil,
				PlayerId: playerId,
			}),
		})

		// 如果摆放完成就开始游戏
		if room.PlayerStatus.IsAllReady() {
			room.Status = model.PlayingStatus
			err = r.SaveRoom(&room)
			if err != nil {
				return
			}

			rsp = append(rsp, buildRsp(ids, model.Message{
				Type: model.Start,
				Raw:  nil,
			})...)

			// 通知谁先手
		}
	}

	return
}

func buildJson(i interface{}) []byte {
	bs, _ := json.Marshal(i)
	return bs
}

func buildRsp(toPlayerIds []int64, msg model.Message) (r []MessageRsp) {
	for _, playerId := range toPlayerIds {
		r = append(r, MessageRsp{
			ToPlayerId: playerId,
			Msg:        msg,
		})
	}

	return
}

func getPlayerIdsInRoom(r repository.Interface, roomId int64) (ids []int64, err error) {
	ps, err := r.GetPlayerByRoomId(roomId)
	if err != nil {
		return
	}

	ids = make([]int64, len(ps))
	for i, v := range ps {
		ids[i] = v.Id
	}
	return
}
