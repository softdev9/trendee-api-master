package repos

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/data"
	"log"
	"os"
	"time"
)

func initRepos() (SelfieRepo, UserRepo, ArticleRepo) {
	dbHost := os.Getenv("MONGODB_TEST")
	DB_NAME = "trendee_test"
	dbSession, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal("Unable to connect to the database " + dbHost)
	}
	return &SelfieMGORepo{DBConn: dbSession.Copy()}, NewMongoUserRepo(dbSession.Copy()),
		NewArticleRepo(dbSession.Copy())
}

func cleanUp() {
	dbHost := os.Getenv("MONGODB_TEST")
	DB_NAME = "trendee_test"
	dbSession, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal("Unable to connect to the database " + dbHost)
	}
	dbSession.DB(DB_NAME).C(usersC).RemoveAll(nil)
	dbSession.DB(DB_NAME).C(selfieC).RemoveAll(nil)
	dbSession.DB(DB_NAME).C(votesC).RemoveAll(nil)
	dbSession.DB(DB_NAME).C(articleC).RemoveAll(nil)
}

func CreateTestSelfie(selfieRepo SelfieRepo, u *data.User) *data.Selfie {
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
			Category: "trousers",
			Brand:    "h&m",
			Color:    "red",
			Position: data.Location{
				PosX: 0.24,
				PosY: 0.23,
			},
		},
	}
	s, err := data.NewSelfie(u.ID, "Adding a test comment", tags, img)
	if err != nil {
		log.Fatal("Could not create selfies", err.Error())
	}
	selfieRepo.SaveSelfie(s)
	return s
}

func CreateTestUserAndSave(userRepo UserRepo) *data.User {
	u := &data.User{
		Username: "username",
		Birthday: time.Now().UTC(),
		UserType: "fashionista",
		Gender:   "male",
	}
	u.DefaultProfilePic()
	err := userRepo.SaveUser(u)
	if err != nil {
		log.Fatal("Could not create a test user")
	}
	return u
}
