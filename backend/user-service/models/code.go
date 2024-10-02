package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthorizationCode struct {
	code     uuid.UUID
	userID   uuid.UUID
	created  time.Time
	codeType CodeType
}

func (c AuthorizationCode) Code() uuid.UUID      { return c.code }
func (c AuthorizationCode) UserID() uuid.UUID    { return c.userID }
func (c AuthorizationCode) CreatedAt() time.Time { return c.created }
func (c AuthorizationCode) Type() CodeType       { return c.codeType }

func NewAuthorizationCode(userID uuid.UUID, codeType CodeType) *AuthorizationCode {
	return &AuthorizationCode{
		code:     uuid.New(),
		userID:   userID,
		created:  time.Now(),
		codeType: codeType,
	}
}

// UnmarshalVerificationCodeFromDatabase should be only called from database layer
func UnmarshalAuthorizationCodeFromDatabase(code, userID uuid.UUID, created time.Time, codeType string) (*AuthorizationCode, error) {
	t, err := NewCodeTypeFromString(codeType)
	if err != nil {
		return &AuthorizationCode{}, err
	}
	return &AuthorizationCode{
		code:     code,
		userID:   userID,
		created:  created,
		codeType: t,
	}, nil
}

type CodeType struct {
	s string
}

func (c CodeType) String() string { return c.s }
func (c CodeType) IsZero() bool   { return c == CodeType{} }

var (
	EmailVerificationCode = CodeType{"email_verification"}
	ResetPasswordCode     = CodeType{"reset_password"}
)

var codeTypeValues = []CodeType{
	EmailVerificationCode, ResetPasswordCode,
}

func NewCodeTypeFromString(codeStr string) (CodeType, error) {
	for _, codeType := range codeTypeValues {
		if codeType.String() == codeStr {
			return codeType, nil
		}
	}
	return CodeType{}, fmt.Errorf("unknown '%s' code type", codeStr)
}
