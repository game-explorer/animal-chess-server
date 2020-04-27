package model

import "encoding/json"

type Message struct {
	Type MessageType     `json:"type"`
	Raw  json.RawMessage `json:"raw"`
}

type MessageType string

const (
	Err MessageType = "err" // 错误
	// 玩家动作
	CreateRoom MessageType = "create_room" // 创建房间
	JoinRoom   MessageType = "join_room"   // 加入房间
	Ready      MessageType = "ready"       // 选择阵营
	SetPiece   MessageType = "set-piece"   // 摆放棋子
	Move       MessageType = "move"        // 移动棋子

	// 系统发送的消息
	Fit MessageType = "fit" // 两个棋子打架了, 消息体中包含结果
	End MessageType = "end" // 结束, 消息体包含结果(谁胜利了)
)

type ErrorMsgRaw struct {
	Msg string `json:"msg"`
}

type LoginMsgRaw struct {
	PlayerId int64 `json:"player_id"`
}

type JoinRoomMsgRawIn struct {
	RoomId int64 `json:"room_id"`
}

type JoinRoomMsgRawOut struct {
	RoomId   int64 `json:"room_id"`
	PlayerId int64 `json:"player_id"`
}

func NewErrorMsg(err error) Message {
	e := Message{Type: Err}
	e.MarshalRaw(ErrorMsgRaw{Msg: err.Error()})
	return e
}

func (m *Message) UnmarshalRaw(i interface{}) (err error) {
	err = json.Unmarshal(m.Raw, i)
	if err != nil {
		return
	}
	return
}

func (m *Message) MarshalRaw(i interface{}) {
	m.Raw, _ = json.Marshal(i)
	return
}

func (m *Message) Unmarshal(bs []byte) (err error) {
	err = json.Unmarshal(bs, m)
	if err != nil {
		return
	}
	return
}

func (m Message) Marshal() (bs []byte) {
	bs, _ = json.Marshal(m)
	return
}
