package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/repos"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	cases := []struct {
		Incoming          LoginReq
		Expected          LoginResp
		CredentialChecker repos.CredentialChecker
	}{
		{
			Incoming: LoginReq{
				Email: "kev.legoff@gmail.com",
				Pwd:   "test1234",
			},
			Expected: LoginResp{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{ErrorCode: 0, Error: ""},
				},
			},
			CredentialChecker: &validCredentialMock{},
		},
		{
			Incoming: LoginReq{
				Email: "",
				Pwd:   "test1234",
			},
			Expected: LoginResp{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{ErrorCode: 101, Error: ""},
				},
			},
			CredentialChecker: &validCredentialMock{},
		},
		{
			Incoming: LoginReq{
				Email: "kev.legoff@gmail.com",
				Pwd:   "",
			},
			Expected: LoginResp{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{ErrorCode: 102, Error: ""},
				},
			},
			CredentialChecker: &validCredentialMock{},
		},
		{
			Incoming: LoginReq{
				Email: "kev.legoff@gmail.com",
				Pwd:   "test1234",
			},
			Expected: LoginResp{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{ErrorCode: 103, Error: ""},
				},
			},
			CredentialChecker: &invalidEmailCredentialMock{},
		},
		{
			Incoming: LoginReq{
				Email: "kev.legoff@gmail.com",
				Pwd:   "test1234",
			},
			Expected: LoginResp{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{ErrorCode: 104, Error: ""},
				},
			},
			CredentialChecker: &invalidPasswordCredentialMock{},
		},
	}
	for _, test := range cases {
		form := url.Values{}
		form.Add("email", test.Incoming.Email)
		form.Add("pwd", test.Incoming.Pwd)
		// Build a request
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		context.Set(req, repos.UserR, test.CredentialChecker)

		rw := httptest.NewRecorder()
		Login(rw, req)
		// Check we got a 200
		if rw.Code != 200 {
			t.Error("200 was expected but got ", rw.Code)
		}
		decoder := json.NewDecoder(rw.Body)
		resp := LoginResp{}
		if err := decoder.Decode(&resp); err != nil {
			t.Errorf("Could not decode the answer")
		}
		if resp.Error.ErrorCode != test.Expected.Error.ErrorCode {
			t.Error("Error code expected ", test.Expected.Error.ErrorCode, "got ", resp.Error.ErrorCode)
		}
	}
}

type validCredentialMock struct{}

func (repo validCredentialMock) FindByEmail(email string) bool {
	return true
}

func (repo validCredentialMock) CheckEmailPassword(email string, password string) (*data.User, *data.AuthToken, error) {
	token := data.NewAuthToken()
	return &data.User{}, &token, nil
}

type invalidEmailCredentialMock struct{}

func (repo invalidEmailCredentialMock) FindByEmail(email string) bool {
	return false
}
func (repo invalidEmailCredentialMock) CheckEmailPassword(email string, password string) (*data.User, *data.AuthToken, error) {
	// Never reached
	return nil, nil, nil
}

type invalidPasswordCredentialMock struct{}

func (repo invalidPasswordCredentialMock) FindByEmail(email string) bool {
	return true
}
func (repo invalidPasswordCredentialMock) CheckEmailPassword(email string, password string) (*data.User, *data.AuthToken, error) {
	return nil, nil, errors.New("user not found")
}
