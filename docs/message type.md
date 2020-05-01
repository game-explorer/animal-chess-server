## WS地址
ws://localhost:9000/ws?player_id=1

## Model & Const
```
const (
	WaitPeopleStatus RoomStatus = 1 // 等待玩家加入
	WaitReadStatus   RoomStatus = 2 // 等待准备
	PlayingStatus    RoomStatus = 3 // 正在游戏中
	EndStatus        RoomStatus = 4 // 游戏结束
)
```

## 上行
#### 创建房间
```
type: create_room
```

#### 加入房间
```
type: join_room
raw: 
  room_id: 1
```
#### 摆放棋子并准备
如果收到join_room消息并且status=2(准备游戏中), 则开始准备(摆放棋子)流程

```
type: set-piece
raw: 
  pieces: 
    "1-2": 1
    "2-3": 0
```
#### 走棋

```
type: move
raw: 
  from: "1-2"
  to: "1-3"
```

## 下行

#### 游戏状态
当用户的ws第一次连接上时 会发送玩家当前状态, 如果是正在房间中的状态则需要询问用户是否回到房间.

```
type: game_status
raw:
  status: 1 // 1: 等待开始游戏, 2: 正在游戏, 0: 没有进行游戏
```

#### 创建房间成功
```
type: create_room
raw: 
  room_id: 1 // 房间号
```

#### 加入房间
自己或者其他玩家加入房间

```
type: join_room
raw:
  player_id: 1 // 加入房间的玩家id，可能是其他玩家
  camp: 'red' // 阵营, red blue
  status: 1 // 当前房间状态
```

#### 摆放棋子并准备
```
type: set-piece
raw: 
  pieces: 
    "1-2": 1
    "2-3": 0
  player_id: 1 
```
#### 开始游戏
当双方都摆放好棋子(准备)之后, 会收到此消息

```
type: start
raw: 
```

#### 该XX走棋

```
type: timeto
raw: 
  player_id: 1
```
