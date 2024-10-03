package events

import (
	"github.com/google/uuid"
)

type UserRegisteredEvent struct {
	Email    string `json:"email" mapstructure:"email"`
	Username string `json:"username" mapstructure:"username"`
	Code     string `json:"code" mapstructure:"code"`
}

func (UserRegisteredEvent) EventName() string { return "user.created" }

type UserVerifiedEvent struct {
	ID       uuid.UUID `json:"userID" mapstructure:"userID"`
	Username string    `json:"username" mapstructure:"username"`
}

func (UserVerifiedEvent) EventName() string { return "user.verified" }

type UserForgotPasswordEvent struct {
	Email    string `json:"email" mapstructure:"email"`
	Username string `json:"username" mapstructure:"username"`
	Code     string `json:"code" mapstructure:"code"`
}

func (UserForgotPasswordEvent) EventName() string { return "user.forgottenpassword" }
