package repository

import (
	"fmt"
	"github.com/game-explorer/animal-chess-server/app/config"
	"github.com/game-explorer/animal-chess-server/internal/pkg/orm"
	util "github.com/game-explorer/animal-chess-server/internal/pkg/rand"
	"github.com/game-explorer/animal-chess-server/model"
	"strings"
	"xorm.io/core"
)

type Mysql struct {
	engine *orm.RWEngine
}

func (m Mysql) CreateRoom(room *model.Room) (roomId int64, err error) {
	room.Status = model.WaitPeopleStatus
	_, err = m.engine.Insert(room)
	if err != nil {
		err = fmt.Errorf("mysql.Insert %w", err)
		return
	}

	roomId = room.Id
	return
}

func (m Mysql) GetRoom(roomId int64) (room model.Room, exist bool, err error) {
	exist, err = m.engine.Where("id=?", roomId).Get(&room)
	if err != nil {
		err = fmt.Errorf("mysql.Get %w", err)
		return
	}
	return
}

func (m Mysql) SaveRoom(room *model.Room) (err error) {
	_, err = m.engine.Where("id=?", room.Id).Update(room)
	if err != nil {
		err = fmt.Errorf("mysql.Update Room %w", err)
		return
	}
	return
}

func (m Mysql) UpdatePlayer(p *model.Player) (err error) {
	exist, err := m.engine.Where("id=?", p.Id).Exist(&model.Player{})
	if err != nil {
		err = fmt.Errorf("mysql.Exist User %w", err)
		return
	}

	// 存在就更新
	if exist {
		_, err = m.engine.Where("id=?", p.Id).Update(p)
		if err != nil {
			err = fmt.Errorf("mysql.Update Player %w", err)
			return
		}

		return
	}

	// 不存在就新建
	p.Uid = GenUid()
	_, err = m.engine.Insert(p)
	if err != nil {
		err = fmt.Errorf("mysql.Insert Player %w", err)
		return
	}

	return
}

func GenUid() string {
	return util.RandomString(32)
}

func (m Mysql) GetPlayerByRoomId(roomId int64) (r []model.Player, err error) {
	err = m.engine.Where("in_room_id=?", roomId).Find(&r)
	if err != nil {
		err = fmt.Errorf("mysql.Find User %w", err)
		return
	}

	return
}

func (m Mysql) GetPlayer(playerId int64) (r model.Player, exist bool, err error) {
	exist, err = m.engine.Where("id=?", playerId).Get(&r)
	if err != nil {
		err = fmt.Errorf("mysql.Get User %w", err)
		return
	}

	return
}

func (m Mysql) GetOrCreatePlayer(uid string) (r model.Player, err error) {
	exist, err := m.engine.Where("uid=?", uid).Get(&r)
	if err != nil {
		err = fmt.Errorf("mysql.Get User %w", err)
		return
	}
	if exist {
		return
	}

	r.Uid = uid
	_, err = m.engine.Insert(&r)
	if err != nil {
		err = fmt.Errorf("mysql.Insert User %w", err)
		return
	}

	return
}

func NewMysql() Interface {
	dsn := config.App.Mysql.AnimalChess
	var err error
	engine, err := orm.New(dsn, dsn)
	if err != nil {
		panic(err)
	}

	return Mysql{
		engine: engine,
	}
}

func InitMysql() (err error) {
	dsn := config.App.Mysql.AnimalChess

	// 创建库
	db, err := core.Open("mysql", strings.Split(dsn, "/")[0]+"/mysql")
	if err != nil {
		return err
	}

	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("create database if not exists `%s` default character set utf8mb4 collate utf8mb4_unicode_ci;", strings.Split(dsn, "/")[1]))
	if err != nil {
		return err
	}

	// 创建表
	engine, err := orm.New(dsn, dsn)
	if err != nil {
		return
	}

	engine.ShowSQL(true)
	engine.ShowExecTime(true)

	err = engine.Sync2(&model.Room{}, &model.Player{})
	if err != nil {
		return
	}
	return
}
