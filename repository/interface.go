package repository

import "github.com/game-explorer/animal-chess-server/model"

type Interface interface {
	CreateRoom(room *model.Room) (roomId int64, err error)
	GetRoom(roomId int64) (room model.Room, exist bool, err error)
	SaveRoom(room *model.Room) (err error)

	UpdatePlayer(p *model.Player) (err error)
	GetPlayerByRoomId(roomId int64) (r []model.Player, err error)
	GetPlayer(playerId int64) (r model.Player, exist bool, err error)

	GetOrCreatePlayer(uid string) (r model.Player, err error)
}
