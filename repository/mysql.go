package repository

import (
	"github.com/game-explorer/animal-chess-server/app/config"
	"github.com/game-explorer/animal-chess-server/lib/orm"
)

func InitMysql() (err error) {
	return
}

var engine *orm.RWEngine

func init() {
	dsn := config.App.Mysql.AnimateChess
	var err error
	engine, err = orm.New(dsn, dsn)
	if err != nil {
		panic(err)
	}
}
