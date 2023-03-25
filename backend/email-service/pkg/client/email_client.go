package client

import (
	"context"
	"errors"

	"github.com/Slimo300/chat-emailservice/pkg/client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type emailClient struct {
	client pb.EmailServiceClient
}

// SendVerificationEmail sends rpc call to grpc server and checks for error
func (e *emailClient) SendVerificationEmail(ctx context.Context, data *pb.EmailData) error {
	res, err := e.client.SendVerificationEmail(ctx, data)
	if err != nil {
		return err
	}

	if res.Error != "" {
		return errors.New(res.Error)
	}

	return nil
}

// SendResetPasswordEmail sends rpc call to grpc server and checks for error
func (e *emailClient) SendResetPasswordEmail(ctx context.Context, data *pb.EmailData) error {
	res, err := e.client.SendResetPasswordEmail(ctx, data)
	if err != nil {
		return err
	}

	if res.Error != "" {
		return errors.New(res.Error)
	}

	return nil
}

// NewGRPCEmailClient returns new email client with established grpc connection
func NewGRPCEmailClient(address string) (EmailClient, error) {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &emailClient{
		client: pb.NewEmailServiceClient(conn),
	}, nil
}
