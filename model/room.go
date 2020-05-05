package model

import (
	"errors"
)

// 房间管理着一局棋局和玩家
type Room struct {
	Id             int64        `json:"id" xorm:"pk autoincr int(11)"`
	PlayerId       int64        `json:"player_id" xorm:"int(11)"` // 房主
	PlayerStatus   PlayerStatus `json:"player_status" xorm:"json"`
	TimeToPlayerId int64        `json:"timeto_player_id" xorm:"int(11)"` // 该谁走棋
	WinPlayerId    int64        `json:"win_player_id" xorm:"int(11)"`    // 胜利方
	Status         RoomStatus   `json:"status" xorm:"tinyint(1)"`        // 游戏状态
	TablePieces    TablePieces  `json:"table_pieces" xorm:"json"`        // 当前桌面上的棋子, 包括两方
}

type RoomStatus int64

const (
	WaitPeopleStatus RoomStatus = 1 // 等待玩家加入
	WaitReadStatus   RoomStatus = 2 // 等待准备
	PlayingStatus    RoomStatus = 3 // 正在游戏中
	EndStatus        RoomStatus = 4 // 游戏结束
)

// 第一个人就是蓝色方, 第二个是红色方
type PlayerStatus []*PlayerStatusOne

type PlayerStatusOne struct {
	PlayerId int64  `json:"player_id"`
	Ready    bool   `json:"ready"`
	Camp     string `json:"camp"` // p1, p2 第一个进入房间的是p1
}

func (po PlayerStatusOne) IsP1() bool {
	return po.Camp == "p1" || po.Camp == "blue"
}

func (ps PlayerStatus) IsAllReady() bool {
	if len(ps) != 2 {
		return false
	}

	for _, v := range ps {
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
func (ps *PlayerStatus) Join(playerId int64) (s *PlayerStatusOne, err error) {
	s, exist := ps.Get(playerId)
	if exist {
		return
	}

	if ps.IsFull() {
		err = errors.New("房间人数已满")
		return
	}
	camp := "p2"
	if len(*ps) == 1 {
		camp = "p1"
	}

	s = &PlayerStatusOne{
		PlayerId: playerId,
		Ready:    false,
		Camp:     camp,
	}
	*ps = append(*ps, s)
	return
}

// 离开房间
func (ps *PlayerStatus) Leave(playerId int64) (err error) {
	x := PlayerStatus{}
	for _, v := range *ps {
		if v.PlayerId == playerId {
			continue
		}
		x = append(x, v)
	}

	*ps = x
	return
}

func (ps *PlayerStatus) Get(playerId int64) (*PlayerStatusOne, bool) {
	for _, v := range *ps {
		if v.PlayerId == playerId {
			return v, true
		}
	}

	return nil, false
}

func (ps *PlayerStatus) GetP1() (*PlayerStatusOne, bool) {
	if len(*ps) < 1 {
		return nil, false
	}

	return (*ps)[0], true
}

func (ps *PlayerStatus) GetP2() (*PlayerStatusOne, bool) {
	if len(*ps) < 2 {
		return nil, false
	}

	return (*ps)[1], true
}

// 下一个人
func (ps *PlayerStatus) Next(playerId int64) (po *PlayerStatusOne, err error) {
	p, exist := ps.Get(playerId)
	if !exist {
		return nil, errors.New("not found playerId")
	}

	if p.IsP1() {
		po, _ = ps.GetP2()
	} else {
		po, _ = ps.GetP1()
	}
	return
}

// 安放棋子并准备开始
func (ps *PlayerStatus) Ready(playerId int64, ) (err error) {
	if one, exist := ps.Get(playerId); exist {
		one.Ready = true
	}
	return
}
