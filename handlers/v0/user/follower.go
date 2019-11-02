package user

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
)

type FollowerListResponse struct {
	v0.Response
	Followers []*descriptor.UserLight `json:"followers"`
}

type FollowingListResponse struct {
	v0.Response
	Following []*descriptor.UserLight `json:"following"`
}

func GetFollowers(rw http.ResponseWriter, req *http.Request) {
	asker := context.Get(req, "user").(*data.User)
	userRepo := context.Get(req, repos.UserR).(repos.ListFinder)
	var target *data.User
	if user, _ := ParseUserId(req, userRepo.(repos.UserRepo)); user != nil {
		target = user
	} else {
		target = asker
	}
	resp := FollowerListResponse{}
	resp.Followers = buildFollowerList(userRepo, asker.ID, target.FollowedBy)
	handlers.Respond(rw, req, http.StatusOK, &resp)
}

func GetFollowing(rw http.ResponseWriter, req *http.Request) {
	asker := context.Get(req, "user").(*data.User)
	userRepo := context.Get(req, repos.UserR).(repos.ListFinder)
	var target *data.User
	if user, _ := ParseUserId(req, userRepo.(repos.UserRepo)); user != nil {
		target = user
	} else {
		target = asker
	}
	resp := FollowingListResponse{}
	resp.Following = buildFollowerList(userRepo, asker.ID, target.Following)
	handlers.Respond(rw, req, http.StatusOK, &resp)
}

func buildFollowerList(userSrc repos.ListFinder, asker bson.ObjectId, listToFormat []bson.ObjectId) []*descriptor.UserLight {
	users, err := userSrc.FindFromList(listToFormat)
	if err != nil {
		log.Printf("[ERROR] Error while retriving the follower list %s \n", err.Error())
		return []*descriptor.UserLight{}
	}
	if len(users) == 0 {
		log.Println("[DEBUG] User has no follower ")
		return []*descriptor.UserLight{}
	}
	resp := make([]*descriptor.UserLight, len(users), len(users))
	for i, u := range users {
		resp[i] = descriptor.NewUserLight(&u, u.IsFollowed(asker), false)
	}
	return resp
}
