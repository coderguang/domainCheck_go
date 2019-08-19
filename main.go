package main

import (
	domainCheckData "domainCheck/src/data"
	domainCheckDb "domainCheck/src/db"
	domainCheckMail "domainCheck/src/email"
	domainCheckGenTxt "domainCheck/src/genTxt"
	domainCheckScanner "domainCheck/src/sacnner"
	"log"
	"os"

	"github.com/coderguang/GameEngine_go/sgcmd"

	"github.com/coderguang/GameEngine_go/sgserver"

	_ "net/http/pprof"

	"github.com/coderguang/GameEngine_go/sglog"
)

func StartCheckByDb(cmd []string) {
	for i := 1; i < 6; i++ {
		go domainCheckScanner.CheckDomainExpireDt("com", i)
		go domainCheckScanner.CheckDomainExpireDt("cn", i)
		go domainCheckScanner.CheckDomainExpireDt("net", i)
	}
}

func StartGenDataAndCheck(cmd []string) {
	file := "./domainTxt/domain.txt"
	domainCheckGenTxt.CreateDominFile(file)
	domainCheckScanner.ScanDomainByFile(file)
	StartCheckByDb(cmd)
}

func registCmd() {
	sgcmd.RegistCmd("StartCheckByDb", "[\"StartCheckByDb\"] will start check by db data", StartCheckByDb)
	sgcmd.RegistCmd("StartGenDataAndCheck", "[\"StartGenDataAndCheck\"] will start gen domainData and check by file data", StartGenDataAndCheck)
}

func main() {

	sgserver.StartLogServer("debug", "./log", log.LstdFlags, true)

	arg_num := len(os.Args) - 1
	if arg_num < 1 {
		sglog.Error("please input config file ")
		return
	}

	sglog.Info("welcome to domain check !any question can ask royalchen@royalchen.com")

	err := domainCheckData.InitConfig(os.Args[1])

	if err != nil {
		os.Exit(1)
	}

	domainCheckDb.InitDbConnection()
	domainCheckMail.InitEmailConnect()

	registCmd()
	sgcmd.StartCmdWaitInputLoop()

	log.Println("end all scanner")
}
