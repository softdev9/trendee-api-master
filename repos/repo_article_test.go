package repos

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	//"os"
	"log"
	"strconv"
	"testing"
)

func TestGetArticleInList(t *testing.T) {
	DB_NAME = "trendee_test"
	defer cleanUp()
	_, _, artRepo := initRepos()
	arts, _ := artRepo.GetAll()
	if len(arts) != 0 {
		log.Println(arts)
		t.Errorf("There should be no artcile in the repo but there is %d", len(arts))
	}
	var bsonIds []bson.ObjectId
	for i := 0; i < 10; i++ {
		art := initArticleForShopTest(data.StoreCOM, strconv.Itoa(i))
		artRepo.SaveArticleToStore(art)
		bsonIds = append(bsonIds, art.Id)
	}
	if len(bsonIds) != 10 {
		t.Errorf("There should be %d in the repos but there is %d ", 10, len(bsonIds))
	}
	arts, err := artRepo.GetArticleInList(bsonIds)
	if err != nil {
		t.Errorf("Unable to retreived the ids ")
	}
	if len(arts) != 10 {
		t.Errorf("We were expecting %d articles to come back but we got %d ", 10, len(arts))
	}
}

func TestInsert(t *testing.T) {
	DB_NAME = "trendee_test"
	defer cleanUp()
	_, _, artRepo := initRepos()
	arts, _ := artRepo.GetAll()
	if len(arts) != 0 {
		log.Println(arts)
		t.Errorf("There should be no artcile in the repo but there is %d", len(arts))
	}
	// Connect to the db

	// Inserte the seconde article
	err := artRepo.SaveArticleToStore(initArticleForShopTest(data.StoreFR, ""))
	if err != nil {
		t.Errorf("Not able to insert the first article %s \n", err.Error())
	}
	art2 := initArticleForShopTest(data.StoreFR, "")
	err = artRepo.SaveArticleToStore(art2)
	if len(art2.Id.Hex()) == 0 {
		t.Errorf("No id has been attributed to updadted article")
	}
	if err != nil {
		t.Errorf("Not able to insert the second article", err.Error())
	}
	arts, _ = artRepo.GetAll()
	if len(arts) != 1 {
		log.Println(arts)
		t.Errorf("There should be one artcile in the repo but there is %d", len(arts))
	}
	art3 := initArticleForShopTest(data.StoreCOM, "")
	err = artRepo.SaveArticleToStore(art3)
	arts, _ = artRepo.GetAll()
	if len(arts) != 2 {
		log.Println(arts)
		t.Errorf("There should be two artcile in the repo but there is %d", len(arts))
	}

	artFound, _ := artRepo.GetArticleWithId(art2.Id.Hex())
	if artFound.NumberOfLike != 5 {
		t.Errorf("We should have 5 like on this article but we got %d", artFound.NumberOfLike)
	}
	cleanUp()
}

func TestRetrieveWithShop(t *testing.T) {
	cleanUp()
	_, _, artRepo := initRepos()
	art1 := initArticleForShopTest(data.StoreFR, "1")
	art2 := initArticleForShopTest(data.StoreCOM, "1")
	art3 := initArticleForShopTest(data.StoreCOM, "2")
	art4 := initArticleForShopTest(data.StoreCOM, "3")
	err := artRepo.SaveArticleToStore(art1)
	err = artRepo.SaveArticleToStore(art2)
	err = artRepo.SaveArticleToStore(art3)
	err = artRepo.SaveArticleToStore(art4)
	if err != nil {
		t.Errorf("Not able to insert test data", err.Error())
	}
	arts, err := artRepo.GetAll()
	if len(arts) != 4 {
		t.Errorf("We should have had %d articles but we got %d", 4, len(arts))
	}
	frLookup := []bson.ObjectId{art1.Id}
	frArts, err := artRepo.GetArticlesInStoreIds(data.StoreFR, frLookup)
	if len(frArts) != 1 {
		t.Errorf("We should have had %d articles in StoreFR but we got %d", 1, len(frArts))
	}
	enLookup := []bson.ObjectId{art2.Id, art3.Id, art4.Id}
	enArts, err := artRepo.GetArticlesInStoreIds(data.StoreCOM, enLookup)
	if len(enArts) != 3 {
		t.Errorf("We should have had %d articles in StoreFR but we got %d", 3, len(enArts))
	}
	cleanUp()
}

func initArticleForShopTest(shop string, uniqueId string) *data.Article {
	art := data.Article{
		Name: "Test article a",
		Source: data.ArticleSource{
			URL: "testsource.com",
			Id:  "2" + uniqueId,
		},
		Store:        shop,
		NumberOfLike: 5,
	}
	return &art
}
