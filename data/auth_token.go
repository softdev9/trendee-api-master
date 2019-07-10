package data

import (
	"github.com/spiderdev86/trendee-api/Godeps/_workspace/src/github.com/pborman/uuid"
	"time"
)

type AuthToken struct {
	Token     string    `json:"token" bson:"token"`
	ExpiresOn time.Time `json:"expires_on" bson:"expires_on"`
}

func NewAuthToken() AuthToken {
	expires := time.Now().AddDate(0, 0, 21).UTC()
	uuid := uuid.New()
	return AuthToken{Token: uuid, ExpiresOn: expires}
}

func (t AuthToken) OK() bool {
	if len(t.Token) == 0 {
		return false
	}
	if time.Now().UTC().After(t.ExpiresOn) == true {
		return false
	}
	return true
}
