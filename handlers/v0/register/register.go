package register

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"github.com/softdev9/trendee-api-master/utils"
	"log"
	"net/http"
	"net/mail"
)

type IncomingCredential struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (cred *IncomingCredential) Ok() error {
	log.Println("[DEBUG] Credential ", cred)
	if len(cred.Email) > 0 {
		if _, err := mail.ParseAddress(cred.Email); err != nil {
			return errors.New("Invalid email given")
		}
	} else {
		return errors.New("No email specified")
	}
	if len(cred.Password) < 4 {
		return errors.New("No password given, password need to be a least 5 char")
	}
	return nil
}

type UserRegistrationResponse struct {
	v0.Response
	User  *descriptor.UserLight `json:"user"`
	Token *data.AuthToken       `json:"auth_token"`
}

func RegisterPost(rw http.ResponseWriter, req *http.Request) {
	incoming := IncomingCredential{}
	resp := UserRegistrationResponse{}
	if err := handlers.DecodeForm(req, &incoming); err != nil {
		resp.Response.Error = handlers.ErrorDescriptor{
			Error:     err.Error(),
			ErrorCode: 401,
		}
	}
	// If the user exist we return an error as well
	mailVerifier := context.Get(req, repos.UserR).(repos.EmailUserVerifier)
	if mailVerifier.FindByEmail(incoming.Email) {
		resp.Response.Error = handlers.ErrorDescriptor{
			Error:     "User already exists",
			ErrorCode: http.StatusForbidden,
		}
	} else {
		// We create the user
		u := data.NewUser(incoming.Email, incoming.Password, utils.RandSeq(32))
		// We save the user
		userSaver := context.Get(req, repos.UserR).(repos.UserSaver)
		userSaver.SaveUser(&u)
		// We create a token
		token := data.NewAuthToken()
		tokenRepo := context.Get(req, repos.TokenR).(repos.TokenRepo)
		if err := tokenRepo.SaveAuthToken(token, u); err != nil {
			resp.Response.Error = handlers.ErrorDescriptor{
				Error:     err.Error(),
				ErrorCode: http.StatusInternalServerError,
			}
		} else {
			resp.Token = &token
		}
		resp.User = descriptor.NewUserLight(&u, false, true)
	}
	handlers.Respond(rw, req, http.StatusOK, &resp)
}
