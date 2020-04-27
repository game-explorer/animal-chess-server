package ws

import (
	"github.com/game-explorer/animal-chess-server/model"
	"golang.org/x/net/websocket"
)

func SendMessage(c *websocket.Conn, msg *model.Message) (err error) {
	err = websocket.Message.Send(c, string(msg.Marshal()))
	if err != nil {
		return
	}

	return
}

func SendMessageRaw(c *websocket.Conn, types model.MessageType, raw interface{}) (err error) {
	msg := model.Message{
		Type: types,
		Raw:  nil,
	}
	err = msg.UnmarshalRaw(raw)
	if err != nil {
		return
	}

	err = websocket.Message.Send(c, string(msg.Marshal()))
	if err != nil {
		return
	}

	return
}
