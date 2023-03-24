package client

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth/pb"
)

type gRPCTokenAuthClient struct {
	client pb.TokenServiceClient
	pubKey rsa.PublicKey
}

// NewGRPCTokenClient is a constructor for grpc client to reach token service
func NewGRPCTokenClient(address string) (TokenClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewTokenServiceClient(conn)

	pubKeyMsg, err := client.GetPublicKey(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, err
	}
	if pubKeyMsg.Error != "" {
		return nil, fmt.Errorf("Message from token service: %v", pubKeyMsg.Error)
	}

	publicKeyParsed, err := x509.ParsePKIXPublicKey(pubKeyMsg.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyParsed.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("PublicKey not of type *rsa.PublicKey")
	}

	return &gRPCTokenAuthClient{
		client: client,
		pubKey: *publicKey,
	}, nil
}

func (c *gRPCTokenAuthClient) NewPairFromUserID(ctx context.Context, userID uuid.UUID) (*pb.TokenPair, error) {
	response, err := c.client.NewPairFromUserID(ctx, &pb.UserID{ID: userID.String()})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (c *gRPCTokenAuthClient) NewPairFromRefresh(ctx context.Context, refresh string) (*pb.TokenPair, error) {
	response, err := c.client.NewPairFromRefresh(ctx, &pb.RefreshToken{Token: refresh})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (c *gRPCTokenAuthClient) DeleteUserToken(ctx context.Context, refresh string) error {
	response, err := c.client.DeleteUserToken(ctx, &pb.RefreshToken{Token: refresh})
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func (c *gRPCTokenAuthClient) GetPublicKey() *rsa.PublicKey {

	return &c.pubKey
}
