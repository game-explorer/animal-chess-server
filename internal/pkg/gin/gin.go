package gin

import (
	"context"
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Engine struct {
	*gin.Engine
	httpServer *http.Server
}

func NewGin(debug bool) *Engine {
	if debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	httpServer := &http.Server{
		Addr:    "",
		Handler: engine,
	}

	return &Engine{
		Engine:     engine,
		httpServer: httpServer,
	}
}

func (g *Engine) Listen(ctx context.Context, addr string, ) {
	go func() {
		<-ctx.Done()
		g.GracefulStop()
	}()

	g.httpServer.Addr = addr
	err := g.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panicf("http.ListenAndServe err:%+v", err)
	}
	return
}

func (g *Engine) GracefulStop() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := g.httpServer.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}
}
