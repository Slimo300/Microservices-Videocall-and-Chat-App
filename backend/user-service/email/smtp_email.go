package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/k3a/html2text"
	gomail "gopkg.in/gomail.v2"
)

type SMTPEmailService struct {
	SMTPDialer *gomail.Dialer
	Templates  map[string]*template.Template
	EmailFrom  string
	Origin     string
}

func NewSMTPService(emailDir, emailFrom, host string, port int, user, pass string) (*SMTPEmailService, error) {

	templateCache := make(map[string]*template.Template)

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", emailDir))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseGlob(fmt.Sprintf("%s/*.layout.html", emailDir))
		if err != nil {
			return nil, err
		}

		templateCache[name] = tmpl
	}

	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &SMTPEmailService{
		SMTPDialer: d,
		Templates:  templateCache,
		EmailFrom:  emailFrom,
	}, nil
}

// Sending email
func (srv SMTPEmailService) SendEmail(tmpl string, data EmailData) error {

	var body bytes.Buffer

	t, ok := srv.Templates[tmpl]
	if !ok {
		return errors.New("Could not get template")
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", srv.EmailFrom)
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	if err := srv.SMTPDialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
