package v0

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"net/http"
)

type UserResponse struct {
	Response
	User *descriptor.UserLight `json:"user"`
}

func GetProfile(rw http.ResponseWriter, req *http.Request) {
	// Create an empty response
	// Check the token is valid
	params := req.URL.Query()
	if len(params["token"]) == 1 {
		tokenRepo := context.Get(req, "tokenRepo").(repos.TokenRepo)
		userRepo := context.Get(req, "userRepo").(*repos.MongoUserRepo)
		token := params["token"][0]
		idUser, err := tokenRepo.GetUserWithToken(token)
		if err != nil {
			//reportError(rw, req, &UserResponse{}, err)
			return
		}
		u, err := userRepo.GetUserById(idUser)
		resp := UserResponse{}
		resp.User = descriptor.NewUserLight(u, false, true)
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}

}
