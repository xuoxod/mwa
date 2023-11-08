package sms

import (
	"github.com/xuoxod/mwa/internal/config"
)

type CreateMessageFeedbackParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) associated with the Message resource for which to create MessageFeedback.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
	//
	Outcome *string `json:"Outcome,omitempty"`
}

func SendVerifyToken(to, msg string, app *config.AppConfig) error {

	return nil
}
