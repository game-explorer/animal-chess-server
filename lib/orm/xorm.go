package orm

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type RWEngine struct {
	*xorm.Engine // 默认引擎
	Read         *xorm.Engine
	Write        *xorm.Engine
}

var debug = true

// 支持读写分离，需要输入读写两个数据源
func New(dsnRead, dsnWrite string) (engine *RWEngine, err error) {
	if dsnRead == "" || dsnWrite == "" {
		err = errors.New("invalid mysql dsn")
		return
	}
	engine = &RWEngine{}
	engine.Read, err = xorm.NewEngine("mysql", dsnRead)
	if err != nil {
		return
	}
	engine.Write, err = xorm.NewEngine("mysql", dsnWrite)
	if err != nil {
		return
	}

	engine.Read.ShowSQL(debug)
	engine.Read.ShowExecTime(debug)
	engine.Write.ShowSQL(debug)
	engine.Write.ShowExecTime(debug)

	// 默认引擎为Read
	engine.Engine = engine.Read
	return
}

func SetDebug(d bool) {
	debug = d
}
