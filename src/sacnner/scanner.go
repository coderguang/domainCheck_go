package domainCheckScanner

import (
	"bufio"
	domainCheckData "domainCheck/src/data"
	domainCheckDb "domainCheck/src/db"
	domainCheckMail "domainCheck/src/email"
	"io"
	"os"

	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sgwhois"

	"github.com/coderguang/GameEngine_go/sglog"
)

func ScanDomainByFile(filename string) {

	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		sglog.Fatal("read file:%s error,err:=%e", filename, err)
		return
	}

	domainList := []string{}
	sum := 0
	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		domain := string([]byte(line))
		domainList = append(domainList, domain)
		sum++
	}
	coroutinesNum := domainCheckData.GetCoroutineNum()
	everyNum := len(domainList) / coroutinesNum
	sglog.Info("readFile ok,sum=%d,every list len=%d", sum, everyNum)
	tmpDomainList := [][]string{}
	index := 0
	for i := 0; i < coroutinesNum; i++ {
		if i == coroutinesNum-1 {
			tmp := domainList[i*everyNum : len(domainList)]
			tmpDomainList = append(tmpDomainList, tmp)
			sglog.Info("index=%d list len=%d", index, len(tmp))
		} else {
			tmp := domainList[i*everyNum : (i+1)*everyNum]
			tmpDomainList = append(tmpDomainList, tmp)
			sglog.Info("index=%d list len=%d", index, len(tmp))
		}
		index++
	}
	chanList := [](chan bool){}
	index = 0
	for _, domainSet := range tmpDomainList {
		tmpChan := make(chan bool)
		go ScanDomainList(domainSet, tmpChan, index)
		chanList = append(chanList, tmpChan)
		index++
	}
	index = 0
	for _, n := range chanList {
		sglog.Info("wait for coroutines ,index=%d", index)
		<-n
		index++
	}

}

func ScanDomainList(domainList []string, flag chan bool, coroutinesIndex int) {
	for _, n := range domainList {
		ScanDomainAndSaveDb(n, coroutinesIndex)
		sgthread.SleepByMillSecond(200)
	}
	flag <- true
}

func ScanDomainAndSaveDb(domain string, coroutinesIndex int) {
	result, err := sgwhois.GetWhoisInfo(domain)
	if err != nil {
		domainCheckDb.InsertOrUpdateData(result)
		return
	}
	sglog.Info("start ScanDomainAndSaveDb coroutinesIndex=%s,domain=%s", coroutinesIndex, domain)
	if sgwhois.SG_WHOIS_STATUS_CHECK_FAILD == result.IsRegist {
		domainCheckDb.InsertOrUpdateData(result)
		return
	}

	sgwhois.ParseWhois(result)
	if sgwhois.SG_WHOIS_STATUS_CHECK_FAILD == result.IsRegist {
		domainCheckDb.InsertOrUpdateData(result)
		return
	}

	domainCheckDb.InsertOrUpdateData(result)

	HightValueNotice(result)
}

func CheckDomainExpireDt(zone string, name_length int) {

	sglog.Info("================start CheckDomainExpireDt,zone=%s,name_length=%d,load data from db", zone, name_length)

	domainList, err := domainCheckDb.GetAllExpireDtData(zone, name_length)

	if err != nil {
		sglog.Error("get domainlist from db error,err=%s", err)
		sgthread.SleepBySecond(5 * 60)
		CheckDomainExpireDt(zone, name_length)
		return
	}

	sglog.Info("================start CheckDomainExpireDt,zone=%s,name_length=%d,need check size=%d", zone, name_length, len(domainList))

	for _, v := range domainList {
		doUpdateDomainInfo(v)
		sgthread.SleepByMillSecond(200)
	}

	if 0 == len(domainList) {
		sleepTime, nextTime := domainCheckDb.GetLasteExpiryDt(zone, name_length)
		sglog.Info("zone:%s,name_length:%d,next will be run in %s", zone, name_length, nextTime.NormalString())
		sgthread.SleepBySecond(int(sleepTime))
	} else {
		sgthread.SleepBySecond(60)
	}
	CheckDomainExpireDt(zone, name_length)
}

func doUpdateDomainInfo(oldInfo *sgwhois.Whois) {
	result, err := sgwhois.GetWhoisInfo(oldInfo.Domain)

	if err != nil {
		return
	}
	if sgwhois.SG_WHOIS_STATUS_CHECK_FAILD == result.IsRegist {
		return
	}

	sgwhois.ParseWhois(result)
	if sgwhois.SG_WHOIS_STATUS_CHECK_FAILD == result.IsRegist {
		return
	}
	if sgwhois.SG_WHOIS_STATUS_CHECK_FAILD == result.IsRegist {
		sglog.Info("fail to check info,domain=%s", result.Domain)
		return
	}
	//sgwhois.ShowWhoisInfo(result)
	if result.IsRegist != sgwhois.SG_WHOIS_STATUS_CAN_REGIST_NOW {
		if result.IsEqual(oldInfo) {
			//sglog.Info("%s wouldn't update", result.Domain)
			return
		}
	}

	domainCheckDb.InsertOrUpdateData(result)
	HightValueNotice(result)

}

func HightValueNotice(result *sgwhois.Whois) {
	if sgwhois.IsHightValueDomainByName(result.Domain) {
		time_now := sgtime.New()
		if sgwhois.SG_WHOIS_STATUS_CAN_REGIST_NOW == result.IsRegist {
			domainCheckMail.SendMailNotice(result.Domain, "  can regist now")
			sglog.Info("luck domain ", result.Domain, " can regist now")
		} else if sgwhois.SG_WHOIS_STATUS_LIMIT_BY_GOVERNMENT != result.IsRegist && result.ExpiryDt.Before(time_now) {
			//SendMailNotice( domain, config, "  can regist "+result.ExpiryDtStr)
			sglog.Info("luck domain %s can regist at %s", result.Domain, result.ExpiryDtStr)
		}
	}
}
