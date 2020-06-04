package sessionhub

import (
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"github.com/game-explorer/animal-chess-server/internal/pkg/ws"
	"github.com/game-explorer/animal-chess-server/model"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

var Players *players

type players struct {
	ps map[int64]*Session
	mu sync.Mutex
}

type Session struct {
	PlayerId int64
	conn     *websocket.Conn
	receiver chan model.Message
}

func (s *Session) Send(messageType model.MessageType, raw interface{}) error {
	return ws.SendMessageRaw(s.conn, messageType, raw)
}

func (s *Session) SetReceiver(r chan model.Message) {
	s.receiver = r
	return
}

func (p *players) Add(playerId int64, Conn *websocket.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 把老的连接挤下线
	if old, exist := p.ps[playerId]; exist {
		old.conn.Close()
	}
	s := &Session{
		PlayerId: playerId,
		conn:     Conn,
		receiver: nil,
	}
	p.ps[playerId] = s

	return
}

// 读
func (p *players) Read(playerId int64) {
	s, exist := p.Get(playerId)
	if !exist {
		return
	}
	for {
		var bs string
		err := websocket.Message.Receive(s.conn, &bs)
		if err != nil {
			log.Errorf("ws Receive err: %v", err)
			break
		}

		var msg model.Message
		err = msg.Unmarshal([]byte(bs))
		if err != nil {
			log.Errorf("Unmarshal err: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if msg.Type == "ping" {
			s.Send("pong", nil)
			continue
		}

		msg.PlayerId = playerId
		select {
		case s.receiver <- msg:
		default:
			s.Send(model.Err, model.ErrorMsgRaw{Msg: "发送消息太快了"})
		}
	}

	log.Infof("c finish %+v", playerId)
}

func (p *players) Del(playerId int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 下线
	if old, exist := p.ps[playerId]; exist {
		old.conn.Close()
	}

	delete(p.ps, playerId)

	return
}

func (p *players) Get(playerId int64) (s *Session, exist bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	s, exist = p.ps[playerId]

	return
}

func newPlayers() *players {
	x := &players{
		ps: map[int64]*Session{},
	}

	return x
}

func init() {
	Players = newPlayers()
}
