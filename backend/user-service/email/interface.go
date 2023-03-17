package email

type EmailService interface {
	SendEmail(tmpl string, data EmailData) error
}

type EmailData struct {
	Subject string
	Email   string
	Name    string
	Code    string
	Origin  string
}
