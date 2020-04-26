package model

// 玩家保存则正在进行的游戏房间
type Player struct {
	Id       int64 `json:"id" xorm:"int(11) index"`
	InRoomId int64 `json:"in_room_id" xorm:"int(11) index"`

	// 积分
}
