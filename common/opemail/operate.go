package opemail

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// Client 邮件处理客户端
type Client interface {
	SendEmail(subject, from, sender string, sendBody *strings.Builder, ctype string, email ...string) error
	SendEmailWithAttach(subject, from, sender string, sendBody *strings.Builder, ctype, fileName string, content []byte, email ...string) error
}

// GetRandEmailConf 获取随机的邮箱配置
type GetRandEmailConf func() *EmailConfig

// EmailConfig 邮箱客户端配置
type EmailConfig struct {
	Port     int    `copier:"SMTPPort"`
	Host     string `copier:"SMTPHost"`
	User     string `copier:"Username"`
	Password string `copier:"Password"`
}

// NewEmailConfig 新建邮箱客户端配置
func NewEmailConfig(host, user, pwd string, port int) *EmailConfig {
	return &EmailConfig{
		Host:     host,
		User:     user,
		Password: pwd,
		Port:     port,
	}
}

// NewEmailConfigs 从多个邮箱配置中获取随机的邮箱配置
func NewEmailConfigs(confs []EmailConfig) GetRandEmailConf {
	return func() *EmailConfig {
		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)

		return &confs[r.Intn(len(confs))]
	}
}

// SendEmail 发送邮件
func (c *EmailConfig) SendEmail(subject, from, sender string, sendBody *strings.Builder, ctype string, email ...string) error {
	m := gomail.NewMessage()
	if from != "" {
		m.SetHeader("From", from)
	}
	if sender != "" {
		m.SetHeader("Sender", sender)
	}
	m.SetHeader("To", email...)
	m.SetHeader("Subject", subject)
	m.SetBody(ctype, sendBody.String())

	d := gomail.NewDialer(c.Host, c.Port, c.User, c.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send %s email error: %w", subject, err)
	}

	return nil
}

// SendEmailWithAttach 发送带附件的邮件
func (c *EmailConfig) SendEmailWithAttach(subject, from, sender string, sendBody *strings.Builder, ctype, fileName string, content []byte, email ...string) error {
	m := gomail.NewMessage()
	if from != "" {
		m.SetHeader("From", from)
	}
	if sender != "" {
		m.SetHeader("Sender", sender)
	}
	m.SetHeader("To", email...)
	m.SetHeader("Subject", subject)
	m.SetBody(ctype, sendBody.String())
	m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(content)
		return err
	}))

	d := gomail.NewDialer(c.Host, c.Port, c.User, c.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send %s email error: %w", subject, err)
	}

	return nil
}
