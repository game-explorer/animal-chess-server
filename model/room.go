package model

import (
	"errors"
)

// 房间管理着一局棋局和玩家
type Room struct {
	Id            int64        `json:"id" xorm:"pk autoincr int(11)"`
	PlayerId      int64        `json:"player_id" xorm:"int(11)"` // 房主
	PlayerStatus  PlayerStatus `json:"player_status" xorm:"json"`
	FirstPlayerId int64        `json:"first_player_id" xorm:"int(11)"` // 先手的人的id
}

// 第一个人就是蓝色方, 第二个是红色方
type PlayerStatus []*PlayerStatusOne

type PlayerStatusOne struct {
	PlayerId int64  `json:"player_id"`
	Ready    bool   `json:"ready"`
	Pieces   Pieces `json:"pieces"`
}

type Pieces struct {
	P map[Point]Piece `json:"p"`
}

// 方便存储到数据库, 使用string表示
type Point string

func (p Point) Int() (int, int) {
	return 0, 0
}

// 数值表示是什么动物
// 0-7分别表示老鼠到大象
type Piece int

// 申请先手
// 如果没有申请先手就会随机一个先手
func (p *Room) First(playerId int64) (firstPlayerId int64, err error) {
	if p.FirstPlayerId == 0 {
		p.FirstPlayerId = playerId
	} else {
		p.FirstPlayerId = 0
	}

	firstPlayerId = p.FirstPlayerId
	return
}

func (p PlayerStatus) IsAllReady() bool {
	if len(p) != 2 {
		return false
	}

	for _, v := range p {
		if !v.Ready {
			return false
		}
	}
	return true
}

// 加入房间, 不是准备状态
func (p *PlayerStatus) Join(playerId int64) (err error) {
	_, exist := p.Get(playerId)
	if exist {
		return
	}

	if len(*p) >= 2 {
		return errors.New("房间人数已满")
	}

	*p = append(*p, &PlayerStatusOne{
		PlayerId: playerId,
		Ready:    true,
	})
	return
}

func (p *PlayerStatus) Get(playerId int64) (*PlayerStatusOne, bool) {
	for _, v := range *p {
		if v.PlayerId == playerId {
			return v, true
		}
	}

	return nil, false
}

// 安放棋子并准备开始
func (p *PlayerStatus) Ready(playerId int64, pi Pieces) (err error) {
	if one, exist := p.Get(playerId); exist {
		one.Pieces = pi
		one.Ready = true
	}
	return
}
