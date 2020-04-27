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
当双方都摆放好棋子之后

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
