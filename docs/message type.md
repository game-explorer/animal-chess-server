## WS地址
ws://localhost:9000/ws?player_id=1

## Model & Const
```
const (
	WaitPeopleStatus RoomStatus = 1 // 等待玩家加入
	WaitReadStatus   RoomStatus = 2 // 等待准备
	PlayingStatus    RoomStatus = 3 // 正在游戏中11
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
type: set_piece
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
  status: 1 // RoomStatus
```

#### 获取ROOM
用于断线重连
```
type: get_room
raw: 
  status: 2
  player_status: 
    - player_id: 1
      read: true
      camp: "p1"
    - player_id: 2
      read: true
      camp: "p2"
  table_pieces: // 棋子
    p1: // p1的棋子
      pieces:
        "1-2": 1
        "2-3": 0
      die: [0, 1, 2] // p1死掉的棋子
    p2: // p2的棋子
      pieces:
        "1-2": 1
        "2-3": 0
      die: [0, 1, 2] // p2死掉的棋子

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
  camp: 'p2' // 阵营, p1/p2
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
type: time_to
raw: 
  player_id: 1
```

#### 走棋结果

```
type: move
raw: 
  from: "1-2"
  to: "1-3"
  player_id: 1
  fit_result: 1 // 打架结果 bothdie/p1win/p2win, 分别表示都死亡/p1赢/p2赢
```

#### 游戏结束

```
type: end
raw: 
  win_player_id: 1
```

#### 执行动作
服务端命令并等待客户端执行动作, 动作目前有两个: 准备 / 下棋

```
type: action
raw: 
  timeout: 30
  type: move / ready
```
