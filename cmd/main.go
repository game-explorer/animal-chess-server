package main

import (
	"github.com/game-explorer/animal-chess-server/app/config"
	"github.com/game-explorer/animal-chess-server/lib/signal"
	"github.com/game-explorer/animal-chess-server/service/gin"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
	"time"
)

func main() {
	c := cli.NewApp()
	c.Name = "animal-chess"
	c.Usage = ""
	c.Version = ""
	c.Action = func(*cli.Context) error {
		ctx, cancel := signal.NewTermContext()
		if config.IsCi {
			// 如果是测试, 则等待1s关闭程序
			go func() {
				time.Sleep(1 * time.Second)
				cancel()
			}()
		}

		if !config.IsCi {
			err := repository.InitMysql()
			if err != nil {
				return err
			}
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			s := gin.New()
			log.Infof("http listen '%s'", config.App.HttpAddr)
			s.Listen(config.App.HttpAddr, ctx)
			log.Info("http shutdown")
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			log.Infof("worker running")
			err := worker.RunWhenPublish(ctx)
			if err != nil {
				panic(err)
			}
			log.Info("worker shutdown")
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
