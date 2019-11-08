package repos

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"testing"
)

func TestVote(t *testing.T) {
	selfieRepo, userRepo, _ := initRepos()
	u := CreateTestUserAndSave(userRepo)
	u2 := CreateTestUserAndSave(userRepo)
	s := CreateTestSelfie(selfieRepo, u)
	err := selfieRepo.Vote(u2.ID, s.ID, VoteValue(OUT))
	if err == nil {
		t.Log("We are good no  error shas been returned")
	} else {
		t.Error("Error when voting ", err.Error())
	}
	selfieUpdated, _ := selfieRepo.GetSelfieById(s.ID)
	//t.Logf("Selfie %#v \n", selfieUpdated)
	if selfieUpdated.NumberOfMinus != 1 {
		t.Error("The number of minus should be 1")
	}
	userUpdated, _ := userRepo.GetUserById(u2.ID)
	var idFound bool = false
	for _, id := range userUpdated.SelfieVoted {
		if id == s.ID {
			idFound = true
		}

	}
	if idFound == false {
		t.Error("The selfie id was not retreieved in the user voted list")
	}

	// Check the vote
	voteRecord, err := selfieRepo.FindVoteFor(s.ID, u2.ID)
	if err != nil {
		t.Error("Error while trying to recover the vote ", err.Error())
	}
	if voteRecord.Value != OUT {
		t.Error("The vote value should be out as given in the call to vote")
	}
	t.Logf("TEST OF VOTE COMPLETED")
	cleanUp()
}

func TestLinkToArticles(t *testing.T) {
	selfieRepo, userRepo, _ := initRepos()
	u := CreateTestUserAndSave(userRepo)
	s := CreateTestSelfie(selfieRepo, u)
	fakeArticleId := []bson.ObjectId{
		bson.NewObjectId(),
		bson.NewObjectId(),
	}
	selfieRepo.Similars(s.ID, fakeArticleId)
	retreived, _ := selfieRepo.GetSelfieById(s.ID)
	for i, artId := range fakeArticleId {
		if artId != retreived.RelatedArticle[i] {
			t.Errorf("Not the same ids %s  -> %s", artId, retreived.RelatedArticle[i])
		}
	}
	cleanUp()
}
