package noaweb

import (
	"errors"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type mailgunConfigs map[string]mailgun.Mailgun

// MailgunFunctions struct
type MailgunFunctions struct{}

// Mailgun variable to be used when calling mailgun functions
var Mailgun MailgunFunctions

// InitMailgun initiates a mailgun instance associated with the noaweb instance.
// Name of mailgun instance will be domain name.
func (MailgunFunctions) Init(domain, apiKey, publicAPIkey string) error {
	noawebinst.mailguns[domain] = mailgun.NewMailgun(domain, apiKey, publicAPIkey)
	return nil
}

// Send forms a mailgun message and sends it.
func (MailgunFunctions) Send(domain, from, subject, text, to string) error {
	if val, ok := noawebinst.mailguns[domain]; ok {
		message := mailgun.NewMessage(
			from,
			subject,
			text,
			to)
		_, _, err := val.Send(message)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Noaweb.SendMailgun: Could not find mailgun config for given domain.")
}
