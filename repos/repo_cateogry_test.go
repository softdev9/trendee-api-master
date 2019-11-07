package repos

/*
import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"os"
	"testing"
)

func TestCanAddCaterogry(t *testing.T) {
	repo, err := Setup()
	if err != nil {
		t.Error("Counl not set test")
	}
	defer repo.ClearAllCategories()
	checkEntitiyCountFor(t, repo, 0)
	catToSave := CategoryRecord{
		Name:          "dress",
		TranslationFR: "robes",
	}
	repo.SaveCategoryRecord(&catToSave)
	checkEntitiyCountFor(t, repo, 1)
}

func TestGetAllCaterogry(t *testing.T) {
	repo, err := Setup()
	if err != nil {
		t.Error("Counl not set test")
	}
	defer repo.ClearAllCategories()
	// Preparing data
	catsToSave := initDataToInsert()
	for _, cat := range catsToSave {
		repo.SaveCategoryRecord(&cat)
	}
	checkEntitiyCountFor(t, repo, 4)
	data, err := repo.GetAllCategories("fr", "male")
	if err != nil {
		t.Error("Unable to get the list categories")
		return
	} else {
		if len(data) != 2 {
			t.Error("The call should have returned 4 categories but got ", len(data))
			return
		}
		if data[0].Name != "shoes" {
			t.Error("Should be sorted by BasicName ", data)
			return
		}
	}
}

func TestGetAllCaterogryFR(t *testing.T) {
	repo, err := Setup()
	if err != nil {
		t.Error("Counl not set test")
	}
	// End of the TEST
	checkEntitiyCountFor(t, repo, 0)
	defer repo.ClearAllCategories()
	catsToSave := initDataToInsert()
	for _, cat := range catsToSave {
		repo.SaveCategoryRecord(&cat)
	}
	checkEntitiyCountFor(t, repo, 4)
	data, err := repo.GetAllCategories("fr", "female")
	if err != nil {
		t.Error("Unable to get the list categories")
		return
	} else {
		if len(data) != 4 {
			t.Error("The call should have returned 4 categories but got ", len(data))
			return
		}
	}

}

func Setup() (*CaterogryRepoMGO, error) {
	dbHost := os.Getenv("MONGODB_TEST")
	dbConn, err := mgo.Dial(dbHost)
	if err != nil {
		return nil, err
	}
	repo := &CaterogryRepoMGO{DBConn: dbConn}
	return repo, nil
}

func initDataToInsert() []CategoryRecord {
	catsToSave := []CategoryRecord{
		{
			Name:          "shoes",
			Gender:        "mixed",
			TranslationFR: "chaussure",
		},
		{
			Name:          "dress",
			Gender:        "female",
			TranslationFR: "robes",
		},
		{
			Name:          "skirt",
			Gender:        "female",
			TranslationFR: "jupes",
		},
		{
			Name:          "jacket",
			Gender:        "mixed",
			TranslationFR: "manteau",
		},
	}
	return catsToSave
}

func checkEntitiyCountFor(t *testing.T, repo *CaterogryRepoMGO, expected int) {
	if count, err := repo.GetCount(); err != nil {
		t.Error("Could not get the count of the collection " + err.Error())
	} else {
		if count != expected {
			t.Error("Expected ", expected, " got ", count)
		}
	}
}
*/
