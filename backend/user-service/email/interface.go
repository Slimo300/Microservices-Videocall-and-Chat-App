package email

type EmailService interface {
	SendVerificationEmail(data VerificationEmailData) error
}

type VerificationEmailData struct {
	Email            string
	Name             string
	VerificationCode string
	UserID           string
}
