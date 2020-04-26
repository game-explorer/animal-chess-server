package repository

import (
	"fmt"
	"github.com/game-explorer/animal-chess-server/model"
)

func CreateRoom(room *model.Room) (roomId int64, err error) {
	_, err = engine.Insert(room)
	if err != nil {
		err = fmt.Errorf("mysql.Insert %w", err)
		return
	}

	roomId = room.Id
	return
}

func GetRoom(roomId int64) (room model.Room, exist bool, err error) {
	_, err = engine.Where("id=?", roomId).Get(&room)
	if err != nil {
		err = fmt.Errorf("mysql.Get %w", err)
		return
	}
	return
}

func SaveRoom(room *model.Room) (err error) {
	_, err = engine.Where("id=?", room.Id).Update(room)
	if err != nil {
		err = fmt.Errorf("mysql.Update %w", err)
		return
	}
	return
}
