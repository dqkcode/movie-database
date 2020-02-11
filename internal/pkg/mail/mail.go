package mail

import (
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dqkcode/movie-database/internal/pkg/config/envconfig"
	"github.com/jordan-wright/email"
)

type (
	Config struct {
		Host     string `envconfig:"SMTP_HOST"`
		Port     string `envconfig:"SMTP_PORT"`
		Count    int    `envconfig:"SMTP_COUNT"`
		UserName string `envconfig:"SMTP_USERNAME"`
		Password string `envconfig:"SMTP_PASSWORD"`
	}
	Email struct {
		From        string
		To          []string
		Bcc         []string
		Cc          []string
		Subject     string
		Text        []byte // Plaintext message (optional)
		HTML        []byte // Html message (optional)
		Headers     textproto.MIMEHeader
		ReadReceipt []string
	}
	Mailer struct {
		pool *email.Pool
	}
)

func LoadConfigFromEnv() Config {
	var conf Config
	envconfig.Load(&conf)
	return conf
}

func NewPool() (*Mailer, error) {
	conf := LoadConfigFromEnv()
	addr := conf.Host + conf.Port
	p, err := email.NewPool(addr, conf.Count, smtp.PlainAuth("", conf.UserName, conf.Password, conf.Host))
	if err != nil {
		logrus.Errorf("err create pool email: %v", err)
		return nil, err
	}
	return &Mailer{
		pool: p,
	}, nil
}

func (m *Mailer) Send(e Email, timeOut time.Duration) {
	em := &email.Email{
		From:    e.From,
		To:      e.To,
		Cc:      e.Cc,
		Bcc:     e.Bcc,
		Headers: e.Headers,
		Subject: e.Subject,
		Text:    e.Text,
		HTML:    e.HTML,
	}
	m.pool.Send(em, timeOut)
}
