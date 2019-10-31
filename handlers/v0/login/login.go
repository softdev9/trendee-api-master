package login

import (
	"errors"
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"net/http"
)

type LoginReq struct {
	Email string `schema:"email"`
	Pwd   string `schema:"pwd"`
}

// Login Response
type LoginResp struct {
	v0.Response
	User  *descriptor.UserLight `json:"user"`
	Token *data.AuthToken       `json:"auth_token"`
}

func (toCheck *LoginReq) Ok() error {
	fmt.Println("To check ", toCheck)
	if len(toCheck.Email) == 0 {
		return errors.New("no email")
	}
	if len(toCheck.Pwd) == 0 {
		return errors.New("no pwd")
	}
	return nil
}

func Login(rw http.ResponseWriter, req *http.Request) {
	resp := LoginResp{}
	loginReq := LoginReq{}
	errorD := handlers.ErrorDescriptor{}
	// Decode the incoming like article request
	if err := handlers.DecodeForm(req, &loginReq); err != nil {
		if err.Error() == "no email" {
			errorD.ErrorCode = 101
			errorD.Error = "Invalid email sent"
		}
		if err.Error() == "no pwd" {
			errorD.ErrorCode = 102
			errorD.Error = "Invalid password sent"
		}
		resp.Response.Error = errorD
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	// Check the user with email exist
	credentialChecker := context.Get(req, repos.UserR).(repos.CredentialChecker)
	if !credentialChecker.FindByEmail(loginReq.Email) {
		errorD.ErrorCode = 103
		errorD.Error = "No user found in db for this email"
		resp.Response.Error = errorD
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	u, token, err := credentialChecker.CheckEmailPassword(loginReq.Email, loginReq.Pwd)
	if err != nil {
		errorD.ErrorCode = 104
		errorD.Error = "Invalid password"
		resp.Response.Error = errorD
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	resp.User = descriptor.NewUserLight(u, false, true)
	resp.Token = token
	// Send a request and check in answer with a 200
	handlers.Respond(rw, req, http.StatusOK, &resp)
}
