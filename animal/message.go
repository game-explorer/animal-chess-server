package animal

import (
	"encoding/json"
	"errors"
	"github.com/game-explorer/animal-chess-server/lib/log"
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
	case model.CreateRoom: // 创建房间
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
	case model.JoinRoom: // 加入房间
		// 获取room
		var m model.JoinRoomMsgRawIn
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

		err = room.PlayerStatus.Join(playerId)
		if err != nil {
			return
		}

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

		log.Infof("ids: %+v", ids)

		rsp = buildRsp(ids, model.Message{
			Type: model.JoinRoom,
			Raw: buildJson(model.JoinRoomMsgRawOut{
				RoomId:   room.Id,
				PlayerId: playerId,
			}),
		})
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
