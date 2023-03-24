package handlers_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Slimo300/chat-tokenservice/internal/handlers"
	repolayer "github.com/Slimo300/chat-tokenservice/internal/repo"
	mockrepo "github.com/Slimo300/chat-tokenservice/internal/repo/mock"
	"github.com/Slimo300/chat-tokenservice/pkg/client/pb"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var PrivateKey *rsa.PrivateKey
var RefreshSecret string

var repo *mockrepo.MockTokenRepository
var service *handlers.TokenService

func TestMain(m *testing.M) {
	// generate key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Cannot generate RSA key\n")
	}
	PrivateKey = priv

	RefreshSecret = "secret"

	repo = new(mockrepo.MockTokenRepository)
	repo.On("SaveToken", mock.AnythingOfType("string"), time.Hour*24).Return(nil)
	repo.On("GetPrivateKey").Return(priv, nil)

	service, err = handlers.NewTokenService(repo, RefreshSecret, time.Hour*24, time.Minute*20)
	if err != nil {
		log.Fatal("Couldn't create token service")
	}

	os.Exit(m.Run())
}

func TestNewPairFromUserID(t *testing.T) {

	userID := uuid.NewString()

	tokens, err := service.NewPairFromUserID(context.Background(), &pb.UserID{ID: userID})
	if err != nil {
		t.Errorf("Service returned error: %s", err.Error())
	}

	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Empty(t, tokens.Error)

	accessToken, err := jwt.ParseWithClaims(tokens.AccessToken, &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return &PrivateKey.PublicKey, nil
		})
	if err != nil {
		t.Errorf("Error occured when parsing token from string: %v: ", err.Error())
	}
	accessTokenUserID := accessToken.Claims.(*jwt.StandardClaims).Subject
	assert.Equal(t, userID, accessTokenUserID)

	refreshToken, err := jwt.ParseWithClaims(tokens.RefreshToken, &jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(RefreshSecret), nil
		})
	if err != nil {
		t.Errorf("Error occured when parsing token from string: %v: ", err.Error())
	}
	refreshTokenUserID := refreshToken.Claims.(*jwt.StandardClaims).Subject
	assert.Equal(t, userID, refreshTokenUserID)

}

