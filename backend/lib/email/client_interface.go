package email

import (
	"context"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/email/pb"
)

// EmailClient wraps Email Service Client and exposes its functionality
type EmailClient interface {
	SendVerificationEmail(ctx context.Context, in *pb.EmailData) error
	SendResetPasswordEmail(ctx context.Context, in *pb.EmailData) error
}
