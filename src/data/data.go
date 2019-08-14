package domainCheckData

import (
	domainCheckDef "domainCheck/src/define"
	"encoding/json"
	"io/ioutil"

	"github.com/coderguang/GameEngine_go/sglog"
)

var globalCfg *domainCheckDef.DbConfig

func InitConfig(configfile string) error {
	config, err := ioutil.ReadFile(configfile)
	if err != nil {
		sglog.Fatal("read config error")
		return err
	}
	globalCfg = new(domainCheckDef.DbConfig)
	p := &globalCfg
	err = json.Unmarshal([]byte(config), p)
	if err != nil {
		sglog.Fatal("parse config error")
		return err
	}
	return nil
}

func GetDbConnectionInfo() (string, string, string, string, string) {
	return globalCfg.Dbuser, globalCfg.Dbpwd, globalCfg.Dburl, globalCfg.Dbport, globalCfg.Dbname
}

func GetMailConnectionInfo() (string, string, string, string, string) {
	return "", globalCfg.EmailFrom, globalCfg.EmailPwd, globalCfg.EmailTo, globalCfg.Smtp
}

func GetCoroutineNum() int {
	return globalCfg.DoroutinesNum
}
