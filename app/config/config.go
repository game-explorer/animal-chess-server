package config

import (
	"github.com/game-explorer/animal-chess-server/internal/pkg/config"
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"github.com/game-explorer/animal-chess-server/internal/pkg/orm"
)

var App struct {
	// 业务Debug
	Debug bool `yaml:"debug"`
	// OrmDebug开启后会打印sql语句
	OrmDebug bool `yaml:"orm_debug"`
	// LogDebug开启后会使用颜色
	LogDebug bool `yaml:"log_debug"`

	HttpAddr string `yaml:"http_addr"`

	Mysql struct {
		AnimalChess string `yaml:"animal_chess"`
	} `yaml:"mysql"`
}

func init() {
	config.Init(&App, config.WithFileName("config/config.yml"))
	log.SetDebug(App.LogDebug)
	orm.SetDebug(App.OrmDebug)
}
