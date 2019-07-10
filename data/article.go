package data

import (
	"github.com/spiderdev86/trendee-api/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"time"
)

// Collection for the article avaiable in france
const StoreFR = "store_fr"
const StoreCOM = "store_com"

type ArticleSource struct {
	URL string
	Id  string
}

type Article struct {
	Id              bson.ObjectId `bson:"_id"`
	Brand           string        `bson:"brand"`
	Keywords        []string      
	Name            string
	Price           float32 `bson:"price"`
	Category        string
	Currency        string `bson:"currency"`
	InStock         bool   `bson:"in_stock"` // See how we can update it in a nice manner
	Gender          string `bson:"gender"`
	Source          ArticleSource
	Image           map[string]string `bson:"image"`
	Description     string            `bson:"description"`
	Retailer        string
	ClickUrl        string
	NumberOfShare   int       `bson:"numberOfShare"`
	NumberOfLike    int       `bson:"numberOfLike"`
	NumberOfDislike int       `bson:"numberOfDislike"`
	RegistredOn     time.Time `bson:"registred_on"`
	SortCoef        float64
	Store           string `bson:"store"`
}

func (art Article) Attributes() []string {
	return art.Keywords
}

func (art *Article) SetCoef(c float64) {
	art.SortCoef = c
}

func (art *Article) Coef() float64 {
	return art.SortCoef
}
