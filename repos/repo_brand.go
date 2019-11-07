package repos

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"log"
)

type BrandRecord struct {
	Name string            `bson:"_id"`
	Logo map[string]string `bson:"logo"`
}

type BrandRepoMGO struct {
	DBConn *mgo.Session
}

type BrandRepo interface {
	GetAllBrands() ([]data.Brand, error)
	GetCount() (int, error)
	SaveBrand(*BrandRecord) error
	ClearAllBrands() error
	HasBrand(string) bool
	FindByName(string) (BrandRecord, error)
}

type BrandNameFinder interface {
	FindByName(string) (BrandRecord, error)
}

func NewBrandRepo(s *mgo.Session) BrandRepo {
	return &BrandRepoMGO{DBConn: s}
}

func (repo *BrandRepoMGO) HasBrand(name string) bool {
	c, err := repo.repoC().Find(bson.M{"_id": name}).Count()
	if err != nil {
		return false
	}
	return c > 0
}

func (repo *BrandRepoMGO) repoC() *mgo.Collection {
	return repo.DBConn.DB(DB_NAME).C(brandC)
}

func (repo *BrandRepoMGO) GetCount() (int, error) {
	return repo.repoC().Count()
}

func (repo *BrandRepoMGO) SaveBrand(brand *BrandRecord) error {
	if brand.Name == "" || brand.Name == " " {
		return nil
	}
	if _, err := repo.repoC().Upsert(bson.M{"_id": brand.Name}, brand); err != nil {
		log.Printf("[ERROR] [DB] BRAND %s", err.Error())
		return err
	}
	log.Printf("[INFO][DB][BRAND]%s brand saved ", brand.Name)
	return nil
}

func (repo *BrandRepoMGO) ClearAllBrands() error {
	_, err := repo.repoC().RemoveAll(nil)
	return err
}

func (repo *BrandRepoMGO) FindByName(name string) (BrandRecord, error) {
	var b BrandRecord
	err := repo.repoC().Find(bson.M{"_id": name}).One(&b)
	return b, err
}

func (repo *BrandRepoMGO) GetAllBrands() ([]data.Brand, error) {
	var results []data.Brand
	err := repo.repoC().Find(nil).Sort("_id").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
