package repos

/*
import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
)

type CaterogryRepoMGO struct {
	DBConn *mgo.Session
}

type CategoryRecord struct {
	Name          string   `bson:"_id"`
	Gender        string   `bson:"gender"`
	TranslationFR string   `bson:"cat_fr"`
	SynonymEN     []string `bson:"synonym_en"`
	SynonymFR     []string `bson:"synonym_fr"`
}

type CateogryLister interface {
	GetAllCategories(string) ([]data.ArticleCategory, error)
}

func NewCatRepo(dbSession *mgo.Session) *CaterogryRepoMGO {
	return &CaterogryRepoMGO{DBConn: dbSession}
}

func (repo *CaterogryRepoMGO) repoC() *mgo.Collection {
	return repo.DBConn.DB(DB_NAME).C(categoryC)
}

func (repo *CaterogryRepoMGO) SaveCategoryRecord(cat *CategoryRecord) error {
	return repo.repoC().Insert(cat)
}

func (repo *CaterogryRepoMGO) GetCount() (int, error) {
	return repo.repoC().Count()
}

func (repo *CaterogryRepoMGO) GetAllCategories(lang string, gender string) ([]CategoryRecord, error) {
	sortCriteria := "_id"
	if lang == "fr" {
		sortCriteria = "translationFR"
	}
	pipeline := bson.D{
		{"gender", "mixed"},
	}
	if gender == "female" {
		pipeline = bson.D{
			{"$or", []interface{}{
				bson.D{{"gender", "mixed"}},
				bson.D{{"gender", "female"}},
			}},
		}
	}
	if gender == "male" {
		pipeline = bson.D{
			{"$or", []interface{}{
				bson.D{{"gender", "mixed"}},
				bson.D{{"gender", "male"}},
			}},
		}
	}
	var results []CategoryRecord
	err := repo.repoC().Find(pipeline).Sort(sortCriteria).All(&results)
	return results, err
}


func (repo *CaterogryRepoMGO) GetAllCategories(lang string) ([]data.ArticleCategory, error) {
	sortCriteria := "_id"
	if lang == "fr" {
		sortCriteria = "translationFR"
	}
	var results []data.ArticleCategory
	err := repo.repoC().Find(nil).Sort(sortCriteria).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}


func (repo *CaterogryRepoMGO) ClearAllCategories() error {
	_, err := repo.repoC().RemoveAll(nil)
	return err
}
*/
