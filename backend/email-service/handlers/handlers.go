package handlers

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/email/pb"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

// Email Service implements EmailServiceServer interface for grpc connections
type EmailService struct {
	*pb.UnimplementedEmailServiceServer
	SMTPDialer *gomail.Dialer
	Templates  map[string]*template.Template
	EmailFrom  string
	Origin     string
}

// NewEmailService is a constructor for EmailService type
func NewEmailService(emailDir, emailFrom, host string, port int, user, pass, origin string) (*EmailService, error) {
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

	return &EmailService{
		SMTPDialer: d,
		Templates:  templateCache,
		EmailFrom:  emailFrom,
		Origin:     origin,
	}, nil
}

// SendVerificationEmail sends email with specified data
func (srv EmailService) SendVerificationEmail(ctx context.Context, data *pb.EmailData) (*pb.Msg, error) {
	return &pb.Msg{}, srv.SendEmail("verification.page.html", EmailData{
		Subject: "Verification Email",
		Email:   data.Email,
		Name:    data.Name,
		Code:    data.Code,
		Origin:  srv.Origin,
	})
}

// SendResetPasswordEmail sends email with specified data
func (srv EmailService) SendResetPasswordEmail(ctx context.Context, data *pb.EmailData) (*pb.Msg, error) {
	return &pb.Msg{}, srv.SendEmail("reset.page.html", EmailData{
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
