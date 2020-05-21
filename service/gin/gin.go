package gin

import (
	"github.com/game-explorer/animal-chess-server/internal/pkg/gin"
	"github.com/game-explorer/animal-chess-server/service/gin/handler"
)

func New(debug bool) *gin.Engine {
	e := gin.NewGin(debug)

	handler.Ws(e)
	handler.Login(e)

	return e
}
