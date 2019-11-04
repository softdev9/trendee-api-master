package repos

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"log"
	"time"
)

// Store the information about the change statistic for a like
type ArticleLikedChanged struct {
	ArticleId     bson.ObjectId
	IncNumLike    int
	IncNumDislike int
}

type ArticleRepoMGO struct {
	DB *mgo.Session
}

type ArticleLister interface {
	GetArticleInList(u *data.User) ([]data.Article, error)
}

type ArticleRepo interface {
	// Save the article given to the article collections
	SaveArticleToStore(article *data.Article) error
	// Save the array of article given as paramter to the article collections
	SaveMultiToStore(articles []interface{}) error
	// Get the article with the id given as paramaters
	GetArticleWithId(articleId string) (data.Article, error)
	// Get a list of articles with ids that are not contained in the list given as parameter
	GetArticleNotIds(articleIds []bson.ObjectId, gender string, category string) ([]data.Article, error)
	// Update the number of like on an article
	UpdateNumberOfLike(change ArticleLikedChanged) (*data.Article, error)
	// Get the artcicle in the store
	GetArticleInStore(idObjectToQuery string) ([]*data.Article, error)
	GetArticlesInStoreIds(storeName string, idObjectToQuery []bson.ObjectId) ([]data.Article, error)
	// Check if an article has been stored for a specific provider the id is the id given by the article provider
	//IsInShopFromProvider(ArticleC string, id string) bool
	// Get the list of articles that match the list of ids given as parameters
	GetArticleInList(idList []bson.ObjectId) ([]data.Article, error)
	// Return the number of Article in the collection
	CountArticleInStore(ArticleC string) (int, error)
	//
	GetAll() ([]data.Article, error)
	// Empty the store
	ClearStore(ArticleC string)
}

func NewArticleRepo(db *mgo.Session) ArticleRepo {
	return &ArticleRepoMGO{DB: db}
}

func (repo *ArticleRepoMGO) GetArticlesInStoreIds(storeName string, idObjectToQuery []bson.ObjectId) ([]data.Article, error) {
	findCriteria := bson.M{
		"$and": []interface{}{
			bson.M{"_id": bson.M{"$in": idObjectToQuery}},
			bson.M{"store": storeName},
		},
	}
	var arts []data.Article
	err := repo.DB.DB(DB_NAME).C(articleC).Find(findCriteria).All(&arts)
	if len(arts) == 0 {
		arts, err = repo.GetArticleInList(idObjectToQuery)

	}
	return arts, err
}

func (repo *ArticleRepoMGO) SaveArticleToStore(article *data.Article) error {
	findCriteria := bson.M{
		"source.id": article.Source.Id,
		"store":     article.Store,
	}
	count, err := repo.DB.DB(DB_NAME).C(articleC).Find(findCriteria).Count()
	if err != nil {
		log.Printf("[ERROR] Could not count the number of article", err.Error())
		return errors.New("Could not count the number of article")
	}
	if count == 0 {
		article.Id = bson.NewObjectId()
		article.RegistredOn = time.Now().UTC()
		err := repo.DB.DB(DB_NAME).C(articleC).Insert(article)
		if err != nil {
			log.Printf("[ERROR] Not able to insert article ", err.Error())
		}
		return err
	} else {
		art := &data.Article{}
		err := repo.DB.DB(DB_NAME).C(articleC).Find(findCriteria).One(art)
		if err != nil {
			log.Printf("[ERROR] Not able to find article to update ", err.Error())
		}
		change := bson.M{
			"$set": bson.M{
				"store":    article.Store,
				"price":    article.Price,
				"currency": article.Currency,
				"in_stock": article.InStock,
			},
		}
		article.Id = art.Id
		err = repo.DB.DB(DB_NAME).C(articleC).UpdateId(art.Id, change)
		if err != nil {
			log.Printf("[ERROR] Not able to update article", err.Error())
			return err
		}
	}
	return nil
}

func (repo *ArticleRepoMGO) GetAll() ([]data.Article, error) {
	var result []data.Article
	err := repo.DB.DB(DB_NAME).C(articleC).Find(nil).All(&result)
	return result, err
}

func (repo *ArticleRepoMGO) SaveMultiToStore(articles []interface{}) error {
	if len(articles) == 0 {
		return nil
	}
	start := time.Now()
	err := repo.DB.DB(DB_NAME).C(articleC).Insert(articles...)
	elapsed := time.Since(start)
	log.Printf("[PERF]  inserting in rep took  %s  for %d \n", elapsed, len(articles))
	return err
}

func (repo *ArticleRepoMGO) GetArticleWithId(articleId string) (data.Article, error) {
	var article data.Article
	artId, err := BuildIdFromString(articleId)
	err = repo.DB.DB(DB_NAME).C(articleC).Find(bson.M{"_id": artId}).One(&article)
	return article, err
}

func (repo *ArticleRepoMGO) GetArticleNotIds(excluded []bson.ObjectId, gender string, category string) ([]data.Article, error) {
	var result []data.Article
	findCriteria := bson.M{
		"$and": []interface{}{
			bson.M{"_id": bson.M{"$nin": excluded}},
			bson.M{"gender": gender},
			bson.M{"category": category},
		},
	}
	err := repo.DB.DB(DB_NAME).C(articleC).Find(findCriteria).All(&result)
	return result, err
}

func (repo *ArticleRepoMGO) UpdateNumberOfLike(change ArticleLikedChanged) (*data.Article, error) {
	var art data.Article
	update := bson.M{"$inc": bson.M{"numberOfLike": change.IncNumLike, "numberOfDislike": change.IncNumDislike}}
	toUpdate := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	_, err := repo.DB.DB(DB_NAME).C(articleC).Find(bson.M{"_id": change.ArticleId}).Apply(toUpdate, &art)
	return &art, err
}

func (repo *ArticleRepoMGO) GetArticleInStore(idObjectToQuery string) ([]*data.Article, error) {
	var articles []*data.Article
	queryCriteria := bson.M{}
	if cursor, err := BuildIdFromString(idObjectToQuery); err == nil {
		if cursor.Valid() {
			queryCriteria = bson.M{"_id": bson.M{"$lt": cursor}}
			log.Println("quering based on ", cursor.Hex(), "criteria ", queryCriteria)
		}
	}
	err := repo.DB.DB(DB_NAME).C(articleC).Find(queryCriteria).Sort("-_id").All(&articles)
	log.Println("criteria ", queryCriteria)
	if err != nil {
		log.Println("Error", err.Error())
	}
	return articles, err
}

func (repo *ArticleRepoMGO) GetArticleInList(idList []bson.ObjectId) ([]data.Article, error) {
	var articles []data.Article
	findCriteria := bson.M{"_id": bson.M{"$in": idList}}
	if err := repo.DB.DB(DB_NAME).C(
		articleC).Find(findCriteria).All(&articles); err != nil {
		log.Println(" [ERROR] ArticleRepoMGO_147:= ", err.Error())
		return nil, err
	}
	return articles, nil
}

func (repo *ArticleRepoMGO) CountArticleInStore(ArticleC string) (int, error) {
	return repo.DB.DB(DB_NAME).C(articleC).Find(bson.M{}).Count()
}

func (repo *ArticleRepoMGO) ClearStore(ArticleC string) {
	repo.DB.DB(DB_NAME).C(articleC).RemoveAll(bson.M{})
}
