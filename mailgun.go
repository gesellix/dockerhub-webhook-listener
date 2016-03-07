package listener

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/sirsean/go-mailgun/mailgun"
)

type Mailgun struct {
	mailGunConfig
}

type mailGunConfig struct {
	From   string
	To     []string
	Name   string
	Key    string
	Domain string
}

func (m *Mailgun) Call(hubMsg HubMessage) error {
	if len(m.To) == 0 {
		log.Print("MailGun: no recipients configured. Nothing to do.")
		return errors.New("MailGun: no recipients configured. Nothing to do.")
	}
	msg := mailgun.Message{
		FromName:      m.Name,
		FromAddress:   m.From,
		Subject:       "Some Subject",
		ToAddress:     m.To[0],
		CCAddressList: m.To[1:],
	}

	body, err := json.Marshal(hubMsg)
	if err != nil {
		log.Print(err)
		return err
	}
	msg.Body = string(body)
	client := mailgun.NewClient(m.Key, m.Domain)
	_, err = client.Send(msg)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
