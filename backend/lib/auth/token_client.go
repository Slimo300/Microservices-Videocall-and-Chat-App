package auth

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
	keyid  string
}

// NewGRPCTokenClient is a constructor for grpc client to reach token service
func NewGRPCTokenClient(port string) (TokenClient, error) {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
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
		keyid:  pubKeyMsg.Iteration,
	}, nil
}

func (c *gRPCTokenAuthClient) NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error) {
	response, err := c.client.NewPairFromUserID(context.Background(), &pb.UserID{ID: userID.String()})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (c *gRPCTokenAuthClient) NewPairFromRefresh(refresh string) (*pb.TokenPair, error) {
	response, err := c.client.NewPairFromRefresh(context.Background(), &pb.RefreshToken{Token: refresh})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (c *gRPCTokenAuthClient) DeleteUserToken(refresh string) error {
	response, err := c.client.DeleteUserToken(context.Background(), &pb.RefreshToken{Token: refresh})
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func (c *gRPCTokenAuthClient) GetPublicKey(keyID string) (*rsa.PublicKey, error) {

	if keyID == c.keyid {
		return &c.pubKey, nil
	}

	pubKeyMsg, err := c.client.GetPublicKey(context.Background(), &pb.Empty{})
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

	return publicKey, nil
}
