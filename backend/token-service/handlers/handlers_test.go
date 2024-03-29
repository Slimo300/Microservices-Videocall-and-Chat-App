package handlers_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/auth"
	dblayer "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/token-service/database"
	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/token-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/token-service/handlers"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var PrivateKey *rsa.PrivateKey
var RefreshSecret string

var db *mockdb.MockTokenDB
var service *handlers.TokenService

func TestMain(m *testing.M) {
	// generate key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Cannot generate RSA key\n")
	}
	PrivateKey = priv

	db = new(mockdb.MockTokenDB)
	db.On("SaveToken", mock.AnythingOfType("string"), time.Hour*24).Return(nil)

	service = handlers.NewTokenService(db, priv, RefreshSecret, time.Hour*24, time.Minute*20)

	os.Exit(m.Run())
}

func TestNewPairFromUserID(t *testing.T) {

	userID := uuid.NewString()

	tokens, err := service.NewPairFromUserID(context.Background(), &auth.UserID{ID: userID})
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
		checkAssertions func(tokens *auth.TokenPair)
		prepare         func(m *mock.Mock)
	}{
		{
			desc:  "PairFromRefresh Success",
			token: token.Token,
			checkAssertions: func(tokens *auth.TokenPair) {
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
			checkAssertions: func(tokens *auth.TokenPair) {
				assert.Empty(t, tokens.AccessToken)
				assert.Empty(t, tokens.RefreshToken)
				assert.NotEmpty(t, tokens.Error)
			},
			prepare: func(m *mock.Mock) {},
		},
		{
			desc:  "PairFromRefresh Blacklisted",
			token: token.Token,
			checkAssertions: func(tokens *auth.TokenPair) {
				assert.Empty(t, tokens.AccessToken)
				assert.Empty(t, tokens.RefreshToken)
				assert.Equal(t, dblayer.ErrTokenBlacklisted.Error(), tokens.Error)
			},
			prepare: func(m *mock.Mock) {
				m.On("IsTokenValid", userID, mock.AnythingOfType("string")).Return(false, dblayer.ErrTokenBlacklisted).Once()
				m.On("InvalidateTokens", userID, mock.AnythingOfType("string")).Return(nil).Once()
			},
		},
		{
			desc:  "PairFromRefresh NotFound",
			token: token.Token,
			checkAssertions: func(tokens *auth.TokenPair) {
				assert.Empty(t, tokens.AccessToken)
				assert.Empty(t, tokens.RefreshToken)
				assert.Equal(t, dblayer.ErrTokenNotFound.Error(), tokens.Error)
			},
			prepare: func(m *mock.Mock) {
				m.On("IsTokenValid", userID, mock.AnythingOfType("string")).Return(false, dblayer.ErrTokenNotFound).Once()
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			tC.prepare(&db.Mock)

			tokens, err := service.NewPairFromRefresh(context.Background(), &auth.RefreshToken{Token: tC.token})
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
			expectedResponse: dblayer.ErrTokenNotFound.Error(),
			prepare: func(m *mock.Mock) {
				m.On("InvalidateToken", userID, mock.AnythingOfType("string")).Return(dblayer.ErrTokenNotFound).Once()
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			tC.prepare(&db.Mock)

			res, err := service.DeleteUserToken(context.Background(), &auth.RefreshToken{Token: token.Token})
			if err != nil {
				t.Errorf("Method returned an error: %v", err.Error())
			}

			assert.Equal(t, tC.expectedResponse, res.Error)
		})
	}

}
