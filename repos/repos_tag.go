package repos

import (
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
)

type TagRepoMgo struct {
	DBConn *mgo.Session
}

// Defines the entiry we store in the database
type TagRecord struct {
	Tag string `bson:"_id"`
}

type TagRepo interface {
	// Insert a tag if he is not present in the collection
	AddTag(tag string) error
	// Retrieve Tag List
	GetAllTags() ([]string, error)
}

func NewTagRepoMgo(dbConn *mgo.Session) *TagRepoMgo {
	return &TagRepoMgo{DBConn: dbConn}
}

func (repo *TagRepoMgo) GetAllTags() ([]string, error) {
	var results []TagRecord
	err := repo.DBConn.DB(DB_NAME).C(tagC).Find(bson.M{}).Sort("_id").All(&results)
	if err != nil {
		return nil, err
	}
	count, err := repo.Size()
	if err != nil {
		return nil, err
	}
	ret := make([]string, count)
	for i, rec := range results {
		ret[i] = rec.Tag
	}
	return ret, nil
}

func (repo *TagRepoMgo) AddTag(tag string) error {
	inserted := TagRecord{
		Tag: tag,
	}
	count, err := repo.DBConn.DB(DB_NAME).C(tagC).Find(bson.M{"_id": tag}).Count()
	if err != nil {
		fmt.Println("Error : ", err.Error())
		return err
	}
	if count == 0 {
		// We use upsert to avoid to have problem when two tags are inserterd
		err := repo.DBConn.DB(DB_NAME).C(tagC).Insert(inserted)
		if err != nil {
			fmt.Println("Error : ", err.Error())
		}
		return err
	}
	return nil
}

func (repo *TagRepoMgo) Size() (int, error) {
	return repo.DBConn.DB(DB_NAME).C(tagC).Find(bson.M{}).Count()
}
