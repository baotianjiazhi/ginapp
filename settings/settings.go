package settings

import (
	"log"

	"github.com/go-ini/ini"
)

type DataBase struct {
	Host      string
	Port      string
	User      string
	Password  string
	DbName    string
	CharSet   string
	ParseTime string
	MaxIdle   int
	MaxOpen   int
}

type App struct {
	Name string
	Mode string
	Port string
}

type Log struct {
	Level      string
	FileName   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type Redis struct {
	Host     string
	Port     string
	DB       int
	PoolSize int
}

var DataBaseSetting = &DataBase{}
var AppSetting = &App{}
var LogSetting = &Log{}
var RedisSetting = &Redis{}

var cfg *ini.File

func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf.ini':%v", err)
	}
	mapTo("mysql", DataBaseSetting)
	mapTo("app", AppSetting)
	mapTo("log", LogSetting)
	mapTo("redis", RedisSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
