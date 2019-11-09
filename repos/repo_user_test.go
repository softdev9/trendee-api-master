package repos

import (
	"testing"
	"time"
)

func TestFollow(t *testing.T) {
	_, userRepo, _ := initRepos()
	u1 := CreateTestUserAndSave(userRepo)
	u2 := CreateTestUserAndSave(userRepo)
	from, to, err := userRepo.EstablishFollow(u1.ID, u2.ID)
	if err != nil {
		t.Error("Could not establish the follow but it should have been")
	}
	if len(from.Following) != 1 {
		t.Error("The follow should have been established the from.Following should be 1")
	}
	if len(to.Following) != 0 {
		t.Error("The follow should have been established the to.Following should be 0")
	}
	if len(to.FollowedBy) != 1 {
		t.Error("The follow should have been established the to.FollowedBy should be 1 but is ", to.FollowedBy)
	}
	if len(from.FollowedBy) != 0 {
		t.Error("The follow should have been established the from.FollowedBy should be 0")
	}
	_, _, err = userRepo.EstablishFollow(u1.ID, u1.ID)
	if err == nil {
		t.Error("The user should not be allowed to follow himself")
	}
	cleanUp()
}

func TestInfoUpdater(t *testing.T) {
	DB_NAME = "trendee_test"
	_, userRepo, _ := initRepos()
	updateMap := map[string]string{
		"description": "a small description of the user",
		"username":    "kevinlegoff",
		"birthday":    "11-08-1988",
		"gender":      "female",
		"website":     "http://www.google.fr",
		"city":        "testcity",
		"country":     "testcountry",
		"email":       "newemail@gmail.com",
	}
	u1 := CreateTestUserAndSave(userRepo)
	uUpdated, err := userRepo.UpdateUser(u1, updateMap)
	if err != nil {
		t.Errorf("Should not have thrown an error")
	}
	if uUpdated.Description != "a small description of the user" {
		t.Errorf("Description not matching update")
	}
	if uUpdated.Username != "kevinlegoff" {
		t.Errorf("Username not matching update")
	}
	expectedDate := time.Date(1988, time.August, 11, 0, 0, 0, 0, time.UTC)
	if uUpdated.Birthday.UTC() != expectedDate {
		t.Errorf("Birthday %s not not matching update %s ", uUpdated.Birthday.UTC(), expectedDate)
	}
	if uUpdated.Website != "http://www.google.fr" {
		t.Errorf("Website %s  %s ", uUpdated.Website, "http://www.google.fr")
	}
	if uUpdated.Gender != "female" {
		t.Errorf("Gender %s  %s ", uUpdated.Gender, "female")
	}
	if uUpdated.City != "testcity" {
		t.Errorf("City %s  %s ", uUpdated.City, "testcity")
	}
	if uUpdated.Country != "testcountry" {
		t.Errorf("Country %s  %s ", uUpdated.City, "testcountry")
	}
	if uUpdated.EmailChange.Email != "newemail@gmail.com" {
		t.Errorf("The email has not been marked for change")
	}
	if len(uUpdated.EmailChange.Random) == 0 {
		t.Errorf("The random for email update has not been generated")
	}
	cleanUp()
}

