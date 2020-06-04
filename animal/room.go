package animal

import (
	"context"
	"errors"
	"github.com/game-explorer/animal-chess-server/app/sessionhub"
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"github.com/game-explorer/animal-chess-server/model"
	"github.com/game-explorer/animal-chess-server/repository"
	"sync"
	"time"
)

// 处理在房间内的所有用户消息
type Room struct {
	players  map[int64]*sessionhub.Session
	receives map[int64]chan model.Message

	rep  repository.Interface
	room model.Room
}

func GetRoom(id int64) (r *Room, exist bool) {
	lock.Lock()
	defer lock.Unlock()

	r, exist = Rooms[id]
	return
}

func CreateRoom(room model.Room) (r *Room) {
	lock.Lock()
	defer lock.Unlock()

	r, exist := Rooms[room.Id]
	if exist {
		return
	}

	r = newRoom(room)
	Rooms[room.Id] = r
	return
}

var (
	ErrTimeOut = errors.New("timeout")
)

func (r *Room) GetPlayerIds() (ids []int64) {
	for k := range r.players {
		ids = append(ids, k)
	}
	return
}

// 接收一个消息, ctx超时控制
func (r *Room) Receive(ctx context.Context, pid int64, validator func(msg *model.Message) bool) (msg model.Message, err error) {
	var f = make(chan struct{})
	go func() {
		defer close(f)
		for {
			msg = <-r.receives[pid]
			if validator != nil {
				if validator(&msg) {
					break
				}
			}
		}
	}()

	select {
	case <-f:
	case <-ctx.Done():
		err = ErrTimeOut
		return
	}

	return
}

// 接收任何一个消息
func (r *Room) ReceiveAll(validator func(msg *model.Message) bool) (msgs []model.Message, err error) {
	var closeOnce sync.Once

	var finish chan bool
	for _, i := range r.players {
		go func(pid int64) {
			select {
			case <-finish:
				return
			case s := <-r.receives[pid]:
				if validator != nil {
					if !validator(&s) {
						break
					}
				}

				msgs = append(msgs, s)
				closeOnce.Do(func() {
					close(finish)
				})
			}
		}(i.PlayerId)
	}

	<-finish

	return
}

func (r *Room) Send(pid int64, messageType model.MessageType, raw interface{}) (err error) {
	err = r.players[pid].Send(messageType, raw)
	return
}

// 当玩家加入房间
// - 发送玩家动作
func (r *Room) Join(playerId int64) (err error) {
	s, exist := sessionhub.Players.Get(playerId)
	if !exist {
		return
	}
	r.players[playerId] = s

	r.receives[playerId] = make(chan model.Message, 10)
	// 接管用户消息
	s.SetReceiver(r.receives[playerId])

	// 发送消息告知其他玩家加入了房间

	err = r.Notify(playerId)
	if err != nil {
		return
	}

	return
}

// 通知用户该进行的动作, 只有目标用户会收到消息
// 用于用户第一次加入房间
func (r *Room) Notify(playerId int64) (err error) {
	ps, _ := r.room.PlayerStatus.Get(playerId)
	// 没有准备就准备
	if !ps.Ready {
		err = r.players[playerId].Send(model.Action, model.ActionMsgRaw{Timeout: 30, Action: model.ActionReady})
		return
	}

	// 是否在游戏中
	switch r.room.Status {
	case model.PlayingStatus:
		// 该谁走棋
		if playerId == r.room.TimeToPlayerId {
			err = r.players[playerId].Send(model.Action, model.ActionMsgRaw{Timeout: 30, Action: model.ActionMove})
			if err != nil {
				return
			}
		}

	case model.EndStatus:
		err = r.players[playerId].Send(model.End, model.EndRaw{})
		if err != nil {
			return
		}
	}

	return
}

