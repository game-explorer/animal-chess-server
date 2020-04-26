package gin

import (
	"github.com/game-explorer/animal-chess-server/lib/gin"
	"github.com/game-explorer/animal-chess-server/service/gin/handler"
)

func New(debug bool) *gin.Engine {
	e := gin.NewGin(debug)

	handler.Ws(e)

	return e
}