/*
func TestCreateNewUserAndRetreiveByEmail(t *testing.T) {
	dbHost := os.Getenv("MONGODB_TEST")
	dbConn, err := mgo.Dial(dbHost)
	if err != nil {
		t.Error("Not able to connect to the db")
	}
	defer dbConn.Close()

	ClearCollection(dbConn)
	repo := MongoUserRepo{DBConn: dbConn}
	repo.FindByEmail("kevin")
	userToSave := data.User{Email: "kev.legoff@gmail.com"}
	repo.SaveUser(&userToSave)
	userRetreived := repo.FindByEmail("kev.legoff@gmail.com")
	if !userRetreived {
		t.Error("The user should have been retreived")
	}
	ClearCollection(dbConn)
}

func TestFindUserNameAndPassword(t *testing.T) {
	dbHost := os.Getenv("MONGODB_TEST")
	dbConn, err := mgo.Dial(dbHost)
	if err != nil {
		t.Error("Not able to connect to the db")
	}
	defer dbConn.Close()
	ClearCollection(dbConn)
	repo := MongoUserRepo{DBConn: dbConn}
	toInsert := CreateFakeUser("kev.legoff@gmail.com", "test1234", "untestsalt")
	repo.SaveUser(&toInsert)

	cases := []struct {
		Password      string
		Email         string
		UserExpected  data.User
		ErrorExpexted error
	}{
		{
			Email:         "kev.legoff@gmail.com",
			Password:      "test1234",
			ErrorExpexted: nil,
			UserExpected:  toInsert,
		}, {
			Email:         "kev@gmail.com",
			Password:      "test1234",
			ErrorExpexted: errors.New("User Not Found"),
		}, {
			Email:         "kev.legoff@gmail.com",
			Password:      "test",
			ErrorExpexted: errors.New("Password Invalid"),
		},
	}
	for _, test := range cases {
		u, _, err := repo.CheckEmailPassword(test.Email, test.Password)
		if err != nil {
			if test.ErrorExpexted != nil {
				if err.Error() != test.ErrorExpexted.Error() {
					t.Errorf("Got error %s expected %s", err.Error(), test.ErrorExpexted.Error())
					return
				}
			} else {
				t.Errorf("No error was expected for test input %s ", test)
			}
		}
		if u != nil {
			if u.Email != test.Email {
				t.Error("User was not retreived correctly ", *u)
			}
		}
	}

}

// Test the like article function
//
func TestLikeArticle(t *testing.T) {
	dbHost := os.Getenv("MONGODB_TEST")
	dbConn, err := mgo.Dial(dbHost)
	if err != nil {
		t.Error("Not able to connect to the db")
	}
	defer dbConn.Close()
	// Create a test user
	u := CreateFakeUser("kev.legoff@gmail.com", "test1234", "mail")
	repo := MongoUserRepo{DBConn: dbConn}
	repo.SaveUser(&u)
	// Add the article to the liked of the user
	artId := bson.NewObjectId()
	// We add the article to the user dislike article this is the worst case scenario
	// Whare a user goes from dislike to like
	u.ArticleDisliked = append(u.ArticleDisliked, artId)
	change, err := repo.MarkArticle(&u, artId, "like")
	if err != nil {
		t.Error("The like should have been established")
	}
	ok := false
	for _, idA := range u.ArticleLiked {
		if idA == artId {
			ok = true
		}
	}

	if !ok {
		fmt.Println(u.ArticleLiked)
		t.Error("The article was not added in the liked list of the user")
	}

	if change.IncNumLike != 1 {
		t.Error("the article change for like should be at 1 it is at ", change.IncNumLike)
	}

	if len(u.ArticleDisliked) != 0 {
		t.Error("the article should not be there anymore")
	}
	change, err = repo.MarkArticle(&u, artId, "like")
	if err == nil {
		t.Error("Should have returned an error as the article is already liked")
	}
	// User change is mind and dislike the article
	change, err = repo.MarkArticle(&u, artId, "dislike")
	ok = false
	for _, idA := range u.ArticleDisliked {
		if idA == artId {
			ok = true
		}
	}
	if !ok {
		fmt.Println(u.ArticleDisliked)
		t.Error("The article was not added in the disliked list of the user")
	}
	if change.IncNumLike != -1 {
		t.Error("the article change for like should be at -1 it is at ", change.IncNumLike)
	}
	if len(u.ArticleDisliked) != 1 {
		t.Error("the article should not be there anymore")
	}
	ClearCollection(dbConn)
}

func TestLikeUnlike(t *testing.T) {
	// Connect to the db
	dbConn := connectDB(t)
	defer dbConn.Close()
	// Create a test user
	u := CreateFakeUser("kev.legoff@gmail.com", "test1234", "mail")
	repo := MongoUserRepo{DBConn: dbConn}
	repo.SaveUser(&u)
	artLiked := bson.NewObjectId()
	change, err := repo.MarkArticle(&u, artLiked, "like")
	if err != nil {
		t.Error("Not able to make the test", err.Error())
	}
	if change.ArticleId != artLiked {
		t.Error("No consistency in the id given and id returned")
	}
	if change.IncNumDislike != 0 {
		t.Error("We should not change the number of dislike on like operation")
	}
	if change.IncNumLike != 1 {
		t.Error("We are liking so it should be 1")
	}
	if !u.HasLiked(artLiked) {
		t.Error("The article should be in the liked list")
	}
	if u.HasDisliked(artLiked) {
		t.Error("The article should NOT be in the liked list")
	}

	// Unlike the article freshlye liked
	change, err = repo.MarkArticle(&u, artLiked, "unlike")
	if err != nil {
		t.Error("Not able to make the test", err.Error())
	}
	if change.ArticleId != artLiked {
		t.Error("No consistency in the id given and id returned")
	}
	if change.IncNumDislike != 0 {
		t.Error("We should not change the number of dislike on like operation")
	}
	if change.IncNumLike != -1 {
		t.Error("We are liking so it should be 1")
	}
	if u.HasLiked(artLiked) {
		t.Error("The article should NOT be in the liked list")
	}
	if u.HasDisliked(artLiked) {
		t.Error("The article should NOT be in the liked list")
	}
	change, err = repo.MarkArticle(&u, artLiked, "like")
	change, err = repo.MarkArticle(&u, artLiked, "dislike")
	if err != nil {
		t.Error("Not able to make the test", err.Error())
	}
	if change.ArticleId != artLiked {
		t.Error("No consistency in the id given and id returned")
	}
	if change.IncNumDislike != 1 {
		t.Error("We should not change the number of dislike on like operation")
	}
	if change.IncNumLike != -1 {
		t.Error("We are liking so it should be 1")
	}
	if u.HasLiked(artLiked) {
		t.Error("The article should NOT be in the liked list")
	}
	if !u.HasDisliked(artLiked) {
		t.Error("The article should NOT be in the liked list")
	}

	change, err = repo.MarkArticle(&u, artLiked, "undislike")
	if err != nil {
		t.Error("Not able to make the test", err.Error())
	}
	if change.ArticleId != artLiked {
		t.Error("No consistency in the id given and id returned")
	}
	if change.IncNumDislike != -1 {
		t.Error("We should not change the number of dislike on like operation")
	}
	if change.IncNumLike != 0 {
		t.Error("We are liking so it should be 1")
	}
	if u.HasLiked(artLiked) {
		t.Error("The article should NOT be in the liked list")
	}
	if u.HasDisliked(artLiked) {
		t.Error("The article should NOT be in the liked list")
	}

}

func connectDB(t *testing.T) *mgo.Session {
	dbHost := os.Getenv("MONGODB_TEST")
	dbConn, err := mgo.Dial(dbHost)
	if err != nil {
		t.Error("Not able to connect to the db")
	}
	return dbConn.Copy()
}

func CreateFakeUser(email string, password string, salt string) data.User {
	return data.NewUser(email, password, salt)

}

func ClearCollection(db *mgo.Session) {
	//db.DB(dbN).C(usersC).RemoveAll(bson.M{})
}
*/
