package main

import (
	"github.com/game-explorer/animal-chess-server/app/config"
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"github.com/game-explorer/animal-chess-server/internal/pkg/signal"
	"github.com/game-explorer/animal-chess-server/repository"
	"github.com/game-explorer/animal-chess-server/service/gin"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
)

func main() {
	c := cli.NewApp()
	c.Name = "animal-chess"
	c.Usage = ""
	c.Version = ""
	c.Action = func(*cli.Context) error {
		ctx, cancel := signal.NewTermContext()
		defer cancel()

		err := repository.InitMysql()
		if err != nil {
			return err
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			s := gin.New(config.App.Debug)
			log.Infof("http listen '%s'", config.App.HttpAddr)
			s.Listen(ctx, config.App.HttpAddr)
			log.Info("http shutdown")
			wg.Done()
		}()

		wg.Wait()

		return nil
	}

	c.Commands = []*cli.Command{
		{
			Usage: "初始化整个项目(上线使用)",
			Name:  "init",
			Action: func(*cli.Context) error {
				err := repository.InitMysql()
				if err != nil {
					return err
				}

				return err
			},
		},
	}

	err := c.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
