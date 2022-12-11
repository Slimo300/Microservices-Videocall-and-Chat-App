package email

type EmailService interface {
	SendVerificationEmail(data EmailData) error
	SendResetPasswordEmail(data EmailData) error
}

type EmailData struct {
	Subject string
	Email   string
	Name    string
	Code    string
}
