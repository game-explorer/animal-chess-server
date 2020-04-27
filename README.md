# animal-chess-server

全部游戏逻辑使用websocket通信

游戏逻辑包括: 新建房间, 加入房间, 开始游戏等

游戏逻辑不包括: 生成playerId, 获取用户信息

优点:
- 不需要双向通信的地方用http协议更简单
