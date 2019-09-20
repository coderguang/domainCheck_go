package domainCheckMail

import (
	domainCheckData "domainCheck/src/data"
	"net/smtp"

	"github.com/coderguang/GameEngine_go/sglog"

	"github.com/coderguang/GameEngine_go/sgmail"
)

var globalMailAuth *smtp.Auth

func InitEmailConnect() {
	mailFrom, mailPwd, _, smtpUrl, _ := domainCheckData.GetMailConnectionInfo()
	if "" == mailFrom {
		sglog.Info("email would not connection")
		return
	}

	globalMailAuth = sgmail.PlainAuth("", mailFrom, mailPwd, smtpUrl)
}

func SendMailNotice(domain string, exMsg string) {
	if "" == domain {
		sglog.Error("null domain ,unknow error")
		return
	}

	mailFrom, _, mailTo, smtpUrl, stmpPort := domainCheckData.GetMailConnectionInfo()

	if "" == mailFrom {
		sglog.Error("email would not send by because mail config not set")
		return
	}

	to := []string{mailTo}
	err := sgmail.SendMail(smtpUrl+":"+stmpPort, globalMailAuth, to, "Domain check Notice:"+domain+"  "+exMsg, "static server", mailFrom, "luck check a unregiest domain:"+domain)
	if err != nil {
		sglog.Error("sg_email.SendMail error,domain=%s,error=%s", domain, err)
	}
}
