package animal

import (
	"encoding/json"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/game-explorer/animal-chess-server/repository"
)

type MessageRsp struct {
	ToPlayerId int64
	Msg        model.Message
}

func HandMessage(playerId int64, msg *model.Message) (rsp []MessageRsp, err error) {
	switch msg.Type {
	case model.CreateRoom:
		roomId, err := repository.CreateRoom(&model.Room{
			PlayerId: playerId,
		})
		if err != nil {
			return
		}
		rsp = []MessageRsp{{
			ToPlayerId: playerId,
			Msg: model.Message{
				Type: model.CreateRoom,
				Raw:  buildJson(map[string]interface{}{"id": roomId}),
			},
		}}
		return

	}

	return
}

func buildJson(i map[string]interface{}) []byte {
	bs, _ := json.Marshal(i)
	return bs
}
