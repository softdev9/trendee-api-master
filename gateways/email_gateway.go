package gateways

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/mailgun/mailgun-go"
	"github.com/softdev9/trendee-api-master/data"
	"os"
)

// Provide all the method needed to interact with MailGun API
const (
	domain = "trendee.co"
	from   = "webmaster@" + domain
)

type TrendeeMailSender interface {
	// Send a forgot password to the user given as parameter
	SendForgotPassword(u data.User) (string, error)
	SendWelcomeMessage(u data.User) (string, error)
	EmailForEmailChange(email string, updateLink string) (string, error)
}

type MailGunClient struct {
	MgClient mailgun.Mailgun
}

func InitMailGun() MailGunClient {
	apiKey := os.Getenv("MAIL_GUN_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")
	return MailGunClient{MgClient: mg}
}

func (client MailGunClient) EmailForEmailChange(email string, updateLink string) (string, error) {
	return client.sendMessage(email, "Changement email trendee",
		"Bonjour, voici le lien cliquable pour confirmer votre changement d'email "+updateLink)
}

func (client MailGunClient) sendMessage(to string, object string, msg string) (string, error) {
	m := client.MgClient.NewMessage(from, object, msg, to)
	_, id, err := client.MgClient.Send(m)
	return id, err
}

func (c MailGunClient) SendForgotPassword(u data.User) (string, error) {
	return c.sendMessage(u.Email, "Mot de passe oublie",
		"Voici le lien pour active reinitialiser votre mot de passe")
}

func (c MailGunClient) SendWelcomeMessage(u data.User) (string, error) {
	return c.sendMessage(u.Email, "Bienvenu dans trendee",
		"Bonjour, Bienvenue dans trendee")
}