func TestNewPairFromRefresh(t *testing.T) {

	userID := uuid.NewString()

	token, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Errorf("Error when generating test refresh token: %v", err)
	}

	testCases := []struct {
		desc            string
		token           string
		checkAssertions func(tokens *pb.TokenPair)
		prepare         func(m *mock.Mock)
	}{
		{
			desc:  "PairFromRefresh Success",
			token: token.Token,
			checkAssertions: func(tokens *pb.TokenPair) {
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
				assert.Empty(t, tokens.Error)
				accessToken, err := jwt.ParseWithClaims(tokens.AccessToken, &jwt.StandardClaims{},
					func(t *jwt.Token) (interface{}, error) {
						return &PrivateKey.PublicKey, nil
					})
				if err != nil {
					t.Errorf("Error occured when parsing token from string: %v: ", err.Error())
				}
				accessTokenUserID := accessToken.Claims.(*jwt.StandardClaims).Subject
				assert.Equal(t, userID, accessTokenUserID)

				refreshToken, err := jwt.ParseWithClaims(tokens.RefreshToken, &jwt.StandardClaims{},
					func(t *jwt.Token) (interface{}, error) {
						return []byte(RefreshSecret), nil
					})
				if err != nil {
					t.Errorf("Error occured when parsing token from string: %v: ", err.Error())
				}
				refreshTokenUserID := refreshToken.Claims.(*jwt.StandardClaims).Subject
				assert.Equal(t, userID, refreshTokenUserID)

			},
			prepare: func(m *mock.Mock) {
				m.On("IsTokenValid", userID, mock.AnythingOfType("string")).Return(true, nil).Once()
				m.On("InvalidateToken", userID, mock.AnythingOfType("string")).Return(nil).Once()
				m.On("SaveToken", mock.AnythingOfType("string"), time.Hour*24).Return(nil).Once()
			},
		},
		{
			desc:  "PairFromRefresh InvalidToken",
			token: "",
			checkAssertions: func(tokens *pb.TokenPair) {
				assert.Empty(t, tokens.AccessToken)
				assert.Empty(t, tokens.RefreshToken)
				assert.NotEmpty(t, tokens.Error)
			},
			prepare: func(m *mock.Mock) {},
		},
		{
			desc:  "PairFromRefresh Blacklisted",
			token: token.Token,
			checkAssertions: func(tokens *pb.TokenPair) {
				assert.Empty(t, tokens.AccessToken)
				assert.Empty(t, tokens.RefreshToken)
				assert.Equal(t, repolayer.TokenBlacklistedError.Error(), tokens.Error)
			},
			prepare: func(m *mock.Mock) {
				m.On("IsTokenValid", userID, mock.AnythingOfType("string")).Return(false, repolayer.TokenBlacklistedError).Once()
				m.On("InvalidateTokens", userID, mock.AnythingOfType("string")).Return(nil).Once()
			},
		},
		{
			desc:  "PairFromRefresh NotFound",
			token: token.Token,
			checkAssertions: func(tokens *pb.TokenPair) {
				assert.Empty(t, tokens.AccessToken)
				assert.Empty(t, tokens.RefreshToken)
				assert.Equal(t, repolayer.TokenNotFoundError.Error(), tokens.Error)
			},
			prepare: func(m *mock.Mock) {
				m.On("IsTokenValid", userID, mock.AnythingOfType("string")).Return(false, repolayer.TokenNotFoundError).Once()
			},
		},
		{
			desc:  "PairFromRefresh TooManyTokens",
			token: token.Token,
			checkAssertions: func(tokens *pb.TokenPair) {
				assert.Empty(t, tokens.AccessToken)
				assert.Empty(t, tokens.RefreshToken)
				assert.Equal(t, repolayer.TooManyTokensFoundError.Error(), tokens.Error)
			},
			prepare: func(m *mock.Mock) {
				m.On("IsTokenValid", userID, mock.AnythingOfType("string")).Return(false, repolayer.TooManyTokensFoundError).Once()
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			tC.prepare(&repo.Mock)

			tokens, err := service.NewPairFromRefresh(context.Background(), &pb.RefreshToken{Token: tC.token})
			if err != nil {
				t.Errorf("Service returned error: %v", err.Error())
			}

			tC.checkAssertions(tokens)

		})
	}

}

func TestDeleteUserToken(t *testing.T) {

	userID := uuid.NewString()

	token, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Error("Couldn't generate refresh token")
	}

	testCases := []struct {
		desc             string
		expectedResponse string
		prepare          func(m *mock.Mock)
	}{
		{
			desc:             "Delete Token Success",
			expectedResponse: "",
			prepare: func(m *mock.Mock) {
				m.On("InvalidateToken", userID, mock.AnythingOfType("string")).Return(nil).Once()
			},
		},
		{
			desc:             "Delete Token No Token",
			expectedResponse: repolayer.TokenNotFoundError.Error(),
			prepare: func(m *mock.Mock) {
				m.On("InvalidateToken", userID, mock.AnythingOfType("string")).Return(repolayer.TokenNotFoundError).Once()
			},
		},
		{
			desc:             "Delete Token Many Tokens",
			expectedResponse: repolayer.TooManyTokensFoundError.Error(),
			prepare: func(m *mock.Mock) {
				m.On("InvalidateToken", userID, mock.AnythingOfType("string")).Return(repolayer.TooManyTokensFoundError).Once()
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			tC.prepare(&repo.Mock)

			res, err := service.DeleteUserToken(context.Background(), &pb.RefreshToken{Token: token.Token})
			if err != nil {
				t.Errorf("Method returned an error: %v", err.Error())
			}

			assert.Equal(t, tC.expectedResponse, res.Error)
		})
	}

}

func TestGetPublicKey(t *testing.T) {

	pubKey, err := service.GetPublicKey(context.Background(), &pb.Empty{})
	if err != nil {
		t.Errorf("Service returned error: %v", err.Error())
	}

	assert.NotEmpty(t, pubKey.PublicKey)
	assert.Empty(t, pubKey.Error)

}
