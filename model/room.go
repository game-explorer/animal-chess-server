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
	Status        RoomStatus   `json:"status" xorm:"tinyint(1)"`       // 游戏状态
	TablePieces   TablePieces  `json:"table_pieces" xorm:"json"`       // 当前桌面上的棋子, 包括两方
}

type RoomStatus int64

const (
	WaitPeopleStatus RoomStatus = 1 // 等待玩家加入
	WaitReadStatus   RoomStatus = 2 // 等待准备
	PlayingStatus    RoomStatus = 3 // 正在游戏中
	EndStatus        RoomStatus = 4 // 游戏结束
)

type TablePieces struct {
	P1    Pieces  `json:"p1"`     // 蓝方的棋子
	P2    Pieces  `json:"p2"`     // 红方的棋子
	P1Die []Piece `json:"p1_die"` // 蓝方死掉的棋子
	P2Die []Piece `json:"p2_die"` // 红方死掉的棋子
}

// 第一个人就是蓝色方, 第二个是红色方
type PlayerStatus []*PlayerStatusOne

type PlayerStatusOne struct {
	PlayerId int64  `json:"player_id"`
	Ready    bool   `json:"ready"`
	Camp     string `json:"camp"` // red, blue 第一个进入房间的是blue
}

func (p PlayerStatusOne) IsP1() bool {
	return p.Camp == "blue"
}

type Pieces map[Point]Piece

func (p *Pieces) Move(from Point, to Point) error {
	pi, ok := (*p)[from]
	if !ok {
		return errors.New("not has point: " + string(from))
	}

	// todo check 点是否能走

	delete(*p, from)
	(*p)[to] = pi
	return nil
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

func (p PlayerStatus) IsFull() bool {
	if len(p) != 2 {
		return false
	}

	return true
}

// 加入房间, 不是准备状态
func (p *PlayerStatus) Join(playerId int64) (s PlayerStatusOne, err error) {
	_, exist := p.Get(playerId)
	if exist {
		return
	}

	if p.IsFull() {
		err = errors.New("房间人数已满")
		return
	}
	camp := "blue"
	if len(*p) == 1 {
		camp = "red"
	}

	s = PlayerStatusOne{
		PlayerId: playerId,
		Ready:    false,
		Camp:     camp,
	}
	*p = append(*p, &s)
	return
}

// 离开房间
func (p *PlayerStatus) Leave(playerId int64) (err error) {
	x := PlayerStatus{}
	for _, v := range *p {
		if v.PlayerId == playerId {
			continue
		}
		x = append(x, v)
	}

	*p = x
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
func (p *PlayerStatus) Ready(playerId int64, ) (err error) {
	if one, exist := p.Get(playerId); exist {
		one.Ready = true
	}
	return
}
