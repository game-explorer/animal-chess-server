package model

import "encoding/json"

type Message struct {
	Type MessageType
	Raw  json.RawMessage `json:"raw"`
}

type MessageType string

const (
	Login MessageType = "login" // 登录, 上传PlayerId
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

type LoginMsg struct {
	PlayerId int64 `json:"player_id"`
}

func (m *Message) To(i interface{}) (err error) {
	err = json.Unmarshal(m.Raw, i)
	if err != nil {
		return
	}
	return
}
