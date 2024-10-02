package app

import (
	"context"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/email"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/amqp"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage/s3"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app/command"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/app/query"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/config"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database/orm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mockqueue "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/mock"
	mockstorage "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/storage/mock"
	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database/mock"
)

type Commands struct {
	ChangePassword       command.ChangePasswordHandler
	SetProfilePicture    command.SetProfilePictureHandler
	DeleteProfilePicture command.DeleteProfilePictureHandler

	ForgotPassword         command.ForgotPasswordHandler
	ResetForgottenPassword command.ResetForgottenPasswordHandler

	RegisterUser command.RegisterUserHandler
	VerifyEmail  command.VerifyEmailHandler

	SignIn       command.SignInHandler
	SignOut      command.SignOutHandler
	RefreshToken command.RefreshTokenHandler
}

type Queries struct {
	GetUser query.GetUserHandler
}

type App struct {
	Commands Commands
	Queries  Queries
}

func NewTestApplication() App {
	repo := new(mockdb.UsersMockRepository)
	emitter := new(mockqueue.MockEmitter)
	storage := new(mockstorage.MockStorage)
	emailClient := new(email.MockEmailClient)
	tokenClient := new(auth.MockTokenClient)

	return App{
		Queries: Queries{
			GetUser: query.NewGetUserHandler(repo),
		},
		Commands: Commands{
			ChangePassword:         command.NewChangePasswordHandler(repo),
			SetProfilePicture:      command.NewSetProfilePictureHandler(repo, storage, emitter),
			DeleteProfilePicture:   command.NewDeleteProfilePictureHandler(repo, storage, emitter),
			ForgotPassword:         command.NewForgotPasswordHandler(repo, emailClient),
			ResetForgottenPassword: command.NewResetForgottenPasswordHandler(repo),
			RegisterUser:           command.NewRegisterUserHandler(repo, emailClient),
			SignIn:                 command.NewSignInHandler(repo, tokenClient),
			SignOut:                command.NewSignOutHandler(tokenClient),
			RefreshToken:           command.NewRefreshTokenHandler(tokenClient),
		},
	}
}

func NewApplication(conf config.Config) App {
	repo, err := orm.NewUsersGormRepository(conf.DBAddress)
	if err != nil {
		panic(err)
	}
	tokenConn, err := grpc.Dial(conf.TokenServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	tokenClient := auth.NewTokenServiceClient(tokenConn)

	emailConn, err := grpc.Dial(conf.EmailServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	emailClient := email.NewEmailServiceClient(emailConn)

	storage, err := s3.NewS3Storage(context.Background(), conf.StorageKeyID, conf.StorageKeySecret, conf.Bucket, s3.WithRegion(conf.StorageRegion))
	if err != nil {
		panic(err)
	}

	builder, err := amqp.NewAMQPBuilder(conf.BrokerAddress)
	if err != nil {
		panic(err)
	}
	emitter, err := builder.GetEmiter(msgqueue.EmiterConfig{
		ExchangeName: "user",
	})
	if err != nil {
		panic(err)
	}

	return App{
		Queries: Queries{
			GetUser: query.NewGetUserHandler(repo),
		},
		Commands: Commands{
			ChangePassword:         command.NewChangePasswordHandler(repo),
			SetProfilePicture:      command.NewSetProfilePictureHandler(repo, storage, emitter),
			DeleteProfilePicture:   command.NewDeleteProfilePictureHandler(repo, storage, emitter),
			ForgotPassword:         command.NewForgotPasswordHandler(repo, emailClient),
			ResetForgottenPassword: command.NewResetForgottenPasswordHandler(repo),
			RegisterUser:           command.NewRegisterUserHandler(repo, emailClient),
			SignIn:                 command.NewSignInHandler(repo, tokenClient),
			SignOut:                command.NewSignOutHandler(tokenClient),
			RefreshToken:           command.NewRefreshTokenHandler(tokenClient),
		},
	}
}
