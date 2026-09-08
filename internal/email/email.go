package email

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/fkdiy/webalert/internal/config"
	"github.com/fkdiy/webalert/internal/monitor"
)

func Send(conf config.EMailConfig, change monitor.Change) error {
	message := fmt.Sprintf(
		"From: Web alert <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: Webalert detected a change at %s\r\n"+
			"\r\n"+
			"Change registered at: %s\r\n"+
			"\r\n"+
			"Selector: %s\r\n"+
			"\r\n"+
			"Current: %s\r\n"+
			"Previous: %s\r\n",
		conf.From,
		strings.Join(conf.Recipients, ", "),
		change.Target.URL,
		change.Target.URL,
		change.Target.Selector,
		change.Current,
		change.Previous,
	)

	return send(conf, []byte(message))
}

func SendDigest(conf config.EMailConfig, changes []monitor.Change) error {
	message := fmt.Sprintf(
		"From: Web alert <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: Webalert detected %d changes\r\n"+
			"\r\n",
		conf.From,
		strings.Join(conf.Recipients, ", "),
		len(changes),
	)

	for _, change := range changes {
		message = message + fmt.Sprintf(
			"Change registered at: %s\r\n"+
				"\r\n"+
				"Selector: %s\r\n"+
				"\r\n"+
				"Current: %s\r\n"+
				"Previous: %s\r\n"+
				"\r\n"+
				"---\r\n"+
				"\r\n",
			change.Target.URL,
			change.Target.Selector,
			change.Current,
			change.Previous,
		)
	}

	return send(conf, []byte(strings.Trim(message, "---\r\n\r\n")))
}

func send(conf config.EMailConfig, message []byte) error {
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

	return smtp.SendMail(
		addr,
		auth,
		conf.From,
		conf.Recipients,
		message,
	)
}
