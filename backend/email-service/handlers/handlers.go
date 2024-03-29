package handlers

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/email-service/templates"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/email"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

// Email Service implements EmailServiceServer interface for grpc connections
type EmailService struct {
	*email.UnimplementedEmailServiceServer
	SMTPDialer *gomail.Dialer
	Templates  map[string]*template.Template
	EmailFrom  string
	Origin     string
}

// NewEmailService is a constructor for EmailService type
func NewEmailService(emailFrom, host string, port int, user, pass, origin string) (*EmailService, error) {

	templateCache := make(map[string]*template.Template)

	pages, err := fs.Glob(templates.FS, "*.page.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := strings.Split(filepath.Base(page), ".")[0]

		tmpl, err := template.ParseFS(templates.FS, page, "*.layout.html")
		if err != nil {
			return nil, err
		}

		templateCache[name] = tmpl
	}

	dialer := gomail.NewDialer(host, port, user, pass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailService{
		SMTPDialer: dialer,
		Templates:  templateCache,
		EmailFrom:  emailFrom,
		Origin:     origin,
	}, nil
}

// SendVerificationEmail sends email with specified data
func (srv EmailService) SendVerificationEmail(ctx context.Context, data *email.EmailData) (*email.Msg, error) {
	return &email.Msg{}, srv.SendEmail("verification", EmailData{
		Subject: "Verification Email",
		Email:   data.Email,
		Name:    data.Name,
		Code:    data.Code,
		Origin:  srv.Origin,
	})
}

// SendResetPasswordEmail sends email with specified data
func (srv EmailService) SendResetPasswordEmail(ctx context.Context, data *email.EmailData) (*email.Msg, error) {
	return &email.Msg{}, srv.SendEmail("reset", EmailData{
		Subject: "Reset Password",
		Email:   data.Email,
		Name:    data.Name,
		Code:    data.Code,
		Origin:  srv.Origin,
	})
}

// EmailData is data passed to templates
type EmailData struct {
	Subject string
	Email   string
	Name    string
	Code    string
	Origin  string
}

// SendEmail is a helper function that wraps setting up email for sending
func (srv EmailService) SendEmail(tmpl string, data EmailData) error {

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
