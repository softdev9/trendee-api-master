package v0

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	// "github.com/softdev9/trendee-api-master/gateways"
	"fmt"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/repos"
	"github.com/softdev9/trendee-api-master/utils"
	"net/http"
	"net/mail"
)

type IncomingCredential struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (cred *IncomingCredential) Ok() error {
	fmt.Println("Credential ", cred)
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

func PostRegistration(rw http.ResponseWriter, req *http.Request) {
	incoming := IncomingCredential{}
	if err := handlers.DecodeForm(req, &incoming); err != nil {
		handlers.Respond(rw, req, http.StatusOK,
			&UserRegistrationResponse{
				Error: handlers.ErrorDescriptor{
					Error:     err.Error(),
					ErrorCode: 400,
				},
			},
		)
		return
	}
	// We get copy of a repo
	userRepo := context.Get(req, "userRepo").(repos.EmailUserVerifier)
	if found := userRepo.FindByEmail(incoming.Email); found {
		handlers.Respond(rw, req,
			http.StatusOK,
			&UserRegistrationResponse{
				Error: handlers.ErrorDescriptor{
					Error:     "An account for this user already exists",
					ErrorCode: 200,
				},
			})
		return
	}
	userSaver := context.Get(req, "userRepo").(repos.UserSaver)
	tokenRepo := context.Get(req, "tokenRepo").(repos.TokenRepo)
	u := data.NewUser(incoming.Email, incoming.Password, utils.RandSeq(32))
	userSaver.SaveUser(&u)
	token := data.NewAuthToken()
	if err := tokenRepo.SaveAuthToken(token, u); err != nil {
		handlers.Respond(
			rw, req,
			http.StatusInternalServerError,
			&UserRegistrationResponse{
				Error: handlers.ErrorDescriptor{
					Error:     "",
					ErrorCode: 0,
				},
			},
		)
	}
	resp := NewRegistrationResponse(u, token)
	handlers.Respond(rw, req, http.StatusOK, resp)
	// mailSender := context.Get(req, "mailSender").(gateways.TrendeeMailSender)
	// mailSender.SendWelcomeMessage(u)
}
