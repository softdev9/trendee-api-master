//+build sasl

package mgo

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/internal/sasl"
)

func saslNew(cred Credential, host string) (saslStepper, error) {
	return sasl.New(cred.Username, cred.Password, cred.Mechanism, cred.Service, host)
}
