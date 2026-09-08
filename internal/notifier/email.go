package email

import (
	"net"
	"net/smtp"
	"strings"

	"github.com/fkdiy/webalert/internal/config"
	"github.com/fkdiy/webalert/internal/monitor"
)

func Send(conf config.EMailConfig, change monitor.Change) error {
	auth := smtp.PlainAuth(
		"",
		conf.SMTP.Username,
		conf.SMTP.Password,
		conf.SMTP.Host,
	)

	addr := net.JoinHostPort(
		conf.SMTP.Host,
		conf.SMTP.Port,
	)

	message := []byte(
		"From: Web alert <" + conf.From + ">\r\n" +
			"To: " + strings.Join(conf.Recipients, ", ") + "\r\n" +
			"Subject: Webalert detected a change at " + change.Target.URL + "\r\n" +
			"\r\n" +
			"Change registered at: " + change.Target.URL + "\r\n\n" +
			"Selector: " + change.Target.Selector + "\r\n\n" +
			"Current: " + change.Current + "\r\n" +
			"Previous: " + change.Previous + "\r\n",
	)

	return smtp.SendMail(
		addr,
		auth,
		conf.From,
		conf.Recipients,
		message,
	)
}
