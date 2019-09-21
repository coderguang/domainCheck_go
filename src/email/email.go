package domainCheckMail

import (
	"crypto/tls"
	domainCheckData "domainCheck/src/data"
	"net"
	"net/smtp"
	"strings"

	"github.com/coderguang/GameEngine_go/sglog"
)

var globalMailAuth smtp.Auth

func InitEmailConnect() {
	mailFrom, mailPwd, _, smtpUrl, _ := domainCheckData.GetMailConnectionInfo()
	if "" == mailFrom {
		sglog.Info("email would not connection")
		return
	}

	globalMailAuth = smtp.PlainAuth("", mailFrom, mailPwd, smtpUrl)
}

func SendMailNotice(domain string, exMsg string) {
	if "" == domain {
		sglog.Error("null domain ,unknow error")
		return
	}

	sendMailUsingTLS(domain, exMsg)
}

//return a smtp client
func dail(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		sglog.Error("Dialing Error:%s", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func sendMailUsingTLS(domain string, exMsg string) error {
	mailFrom, _, mailTo, smtpUrl, stmpPort := domainCheckData.GetMailConnectionInfo()

	if "" == mailFrom {
		sglog.Error("email would not send by because mail config not set")
		return nil
	}
	to := []string{mailTo}

	content_type := "Content-Type: text/plain; charset=UTF-8"
	subject := "Domain check Notice:" + domain + "  " + exMsg
	content := "luck check a unregiest domain:" + domain
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + "aliyun server" +
		"<" + mailFrom + ">\r\nSubject:" + subject + "\r\n" + content_type + "\r\n\r\n" + content)
	addr := smtpUrl + ":" + stmpPort

	//参考net/smtp的func SendMail()
	//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
	//len(to)>1时,to[1]开始提示是密送

	c, err := dail(addr)
	if err != nil {
		sglog.Error("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if globalMailAuth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(globalMailAuth); err != nil {
				sglog.Error("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(mailFrom); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
