package follow

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
)

type IncomingFollow struct {
	Token string `schema:"token"`
	To    string `schema:"to"`
}

func (incFollow *IncomingFollow) Ok() error {
	if _, err := repos.BuildIdFromString(incFollow.To); err != nil {
		log.Println("Error in Ok()", err.Error())
		return errors.New("invalid_to_id")
	}
	return nil
}

type FollowedResponse struct {
	v0.Response
	Followed *descriptor.FolloweeDescriptor `json:"followed"`
}

func PostFollow(rw http.ResponseWriter, req *http.Request) {
	var incoming IncomingFollow
	var response FollowedResponse
	userRepo := context.Get(req, "userRepo").(repos.UserIdGetter)
	fromUser := context.Get(req, "user").(*data.User)
	if err := handlers.DecodeForm(req, &incoming); err != nil {
		log.Printf("Error detected %s \n", err.Error())
		response.Error = handlers.ErrorDescriptor{ErrorCode: 101, Error: err.Error()}
		handlers.Respond(rw, req, http.StatusOK, &response)
		return
	}
	uId, _ := repos.BuildIdFromString(incoming.To)
	toU, err := userRepo.GetUserById(uId)
	if err != nil {
		log.Println("Not able to find the to user to follow")
		response.Error = handlers.ErrorDescriptor{ErrorCode: 102, Error: err.Error()}
		handlers.Respond(rw, req, http.StatusOK, &response)
		return
	}
	log.Println("From receieved ", fromUser.Following, "to received ", toU.ID)
	// Test the user is not following to
	if fromUser.IsFollowing(toU.ID) {
		log.Println("User is already followed")
		response.Error = handlers.ErrorDescriptor{ErrorCode: 103, Error: "User is already followed"}
		handlers.Respond(rw, req, http.StatusOK, &response)
		return
	}
	followCreator := context.Get(req, "userRepo").(repos.FollowEstablisher)
	_, to, err := followCreator.EstablishFollow(fromUser.ID, toU.ID)
	if err != nil {
		log.Println("Unable to save follow relationship")
		response.Error = handlers.ErrorDescriptor{ErrorCode: 104, Error: "DB problem " + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, &response)
		return
	}
	response.Followed = descriptor.DescribeFollowee(to)
	handlers.Respond(rw, req, http.StatusOK, &response)
}
