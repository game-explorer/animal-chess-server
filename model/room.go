package model

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
)

// 房间管理着一局棋局和玩家
type Room struct {
	Id int64 `json:"id" xorm:"pk autoincr int(11)"`
	PlayerStatus PlayerStatus `json:"player_status" xorm:"json"`
}

type PlayerStatus []PlayerStatusOne

type PlayerStatusOne struct {
	PlayerId int64 `json:"player_id"`
	Conn     *websocket.Conn
	Ready    bool   `json:"ready"`
	Camp     string `json:"camp"` // red, blue
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

func (p *PlayerStatus) Ready(playerId int64, camp string) (err error) {
	if len(*p) >= 2 {
		return errors.New("房间人数已满")
	}

	if camp != "red" && camp != "blue" {
		return errors.New("请传递红蓝方")
	}

	for _, v := range *p {
		if v.Camp == camp {
			return fmt.Errorf("%v方已被对方选择, 请选择另一方", camp)
		}
	}

	*p = append(*p, PlayerStatusOne{
		PlayerId: playerId,
		Ready:    true,
		Camp:     camp,
	})
	return
}
