package email

import (
	"bytes"
	"crypto/tls"
	"log"
	"os"
	"path/filepath"
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

func NewSMTPService(emailDir, emailFrom, host string, port int, user, pass string) (*SMTPEmailService, error) {
	var paths []string

	if err := filepath.Walk(emailDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	templates, err := template.ParseFiles(paths...)
	if err != nil {
		return nil, err
	}

	return &SMTPEmailService{
		EmailFrom: emailFrom,
		SMTPHost:  host,
		SMTPPort:  port,
		SMTPUser:  user,
		SMTPPass:  pass,
		templates: templates,
	}, nil
}

// 👇 Email template parser
func (srv SMTPEmailService) SendVerificationEmail(data EmailData) error {

	var body bytes.Buffer

	if err := srv.templates.ExecuteTemplate(&body, "verification.html", data); err != nil {
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
