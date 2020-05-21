package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type opt struct {
	fileName string // 要向上查找的config文件名字, 可以带路径, 如/config/config.yaml
}

type Option func(*opt)

func WithFileName(name string) Option {
	return func(o *opt) {
		o.fileName = name
	}
}

// config: 结构体指针
func Init(config interface{}, ops ...Option) {
	p := opt{
		fileName: "config.yml",
	}

	for _, v := range ops {
		v(&p)
	}

	configFileName := p.fileName
	var findConfig bool
	// 向上层查找配置文件
	// 在项目的任何地方运行(test时)都能加载到配置文件
	// 优先使用最近的配置文件
	for i := 0; i < 10; i++ {
		_, err := os.Stat(configFileName)
		if err != nil {
			if os.IsNotExist(err) {
				configFileName = "../" + configFileName
			} else {
				panic(err)
			}
		} else {
			findConfig = true
			break
		}
	}

	if !findConfig {
		panic(fmt.Sprintf("can't find config, Please write the config file in %s", p.fileName))
	}

	configFileName, err := filepath.Abs(configFileName)
	if err != nil {
		panic(err)
	}

	log.Printf("found config in `%s`", configFileName)

	bs, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Panicf("can't read '%s'", configFileName)
		return
	}

	err = yaml.Unmarshal(bs, config)
	if err != nil {
		log.Panicf("yaml.Unmarshal err:%v; row:%s", err, bs)
		return
	}
}
