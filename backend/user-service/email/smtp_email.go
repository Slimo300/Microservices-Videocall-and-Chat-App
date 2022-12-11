package email

import (
	"bytes"
	"crypto/tls"
	"log"
	"path/filepath"
	"text/template"

	"github.com/k3a/html2text"
	gomail "gopkg.in/gomail.v2"
)

type SMTPEmailService struct {
	EmailFrom             string
	SMTPHost              string
	SMTPPass              string
	SMTPPort              int
	SMTPUser              string
	VerificationTemplate  *template.Template
	ResetPasswordTemplate *template.Template
}

func NewSMTPService(emailDir, emailFrom, host string, port int, user, pass string) (*SMTPEmailService, error) {

	verificationTemplates, err := template.ParseFiles(
		filepath.Join(emailDir, "verification.html"),
		filepath.Join(emailDir, "base.html"),
		filepath.Join(emailDir, "styles.html"),
	)
	if err != nil {
		return nil, err
	}

	resetPasswordTemplates, err := template.ParseFiles(
		filepath.Join(emailDir, "reset.html"),
		filepath.Join(emailDir, "base.html"),
		filepath.Join(emailDir, "styles.html"),
	)
	if err != nil {
		return nil, err
	}

	return &SMTPEmailService{
		EmailFrom:             emailFrom,
		SMTPHost:              host,
		SMTPPass:              pass,
		SMTPPort:              port,
		SMTPUser:              user,
		VerificationTemplate:  verificationTemplates,
		ResetPasswordTemplate: resetPasswordTemplates,
	}, nil
}

// 👇 Email template parser
func (srv SMTPEmailService) SendVerificationEmail(data EmailData) error {

	var body bytes.Buffer

	if err := srv.VerificationTemplate.ExecuteTemplate(&body, "verification.html", data); err != nil {
		log.Fatal("Could not execute template", err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", srv.EmailFrom)
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", data.Subject)
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

// 👇 Email template parser
func (srv SMTPEmailService) SendResetPasswordEmail(data EmailData) error {

	var body bytes.Buffer

	if err := srv.ResetPasswordTemplate.ExecuteTemplate(&body, "reset.html", data); err != nil {
		log.Fatal("Could not execute template", err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", srv.EmailFrom)
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", data.Subject)
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
