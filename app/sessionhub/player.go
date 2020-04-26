package sessionhub

import (
	"golang.org/x/net/websocket"
	"sync"
)

var Player *player

type player struct {
	ps map[int64]*Session
	mu sync.Mutex
}

type Session struct {
	PlayerId int64
	Conn     *websocket.Conn
}

func (p *player) Add(playerId int64, Conn *websocket.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 把老的连接挤下线
	if old, exist := p.ps[playerId]; exist {
		old.Conn.Close()
	}

	p.ps[playerId] = &Session{
		PlayerId: playerId,
		Conn:     Conn,
	}

	return
}

func (p *player) Del(playerId int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 下线
	if old, exist := p.ps[playerId]; exist {
		old.Conn.Close()
	}

	delete(p.ps, playerId)

	return
}

func (p *player) Get(playerId int64) (s *Session, exist bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	s, exist = p.ps[playerId]

	return
}

func newPlayer() *player {
	return &player{
		ps: map[int64]*Session{},
	}
}

func init() {
	Player = newPlayer()
}
