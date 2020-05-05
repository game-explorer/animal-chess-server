package model

import (
	"errors"
	"strconv"
	"strings"
)

type TablePieces struct {
	P1 *TablePiecesOne `json:"p1"` // 蓝方的棋子
	P2 *TablePiecesOne `json:"p2"` // 红方的棋子
}

type TablePiecesOne struct {
	Pieces Pieces  `json:"pieces"`
	Die    []Piece `json:"die"` // 死掉的棋子
}

const (
	FitResultAllDie = "bothdie"
	FitResultP1Win  = "p1win"
	FitResultP2Win  = "p2win"
)

// camp: p1/p2
func (tp TablePieces) Move(camp string, from Point, to Point) (fitResult string, err error) {
	attack := tp.P2 // 进攻方
	other := tp.P1

	switch camp {
	case "p1":
		attack = tp.P1
		other = tp.P2
	case "p2":
		other = tp.P1
		attack = tp.P2
	default:
		err = errors.New("bad camp")
		return
	}

	err = attack.Pieces.Move(from, to)
	if err != nil {
		return
	}

	// 检查是否打架
	if o, exist := other.Pieces[to]; exist {
		win, allDie := attack.Pieces[to].Fit(o)
		if allDie {
			fitResult = FitResultAllDie

			var die Piece
			// 删除死掉的棋子
			die, err = attack.Pieces.Remove(to)
			if err != nil {
				return
			}
			attack.Die = append(attack.Die, die)
			die, err = other.Pieces.Remove(to)
			if err != nil {
				return
			}
			other.Die = append(other.Die, die)

			return
		}
		if win {
			if camp == "p1" {
				fitResult = FitResultP1Win
			} else {
				fitResult = FitResultP2Win
			}
			var die Piece
			die, err = other.Pieces.Remove(to)
			if err != nil {
				return
			}
			other.Die = append(other.Die, die)

			return
		} else {
			if camp == "p1" {
				fitResult = FitResultP2Win
			} else {
				fitResult = FitResultP1Win
			}
			var die Piece
			die, err = attack.Pieces.Remove(to)
			if err != nil {
				return
			}
			attack.Die = append(attack.Die, die)
			return
		}
	}
	return
}

// 返回赢方的camp: p1/p2
func (tp TablePieces) IsWin() string {
	if _, exist := tp.P1.Pieces[P2Lair]; exist {
		return "p1"
	}

	if _, exist := tp.P2.Pieces[P1Lair]; exist {
		return "p2"
	}

	return ""
}

type Pieces map[Point]Piece

func (p Pieces) Move(from Point, to Point) error {
	pi, ok := (p)[from]
	if !ok {
		return errors.New("not has point: " + string(from))
	}

	// todo check 点是否能走

	delete(p, from)
	(p)[to] = pi
	return nil
}

func (p Pieces) Remove(from Point) (pi Piece, err error) {
	pi, ok := (p)[from]
	if !ok {
		err = errors.New("not has point: " + string(from))
		return
	}

	delete(p, from)
	return
}

// 检查是否摆放完整
// camp: p1/p2
func (p Pieces) ValidateSet(camp string) (err error) {
	if len(p) != 16 {
		return errors.New("请摆放完全部棋子")
	}
	// 是否每个棋子都是两个
	pCount := map[Piece]int{}
	for po, pi := range p {
		pCount[pi]++

		if camp == "p1" {
			// 是否有不合规的点, 兽血和山洞不能放置棋子
			_, y := po.Int()
			if y == 0 {
				return errors.New("边界不能摆放棋子")
			}
			if po == P1CaveLeft || po == P1CaveRight || po == P1Lair {
				return errors.New("兽穴和山洞不能摆放棋子")
			}
		} else {
			_, y := po.Int()
			if y == 12 {
				return errors.New("边界不能摆放棋子")
			}
			if po == P2CaveLeft || po == P2CaveRight || po == P2Lair {
				return errors.New("兽穴和山洞不能摆放棋子")
			}
		}
	}

	return
}

// 方便存储到数据库, 使用string表示
type Point string

func (p Point) Int() (int, int) {
	ss := strings.Split(string(p), "-")
	x, _ := strconv.ParseInt(ss[0], 10, 64)
	y, _ := strconv.ParseInt(ss[1], 10, 64)
	return int(x), int(y)
}

const (
	P1Lair Point = "0-4"
	P2Lair Point = "12-4"

	P1CaveLeft  Point = "4-2"
	P1CaveRight Point = "4-6"

	P2CaveLeft  Point = "8-2"
	P2CaveRight Point = "8-6"
)

// 数值表示是什么动物
// 0-7分别表示老鼠到大象
type Piece int

func (p1 Piece) Fit(p2 Piece) (win bool, allDie bool) {
	if p1 == 0 && p2 == 7 {
		win = true
		return
	}
	if p1 == 7 && p2 == 0 {
		win = false
		return
	}
	if p1 == p2 {
		allDie = true
		return
	}

	win = p1 > p2
	return
}
