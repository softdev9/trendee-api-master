package repos

/*
import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/data"
	"os"
	"testing"
)

func TestEmptyCollection(t *testing.T) {
	dbRepo, err := SetupBrandTest()
	if err != nil {
		t.Error("fail to connect to the db")
	}
	count, err := dbRepo.GetCount()
	if err != nil {
		t.Error("Error in the count")
	}
	if count != 0 {
		t.Error("The db should be empty")
	}
}

func TestSaveAndRetreiveBrand(t *testing.T) {
	dbRepo, err := SetupBrandTest()
	if err != nil {
		t.Error("Could not connect to the database")
	}
	err = dbRepo.SaveBrand(data.NewBrand("test", 123, "h&m.png", "http://www.logo.com"))
	if err != nil {
		t.Error("Could not save the brand")
	}
	err = dbRepo.SaveBrand(data.NewBrand("test", 123, "h&m.png", "http://www.logo.com"))
	if err != nil {
		t.Error("Could not save the brand ", err.Error())
	}
	err = dbRepo.SaveBrand(data.NewBrand("test2", 123, "h&m.png", "http://www.logo.com"))
	if err != nil {
		t.Error("Could not save the brand ", err.Error())
	}
	// Check we can retreive the brand 1 and brand 2
	brands, err := dbRepo.GetAllBrands()
	if err != nil {
		t.Error("Could not retreive the brand ", err.Error())
	}
	if brands[0].Name != "test" {
		t.Error("Not able to find the brand")
	}
	if brands[1].Name != "test2" {
		t.Error("Not able to list the brand")
	}
	// Check we got 2 brands here
	count, err := dbRepo.GetCount()
	if count > 2 {
		t.Error("We should have only one brand for the name ")
	}
	dbRepo.ClearAllBrands()
}

func SetupBrandTest() (*BrandRepoMGO, error) {
	dbHost := os.Getenv("MONGODB_TEST")
	dbConn, err := mgo.Dial(dbHost)
	if err != nil {
		return nil, err
	}
	repo := &BrandRepoMGO{DBConn: dbConn}
	repo.ClearAllBrands()
	return repo, nil
}
*/
