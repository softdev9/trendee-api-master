package follow

/*
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	//"errors"
)

var userCount int = 0

func TestFollowHandler(t *testing.T) {
	uRepo := &URepo{
		Users: make(map[string]data.User),
	}
	to := createUser(bson.NewObjectId())
	uRepo.Users[to.ID.Hex()] = *to
	form := url.Values{}
	form.Add("to", to.ID.Hex())
	// Build the request
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	context.Set(req, repos.UserR, uRepo)
	// Record the resul
	rw := httptest.NewRecorder()
	// Post the follow
	PostFollow(rw, req)
	// Decode the response
	if rw.Code != http.StatusOK {
		t.Errorf("Error code expected %d but got %d", http.StatusOK, rw.Code)
	}
	decoder := json.NewDecoder(rw.Body)
	resp := FollowedResponse{}
	if err := decoder.Decode(&resp); err != nil {
		t.Errorf("Could not decode the answer ", err)
	}
	if resp.Error.ErrorCode != 0 {
		t.Errorf("L35 Error code expected %d but got %d", 0, resp.Error.ErrorCode)
	}
}

func createUser(id bson.ObjectId) *data.User {
	userCount = userCount + 1
	return &data.User{
		ID:        id,
		Lastname:  fmt.Sprintf("TestLN%d", userCount),
		Firstname: fmt.Sprintf("TestFN%d", userCount),
	}
}

type URepo struct {
	Users map[string]data.User
}

func (repo *URepo) GetUserById(id bson.ObjectId) (*data.User, error) {
	u := repo.Users[id.Hex()]
	return &u, nil
	//return nil, errors.New("not found")
}
*/
