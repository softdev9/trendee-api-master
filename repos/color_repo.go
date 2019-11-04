package repos

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
)

type ColorRepoMGO struct {
	DBConn *mgo.Session
}

type ColorRecord struct {
	HexCode       string `bson:"hexcode"`
	Name          string `bson:"_id"`
	TranslationFr string `bson:"name_fr"`
}

func NewColorRepo(conn *mgo.Session) *ColorRepoMGO {
	return &ColorRepoMGO{DBConn: conn}
}

func (repo *ColorRepoMGO) SaveColor(color *ColorRecord) error {
	return repo.DBConn.DB(DB_NAME).C(colorC).Insert(color)
}

func (repo *ColorRepoMGO) GetAllColors() ([]ColorRecord, error) {
	var results []ColorRecord
	err := repo.DBConn.DB(DB_NAME).C(colorC).Find(nil).All(&results)
	return results, err
}
