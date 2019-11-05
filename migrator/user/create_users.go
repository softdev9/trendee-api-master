package main

import (
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Creating test users and test tokens")
	dbUrl := os.Getenv("MONGODB_TEST")
	//dbUrl := os.Getenv("MONGODB_PROD")
	dbSession, err := mgo.Dial(dbUrl)
	if err != nil {
		log.Fatal("Not able to connect to the db ", dbUrl)
	}
	defer dbSession.Close()
	userRepo := repos.NewMongoUserRepo(dbSession)
	tokenRepo := repos.NewTokenRepo(dbSession)
	//selfieRepo := repos.NewSelfieRepoMGO(dbSession)

	for i := 0; i < 15; i++ {
		gender := "female"
		if i%2 == 0 {
			gender = "male"
		}
		// Create a test user
		u := &data.User{
			Username: "Username With Spaces And super long and need to be kfoofofokfl flflfl.... " + strconv.Itoa(i),
			Birthday: time.Now().UTC(),
			UserType: "fashionista",
			Gender:   gender,
		}
		u.DefaultProfilePic()
		// Save the user and create a token
		if err := userRepo.SaveUser(u); err != nil {
			log.Fatal("Not able to save the user")
		}
		// Create the token
		token := data.NewAuthToken()
		// Save the token to the db
		if err := tokenRepo.SaveAuthToken(token, *u); err != nil {
			log.Fatal("Not able to save the user")
		}
		// for the first 5 user create 3 selfies
		if i < 5 {
			//createThreeSelfies(u, selfieRepo)
		}
		fmt.Print("user created \n\t", u, "\nToken\n\t", token, "\n\n")
	}
}

func createThreeSelfies(u *data.User, selfieRepo *repos.SelfieMGORepo) {
	for i := 0; i < 3; i++ {
		img := map[string]string{
			"xlarge": "https://www.selfievote.com/assets/selfievote/og.jpg",
			"large":  "https://www.selfievote.com/assets/selfievote/og.jpg",
			"medium": "https://www.selfievote.com/assets/selfievote/og.jpg",
			"small":  "https://www.selfievote.com/assets/selfievote/og.jpg",
			"xsmall": "https://www.selfievote.com/assets/selfievote/og.jpg",
		}
		tags := []data.Tag{
			{
				Gender:   "female",
				Category: "dresses",
				Brand:    "boots",
				Color:    "black",
				Position: data.Location{
					PosX: 0.24,
					PosY: 0.23,
				},
			},
			{
				Gender:   "male",
				Category: "mens-shirts",
				Brand:    "addidas",
				Color:    "black",
				Position: data.Location{
					PosX: 0.24,
					PosY: 0.70,
				},
			},
		}
		log.Println("ID from user ", u.ID.Hex())
		// We need to create the image onlydata.NewSelfie(u, categories, brands, colors, img)
		s, _ := data.NewSelfie(u.ID, "Adding a test comment", tags, img)
		selfieRepo.SaveSelfie(s)
	}
}
