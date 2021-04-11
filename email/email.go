package email

import (
	"github.com/go-gomail/gomail"
)

const (
	fromName    = "From"
	sendTo      = "To"
	sendSubject = "subject"
	sendBody    = "text/html"

	defServerName = "sta-golang"
	ServerMsg     = "STA音乐推荐系统消息 Level[%s]"
)

type EmailConfig struct {
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Email       string `json:"email" yaml:"email"`
	Pwd         string `json:"pwd" yaml:"pwd"`
	ServerName  string `json:"serverName" yaml:"serverName"`
	ContentType string `json:"contentType" yaml:"contentType"`
}

type emailClient struct {
	cfg    *EmailConfig
	helper *gomail.Dialer
}

var NewEmailClient = func(cfg *EmailConfig) *emailClient {
	if cfg.ServerName == "" {
		cfg.ServerName = defServerName
	} else if cfg.ContentType == "" {
		cfg.ContentType = sendBody
	}
	return &emailClient{
		cfg:    cfg,
		helper: gomail.NewDialer(cfg.Host, cfg.Port, cfg.Email, cfg.Pwd),
	}
}

func (es *emailClient) newMessage(subject, body string, mailTo ...string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader(fromName, es.cfg.Email, es.cfg.ServerName)
	message.SetHeader(sendTo, mailTo...)
	message.SetHeader(sendSubject, subject)
	message.SetBody(sendBody, body)
	return message
}

func (es *emailClient) SendEmail(subject, body string, mailTo ...string) error {
	err := es.helper.DialAndSend(es.newMessage(subject, body, mailTo...))
	if err != nil {
		return err
	}
	return nil
}
