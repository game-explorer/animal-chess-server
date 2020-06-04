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

func getRoomByPlayer(r repository.Interface, playerId int64) (room model.Room, exist bool, err error) {
	p, _, e := r.GetPlayer(playerId)
	if e != nil {
		err = e
		return
	}

	room, exist, err = r.GetRoom(p.InRoomId)
	if err != nil {
		return
	}

	return
}
func getRoomByPlayerMust(r repository.Interface, playerId int64) (room model.Room, err error) {
	var exist bool
	room, exist, err = getRoomByPlayer(r, playerId)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("not found room")
		return
	}
	return
}

func getPlayerIdsByPlayer(r repository.Interface, playerId int64) (ids []int64, err error) {
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

	return getPlayerIdsInRoom(r, room.Id)
}
