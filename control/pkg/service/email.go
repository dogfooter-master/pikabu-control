package service

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

type mailConfigInfo struct {
	SMTPServer string `mapstructure:"smtp_server"`
	Port       string `mapstructure:"port"`
	ID         string `mapstructure:"id"`
	Password   string `mapstructure:"password"`
	From       string `mapstructure:"from"`
	Support    string `mapstructure:"support"`
}

func (c *mailConfigInfo) ReadConfig() (err error) {
	if err = LoadConfig(); err != nil {
		return
	}

	viper.GetStringMap("mail")
	viper.UnmarshalKey("mail", &c)

	return
}

type Email struct {
	SMTPServer          string
	Port                string
	ID                  string
	Password            string
	To                  string
	ToAlias             string
	ToMap               map[string]string
	ToName              string
	From                string
	FromName            string
	CC                  string
	CCAlias             string
	CCMap               map[string]string
	Subject             string
	Contents            string
	AttachmentFilePath  string
	Uri                 string
	CertificationNumber string
}

func (m *Email) SendEmail() (err error) {

	sendMail := gomail.NewMessage()
	sendMail.SetHeader("From", m.From)

	if len(m.ToMap) > 0 {
		for k, v := range m.ToMap {
			if len(v) > 0 {
				sendMail.SetAddressHeader("To", k, v)
			} else {
				sendMail.SetHeader("To", k)
			}
		}
	} else {
		if len(m.ToAlias) > 0 {
			sendMail.SetAddressHeader("To", m.To, m.ToAlias)
		} else {
			sendMail.SetHeader("To", m.To)
		}
	}

	if len(m.CCMap) > 0 {
		for k, v := range m.CCMap {
			if len(v) > 0 {
				sendMail.SetAddressHeader("Cc", k, v)
			} else {
				sendMail.SetHeader("Cc", k)
			}
		}
	} else {
		if len(m.CC) > 0 {
			if len(m.CCAlias) > 0 {
				sendMail.SetAddressHeader("Cc", m.CC, m.CCAlias)
			} else {
				sendMail.SetHeader("Cc", m.CC)
			}
		}
	}

	sendMail.SetHeader("Subject", m.Subject)
	//sendMail.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	sendMail.SetBody("text/html", m.Contents)
	if len(m.AttachmentFilePath) > 0 {
		sendMail.Attach(m.AttachmentFilePath)
	}

	port, _ := strconv.Atoi(m.Port)
	d := gomail.NewDialer(m.SMTPServer, port, m.ID, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(sendMail); err != nil {
		return fmt.Errorf("Send mail error: %v", err)
	}

	return
}
func (m *Email) SendEmailCertificationNumber() (err error) {

	var configInfo mailConfigInfo
	if err = configInfo.ReadConfig(); err != nil {
		err = fmt.Errorf("config error: %v", err)
		return
	}
	t := Template{
		Name:                m.ToAlias,
		Url:                 m.Uri,
		CertificationNumber: m.CertificationNumber,
	}
	

	m.SMTPServer = configInfo.SMTPServer
	m.Port = configInfo.Port
	m.ID = configInfo.ID
	m.Password = configInfo.Password
	m.From = configInfo.From
	m.Subject = "Verify Certification Code"
	m.Contents = t.Html()

	//mail := Email{
	//	SMTPServer: configInfo.SMTPServer,
	//	Port:       configInfo.Port,
	//	ID:         configInfo.ID,
	//	Password:   configInfo.Password,
	//	Subject:    "Reset your password",
	//	From:       configInfo.From,
	//	To:         to,
	//	Contents:   t.Html(),
	//}

	sendMail := gomail.NewMessage()
	sendMail.SetHeaders(map[string][]string{
		"From": {sendMail.FormatAddress(m.From, "The Dermaster Team")},
	})

	if len(m.ToMap) > 0 {
		for k, v := range m.ToMap {
			if len(v) > 0 {
				sendMail.SetAddressHeader("To", k, v)
			} else {
				sendMail.SetHeader("To", k)
			}
		}
	} else {
		if len(m.ToAlias) > 0 {
			sendMail.SetAddressHeader("To", m.To, m.ToAlias)
		} else {
			sendMail.SetHeader("To", m.To)
		}
	}

	sendMail.SetHeader("Subject", m.Subject)

	sendMail.SetBody("text/html", m.Contents)
	if len(m.AttachmentFilePath) > 0 {
		sendMail.Attach(m.AttachmentFilePath)
	}

	sendMail.Embed(os.Getenv("DERMASTER_HOME") + "/dermaster-control/img/dermaster-logo.png")

	port, _ := strconv.Atoi(m.Port)
	d := gomail.NewDialer(m.SMTPServer, port, m.ID, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(sendMail); err != nil {
		return fmt.Errorf("Send mail error: %v", err)
	}

	return
}
func (m *Email) SendEmailPasswordToken() (err error) {

	sendMail := gomail.NewMessage()
	sendMail.SetHeaders(map[string][]string{
		"From": {sendMail.FormatAddress(m.From, "The LazyBird Team")},
	})
	//sendMail.SetHeader("From", m.From)

	if len(m.ToMap) > 0 {
		for k, v := range m.ToMap {
			if len(v) > 0 {
				sendMail.SetAddressHeader("To", k, v)
			} else {
				sendMail.SetHeader("To", k)
			}
		}
	} else {
		if len(m.ToAlias) > 0 {
			sendMail.SetAddressHeader("To", m.To, m.ToAlias)
		} else {
			sendMail.SetHeader("To", m.To)
		}
	}

	sendMail.SetHeader("Subject", m.Subject)
	//sendMail.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	sendMail.SetBody("text/html", m.Contents)
	if len(m.AttachmentFilePath) > 0 {
		sendMail.Attach(m.AttachmentFilePath)
	}

	sendMail.Embed(os.Getenv("WORKDOMO_HOME") + "/srv-email" + "/img/workdomo-logo@2x.png")

	port, _ := strconv.Atoi(m.Port)
	d := gomail.NewDialer(m.SMTPServer, port, m.ID, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(sendMail); err != nil {
		return fmt.Errorf("Send mail error: %v", err)
	}

	return
}
func (m *Email) SendEmailBeta() (err error) {
	sendMail := gomail.NewMessage()
	sendMail.SetHeaders(map[string][]string{
		"From": {sendMail.FormatAddress(m.From, "The LazyBird Team")},
	})

	if len(m.ToMap) > 0 {
		for k, v := range m.ToMap {
			if len(v) > 0 {
				sendMail.SetAddressHeader("To", k, v)
			} else {
				sendMail.SetHeader("To", k)
			}
		}
	} else {
		if len(m.ToAlias) > 0 {
			sendMail.SetAddressHeader("To", m.To, m.ToAlias)
		} else {
			sendMail.SetHeader("To", m.To)
		}
	}

	sendMail.SetHeader("Subject", m.Subject)
	//sendMail.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	sendMail.SetBody("text/html", m.Contents)
	if len(m.AttachmentFilePath) > 0 {
		sendMail.Attach(m.AttachmentFilePath)
	}

	sendMail.Embed(os.Getenv("WORKDOMO_HOME") + "/srv-email" + "/img/workdomo-logo@2x.png")

	port, _ := strconv.Atoi(m.Port)
	d := gomail.NewDialer(m.SMTPServer, port, m.ID, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(sendMail); err != nil {
		return fmt.Errorf("Send mail error: %v", err)
	}

	return
}
func (m *Email) SendEmailContactUs() (err error) {

	sendMail := gomail.NewMessage()
	sendMail.SetHeader("From", m.From)

	if len(m.ToMap) > 0 {
		for k, v := range m.ToMap {
			if len(v) > 0 {
				sendMail.SetAddressHeader("To", k, v)
			} else {
				sendMail.SetHeader("To", k)
			}
		}
	} else {
		if len(m.ToAlias) > 0 {
			sendMail.SetAddressHeader("To", m.To, m.ToAlias)
		} else {
			sendMail.SetHeader("To", m.To)
		}
	}

	sendMail.SetHeader("Subject", m.Subject)
	//sendMail.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	sendMail.SetBody("text/html", m.Contents)
	if len(m.AttachmentFilePath) > 0 {
		sendMail.Attach(m.AttachmentFilePath)
	}

	port, _ := strconv.Atoi(m.Port)
	d := gomail.NewDialer(m.SMTPServer, port, m.ID, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(sendMail); err != nil {
		return fmt.Errorf("Send mail error: %v", err)
	}

	return
}
func (m *Email) SendEmailJoinUs() (err error) {

	sendMail := gomail.NewMessage()
	sendMail.SetHeaders(map[string][]string{
		"From": {sendMail.FormatAddress(m.From, m.FromName)},
	})
	sendMail.SetHeaders(map[string][]string{
		"To": {sendMail.FormatAddress(m.To, m.ToName)},
	})

	sendMail.SetHeader("Subject", m.Subject)

	sendMail.SetBody("text/html", m.Contents)
	sendMail.Embed(os.Getenv("WORKDOMO_HOME") + "/srv-email" + "/img/workdomo-logo.png")

	port, _ := strconv.Atoi(m.Port)
	d := gomail.NewDialer(m.SMTPServer, port, m.ID, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(sendMail); err != nil {
		return fmt.Errorf("Send mail error: %v", err)
	}

	return
}
