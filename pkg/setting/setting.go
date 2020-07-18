package setting

import (
	"log"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string

	Peers  []string
	BcConf string
)

func init() {

	var err error
	var err1 error

	Cfg, err = ini.Load("config/app.ini")
	if err != nil {
		//如果是test修改测试路径
		Cfg, err1 = ini.Load("../config/app.ini")
		if err1 != nil {
			log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
		}
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadBlockchain()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func LoadBlockchain() {
	sec, err := Cfg.GetSection("blockchain")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	p := sec.Key("PEERS").MustString("peer0.org1.lzawt.com")
	Peers = strings.Split(p, ",")
	BcConf = sec.Key("CONFIG").MustString("./config.yaml")
}
