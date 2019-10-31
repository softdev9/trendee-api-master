package home

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	//"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"encoding/json"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestSelfPosted(t *testing.T) {
	repos.DB_NAME = "trendee_test_bed"
	// Create a request
	req, err := http.NewRequest("GET", "http://", nil)
	if err != nil {
		t.Error(err.Error())
	}
	// Create a user with a token
	testServerForReq(req)
	createFakeUser(t, req)
	// Create a selfie
	createAndSaveSelfie(req)
	// Load the home
	rw := httptest.NewRecorder()
	GetHome(rw, req)
	// Check there is not the same selfie
	log.Println(rw)
	// Decode the response
	hompeResp := &HomeResponse{}
	json.NewDecoder(rw.Body).Decode(hompeResp)
	log.Printf("Response decoded %v ", hompeResp)
	// Check the
	if len(hompeResp.Page.Selfie) != 0 {
		t.Error("The response contains the selfie posted by the user it should not")
	}
	// Clean up the test
	cleanUp()
}

func testServerForReq(r *http.Request) {
	dbHost := os.Getenv("MONGODB_TEST_BED")
	dbConn, _ := mgo.Dial(dbHost)
	dbCopy := dbConn.Copy()
	context.Set(r, repos.TokenR, repos.NewTokenRepo(dbCopy.Copy()))
	context.Set(r, repos.SelfieR, repos.NewSelfieRepoMGO(dbCopy.Copy()))
	context.Set(r, repos.BrandR, repos.NewBrandRepo(dbCopy.Copy()))
	context.Set(r, repos.UserR, repos.NewMongoUserRepo(dbCopy.Copy()))
	context.Set(r, repos.ColorsR, repos.NewColorRepo(dbCopy.Copy()))
	context.Set(r, repos.ArticleR, repos.NewArticleRepo(dbCopy.Copy()))
}

func createFakeUser(t *testing.T, req *http.Request) {
	userRepo := context.Get(req, repos.UserR).(*repos.MongoUserRepo)
	tokenRepo := context.Get(req, repos.TokenR).(repos.TokenRepo)
	// Create a test user
	u := &data.User{
		Username: "Username With Spaces",
		Birthday: time.Now().UTC(),
		UserType: "fashionista",
		Gender:   "male",
	}
	u.DefaultProfilePic()
	if err := userRepo.SaveUser(u); err != nil {
		log.Fatalf("Not able to save the user %s", err.Error())
	}
	// Create the token
	token := data.NewAuthToken()
	// Save the token to the db
	if err := tokenRepo.SaveAuthToken(token, *u); err != nil {
		t.Error("Not able to save the user ", err.Error())
	}
	context.Set(req, "user", u)
}

func createAndSaveSelfie(req *http.Request) {
	u := context.Get(req, "user").(*data.User)
	selfieRepo := context.Get(req, repos.SelfieR).(repos.SelfieRepo)
	img := map[string]string{
		"xlarge": "https://www.selfievote.com/assets/selfievote/og.jpg",
		"large":  "https://www.selfievote.com/assets/selfievote/og.jpg",
		"medium": "https://www.selfievote.com/assets/selfievote/og.jpg",
		"small":  "https://www.selfievote.com/assets/selfievote/og.jpg",
		"xsmall": "https://www.selfievote.com/assets/selfievote/og.jpg",
	}
	tags := []data.Tag{
		{
			Gender:   "woman",
			Category: "dresses",
			Brand:    "boots",
			Color:    "black",
			Position: data.Location{
				PosX: 0.24,
				PosY: 0.23,
			},
		},
		{
			Gender:   "men",
			Category: "mens-shirts",
			Brand:    "addidas",
			Color:    "black",
			Position: data.Location{
				PosX: 0.24,
				PosY: 0.70,
			},
		},
	}

	// We need to create the image onlydata.NewSelfie(u, categories, brands, colors, img)
	s, _ := data.NewSelfie(u.ID, "Adding a test comment", tags, img)
	selfieRepo.SaveSelfie(s)
}

func cleanUp() {
	log.Printf("[DEBUG] In clean up for dbname %s", repos.DB_NAME)
	dbHost := os.Getenv("MONGODB_TEST_BED")
	dbConn, _ := mgo.Dial(dbHost)
	dbConn.DB(repos.DB_NAME).C("users").Remove(nil)
	dbConn.DB(repos.DB_NAME).C("selfies").Remove(nil)
	dbConn.DB(repos.DB_NAME).C("tokens").Remove(nil)
}
