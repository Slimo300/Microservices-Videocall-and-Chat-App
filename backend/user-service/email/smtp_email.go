package email

import (
	"bytes"
	"crypto/tls"
	"log"
	"text/template"

	"github.com/k3a/html2text"
	gomail "gopkg.in/gomail.v2"
)

type SMTPEmailService struct {
	EmailFrom string
	SMTPHost  string
	SMTPPass  string
	SMTPPort  int
	SMTPUser  string
	templates *template.Template
}

// 👇 Email template parser
func (srv SMTPEmailService) SendVerificationEmail(data VerificationEmailData) error {

	var body bytes.Buffer

	if err := srv.templates.ExecuteTemplate(&body, "../templates/verification.html", data); err != nil {
		log.Fatal("Could not execute template", err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", srv.EmailFrom)
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(srv.SMTPHost, srv.SMTPPort, srv.SMTPUser, srv.SMTPPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
