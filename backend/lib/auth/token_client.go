package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/pb"
)

type gRPCTokenAuthClient struct {
	client    pb.TokenServiceClient
	pubKey    rsa.PublicKey
	iteration string
}

func NewGRPCTokenAuthClient() (*gRPCTokenAuthClient, error) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		return &gRPCTokenAuthClient{}, err
	}
	client := pb.NewTokenServiceClient(conn)

	pubKeyMsg, err := client.GetPublicKey(context.Background(), &pb.Empty{})
	if err != nil {
		return &gRPCTokenAuthClient{}, err
	}
	if pubKeyMsg.Error != "" {
		return &gRPCTokenAuthClient{}, errors.New(fmt.Sprintf("Message from token service: %v", pubKeyMsg.Error))
	}

	publicKeyParsed, err := x509.ParsePKIXPublicKey(pubKeyMsg.PublicKey)
	if err != nil {
		return &gRPCTokenAuthClient{}, err
	}

	publicKey, ok := publicKeyParsed.(*rsa.PublicKey)
	if !ok {
		return &gRPCTokenAuthClient{}, errors.New("PublicKey not of type *rsa.PublicKey")
	}

	return &gRPCTokenAuthClient{
		client: client,
		pubKey: *publicKey,
	}, nil
}

func (grpc *gRPCTokenAuthClient) NewPairFromUserID(userID uuid.UUID) (*pb.TokenPair, error) {
	response, err := grpc.client.NewPairFromUserID(context.Background(), &pb.UserID{ID: userID.String()})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (grpc *gRPCTokenAuthClient) NewPairFromRefresh(refresh string) (*pb.TokenPair, error) {
	response, err := grpc.client.NewPairFromRefresh(context.Background(), &pb.RefreshToken{Token: refresh})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (grpc *gRPCTokenAuthClient) DeleteUserToken(refresh string) error {
	response, err := grpc.client.DeleteUserToken(context.Background(), &pb.RefreshToken{Token: refresh})
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func (grpc *gRPCTokenAuthClient) GetPublicKey() *rsa.PublicKey {
	return &grpc.pubKey
}