// stop: 当游戏结束时返回true
func (r *Room) oneLoop() (stop bool, err error) {
	// 通知走棋
	currPlayer := r.room.TimeToPlayerId
	playerIds := r.GetPlayerIds()

	err = r.Send(currPlayer, model.Action, model.ActionMsgRaw{Action: model.ActionMove})
	if err != nil {
		log.Errorf("send Msg %+v", err)
		return
	}

	// 等待走棋
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	msg, err := r.Receive(ctx, currPlayer, func(msg *model.Message) bool {
		if msg.Type == model.Move {
			return true
		}
		return false
	})
	cancel()
	if err != nil {
		if err == ErrTimeOut {
			// TODO自动走棋
			msg = model.Message{
				Type: model.Move,
			}
		} else {
			log.Errorf("Receive err %+v", err)
			return
		}
	}

	// 处理走棋
	var m model.MoveMsgRaw
	msg.UnmarshalRaw(&m)

	ps, _ := r.room.PlayerStatus.Get(currPlayer)
	fitResult := ""
	if ps.IsP1() {
		fitResult, err = r.room.TablePieces.Move("p1", m.From, m.To)
		if err != nil {
			return
		}
	} else {
		fitResult, err = r.room.TablePieces.Move("p2", m.From, m.To)
		if err != nil {
			return
		}
	}

	// 发送走棋给房间所有人
	rsp := buildRsp(playerIds, model.Message{
		Type: model.Move,
		Raw: buildJson(model.MoveMsgRaw{
			From:      m.From,
			To:        m.To,
			PlayerId:  currPlayer,
			FitResult: fitResult,
		}),
	})

	Response(rsp)

	// 判断游戏结束, 如果没结束就通知下架走棋
	win := r.room.TablePieces.IsWin()
	if win == "" {
		// 通知下家走棋
		next, e := r.room.PlayerStatus.Next(currPlayer)
		if e != nil {
			err = e
			return
		}
		r.room.TimeToPlayerId = next.PlayerId
		Response(buildRsp(playerIds, model.Message{
			Type: model.Action,
			Raw:  buildJson(model.ActionMsgRaw{Action: model.ActionMove}),
		}))
	} else {
		r.room.Status = model.EndStatus
		if win == "p1" {
			winner, _ := r.room.PlayerStatus.GetP1()
			r.room.WinPlayerId = winner.PlayerId
			Response(buildRsp(playerIds, model.Message{
				Type: model.End,
				Raw:  buildJson(model.EndRaw{WinPlayerId: winner.PlayerId}),
			}))
		} else {
			winner, _ := r.room.PlayerStatus.GetP2()
			r.room.WinPlayerId = winner.PlayerId
			Response(buildRsp(playerIds, model.Message{
				Type: model.End,
				Raw:  buildJson(model.EndRaw{WinPlayerId: winner.PlayerId}),
			}))
		}

		// 游戏结束
		stop = true
	}

	err = r.rep.SaveRoom(&r.room)
	if err != nil {
		return
	}

	return
}

func Response(rsp []MessageRsp) {
	for _, v := range rsp {
		session, exist := sessionhub.Players.Get(v.ToPlayerId)
		if exist {
			e := session.Send(v.Msg.Type, v.Msg.Raw)
			if e != nil {
				log.Errorf("SendMessageRaw %v", e)
			}
		} else {
			// 不存在可能是掉线了, 直接跳过发送
			log.Warningf("player %d is offline", v.ToPlayerId)
		}
	}
}

// 开始游戏循环
func (r *Room) Start() {
	// 通知两方游戏开始
	Response(buildRsp(r.GetPlayerIds(), model.Message{
		Type: model.Start,
		Raw:  nil,
	}))

	// 循环直到游戏结束
	for {
		stop, err := r.oneLoop()
		if err != nil {
			log.Errorf("%+v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if stop {
			break
		}
	}
}

// 接受两个人的 准备 消息
func (r *Room) Service() {
	// 循环直到游戏开始
	isWait := true
	for isWait {
		// 如果游戏是开始状态, 则直接开始
		if r.room.Status == model.PlayingStatus {
			break
		}

		msgs, err := r.ReceiveAll(nil)
		if err != nil {
			log.Errorf("%+v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 接受准备消息
		for _, msg := range msgs {
			switch msg.Type {
			case model.SetPiece:
				var m model.SetPieceMsgRaw
				msg.UnmarshalRaw(&m)

				ps, _ := r.room.PlayerStatus.Get(msg.PlayerId)

				if ps.IsP1() {
					err = m.Pieces.ValidateSet("p1")
					if err != nil {
						return
					}
				} else {
					err = m.Pieces.ValidateSet("p2")
					if err != nil {
						return
					}
				}

				err = r.room.PlayerStatus.Ready(msg.PlayerId)
				if err != nil {
					return
				}

				// 更新房间中的棋子
				if ps.IsP1() {
					r.room.TablePieces.P1 = &model.TablePiecesOne{
						Pieces: m.Pieces,
					}
				} else {
					r.room.TablePieces.P2 = &model.TablePiecesOne{
						Pieces: m.Pieces,
					}
				}

				Response(buildRsp(r.GetPlayerIds(), model.Message{
					Type: model.SetPiece,
					Raw: buildJson(model.SetPieceMsgRaw{
						Pieces:   m.Pieces, // 删除对方棋子信息
						PlayerId: msg.PlayerId,
					}),
				}))

				// 如果摆放完成就开始游戏
				if r.room.PlayerStatus.IsAllReady() {
					isWait = false
				}

				// 记得保存房间
				err = r.rep.SaveRoom(&r.room)
				if err != nil {
					return
				}
			}
		}
	}

	// 循环直到游戏结束
	r.Start()
}

func newRoom(room model.Room) *Room {
	r := Room{
		rep:      repository.NewMysql(),
		room:     room,
		receives: map[int64]chan model.Message{},
		players:  map[int64]*sessionhub.Session{},
	}
	go r.Service()
	return &r
}

var Rooms map[int64]*Room
var lock sync.Mutex

func init() {
	Rooms = map[int64]*Room{}
}
