package vote

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"net/http"
	//"log"
	"errors"
)

type IncomingVote struct {
	Selfie string          `schema:"selfie"`
	Vote   repos.VoteValue `schema:"vote"`
}

type OutVote struct {
	v0.Response
	Selfie *descriptor.SelfieDescriptor `json:"selfie"`
}

func (inVote *IncomingVote) Ok() error {
	if len(inVote.Selfie) < 1 {
		return errors.New("no selfie id sent")
	}
	if !repos.IsValidBsonObjectId(inVote.Selfie) {
		return errors.New("invalid selfie id sent")
	}
	return inVote.Vote.Ok()
}

func PostVote(rw http.ResponseWriter, req *http.Request) {
	// Decode the data coming from the form
	inVote := IncomingVote{Vote: -1}
	resp := OutVote{}
	if err := handlers.DecodeForm(req, &inVote); err != nil {
		resp.Response.Error = handlers.ErrorDescriptor{
			ErrorCode: 101,
			Error:     " Unable to decode incoming vote. Reason " + err.Error(),
		}
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	selfieRepo := context.Get(req, repos.SelfieR).(repos.SelfieRepo)
	userRepo := context.Get(req, repos.UserR).(repos.UserRepo)
	fromUser := context.Get(req, "user").(*data.User)
	idSelfie, _ := repos.BuildIdFromString(inVote.Selfie)
	//Vote(user bson.ObjectId, selfie bson.ObjectId, voteValue VoteValue) error
	err := selfieRepo.Vote(fromUser.ID, idSelfie, inVote.Vote)
	if err != nil {
		resp.Response.Error = handlers.ErrorDescriptor{
			ErrorCode: 101,
			Error:     "Impossible to vote on the selfie. " + err.Error(),
		}
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	selfieR, _ := selfieRepo.GetSelfieById(idSelfie)
	// (selfie repos.SelfieRecord, requesting *data.User, userDatasource repos.UserIdGetter, vote repos.VoteValue)
	selfieD := descriptor.DescribeSelfie(selfieR, fromUser, userRepo, inVote.Vote, nil)
	resp.Selfie = &selfieD
	handlers.Respond(rw, req, http.StatusOK, &resp)
	//log.Printf("Vote taken into account \n")
}
