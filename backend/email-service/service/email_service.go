package service

import (
	"bytes"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/templates"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

// Email Service implements EmailServiceServer interface for grpc connections
type EmailService struct {
	SMTPDialer *gomail.Dialer
	Templates  map[string]*template.Template
	EmailFrom  string
	Origin     string
}

func NewEmailService(emailFrom, host string, port int, user, pass, origin string) (EmailService, error) {
	templateCache := make(map[string]*template.Template)
	pages, err := fs.Glob(templates.FS, "*.page.html")
	if err != nil {
		return EmailService{}, err
	}

	for _, page := range pages {
		name := strings.Split(filepath.Base(page), ".")[0]
		tmpl, err := template.ParseFS(templates.FS, page, "*.layout.html")
		if err != nil {
			return EmailService{}, err
		}
		templateCache[name] = tmpl
	}

	return EmailService{
		SMTPDialer: gomail.NewDialer(host, port, user, pass),
		Templates:  templateCache,
		EmailFrom:  emailFrom,
		Origin:     origin,
	}, nil
}

func (srv EmailService) SendVerificationEmail(email, username, code string) error {
	return srv.SendEmail("verification", EmailData{
		Subject: "Verification Email",
		Email:   email,
		Name:    username,
		Code:    code,
		Origin:  srv.Origin,
	})
}

func (srv EmailService) SendResetPasswordEmail(email, username, code string) error {
	return srv.SendEmail("reset", EmailData{
		Subject: "Reset Password",
		Email:   email,
		Name:    username,
		Code:    code,
		Origin:  srv.Origin,
	})
}

type EmailData struct {
	Subject string
	Email   string
	Name    string
	Code    string
	Origin  string
}

func (srv EmailService) SendEmail(tmpl string, data EmailData) error {
	var body bytes.Buffer
	t, ok := srv.Templates[tmpl]
	if !ok {
		return errors.New("could not get template")
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
