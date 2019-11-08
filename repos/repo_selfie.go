package repos

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/txn"
	"github.com/softdev9/trendee-api-master/data"
	"log"
	"time"
)

type SelfieMGORepo struct {
	DBConn *mgo.Session
}

type VoteValue int

const (
	OUT VoteValue = 1
	OK  VoteValue = 2
	TOP VoteValue = 3
)

func (v VoteValue) Ok() error {
	if v != OUT && v != OK && v != TOP {
		return errors.New("Invalid vote value")
	}
	return nil
}

type SelfieRepo interface {
	SaveSelfie(s *data.Selfie)
	GetSelfiePublishedBy(uId bson.ObjectId) ([]SelfieRecord, error)
	GetSelfieById(selfieId bson.ObjectId) (SelfieRecord, error)
	GetSelfieFromOther(notPublishedById bson.ObjectId,
		cursor *bson.ObjectId, listLen int, page int) ([]SelfieRecord, error)
	Vote(user bson.ObjectId, selfie bson.ObjectId, voteValue VoteValue) error
	FindVoteFor(selfie bson.ObjectId, user bson.ObjectId) (VoteRecord, error)
	Similars(selfie bson.ObjectId, articles []bson.ObjectId) error
}

type SelfieRecord struct {
	ID               bson.ObjectId     `bson:"_id"`
	Author           bson.ObjectId     `bson:"author"`
	Image            map[string]string `bson:"images"`
	Tags             []data.Tag        `bson:"tags"`
	Gender           string            `bson:"gender"`
	Comment          string            `bson:"comment"`
	Keywords         []string          `bson:"keywords"`
	NumberOfComments int               `bson:"number_of_comments"`
	NumberOfPlus     int               `bson:"number_of_plus"`
	NumberOfNeutral  int               `bson:"number_of_neutral"`
	NumberOfMinus    int               `bson:"number_of_minus"`
	CreatedOn        time.Time         `bson:"created_on"`
	RelatedArticle   []bson.ObjectId   `bson:"similars"`
}

type VoteRecord struct {
	ID        bson.ObjectId `bson:"_id"`
	Voter     bson.ObjectId `bson:"voter"`
	Selfie    bson.ObjectId `bson:"selfie"`
	Value     VoteValue     `bson:"vote_value"`
	CreatedOn time.Time     `bson:"created_on"`
}

func NewSelfieRepoMGO(db *mgo.Session) *SelfieMGORepo {
	return &SelfieMGORepo{DBConn: db}
}

func (repo *SelfieMGORepo) GetSelfieById(selfieId bson.ObjectId) (SelfieRecord, error) {
	var selfie SelfieRecord
	err := repo.DBConn.DB(DB_NAME).C(selfieC).FindId(selfieId).One(&selfie)
	return selfie, err
}

func (repo *SelfieMGORepo) Similars(selfie bson.ObjectId, articles []bson.ObjectId) error {
	change := bson.M{
		"$set": bson.M{
			"similars": articles,
		},
	}
	err := repo.DBConn.DB(DB_NAME).C(selfieC).UpdateId(selfie, change)
	return err
}

func (repo *SelfieMGORepo) SaveSelfie(s *data.Selfie) {
	// Increment the number of selfie for the author
	change := bson.M{"$inc": bson.M{"selfie-posted": 1, "suggestion-index": 1}}
	if err := repo.DBConn.DB(DB_NAME).C(usersC).UpdateId(s.Author, change); err != nil {
		log.Println("Error updating user with id ", err.Error())
	}
	repo.DBConn.DB(DB_NAME).C(selfieC).Insert(fromSelfie(s))
}

func (repo *SelfieMGORepo) GetSelfiePublishedBy(uId bson.ObjectId) ([]SelfieRecord, error) {
	findCriteria := bson.M{"author": uId}
	var result []SelfieRecord
	err := repo.DBConn.DB(DB_NAME).C(selfieC).Find(findCriteria).Sort("-created_on").All(&result)
	if err != nil {
		log.Println("Error while retreiving selfies from author ", err.Error())
	}
	return result, err
}

func (repo *SelfieMGORepo) GetSelfieFromOther(notPublishedById bson.ObjectId,
	cursor *bson.ObjectId, listLen int, page int) ([]SelfieRecord, error) {
	c := repo.DBConn.DB(DB_NAME).C(selfieC)
	findCriteria := bson.M{"author": bson.M{"$ne": notPublishedById}}
	if cursor != nil {
		findCriteria = bson.M{"$and": []interface{}{
			bson.M{"author": bson.M{"$ne": notPublishedById}},
			bson.M{"_id": bson.M{"$lt": cursor}},
		},
		}
	}
	var results []SelfieRecord
	res := c.Find(findCriteria).Sort("-created_on").Limit(listLen)
	if page > 0 {
		res = res.Skip(page * listLen)
	}
	err := res.All(&results)
	//err := c.Find(findCriteria).Sort("-created_on").Limit(listLen).All(&results)
	log.Println("selfies : ", len(results), "criteria ", findCriteria)
	if err != nil {
		log.Println("error ", err.Error())
		return nil, err

	}
	return results, nil
}

// Register a vote on a selfie
func (repo *SelfieMGORepo) Vote(user bson.ObjectId, selfie bson.ObjectId, voteValue VoteValue) error {
	var selfieModif bson.M
	if voteValue == TOP {
		selfieModif = bson.M{"$inc": bson.M{"number_of_plus": 1}}
	}
	if voteValue == OK {
		selfieModif = bson.M{"$inc": bson.M{"number_of_neutral": 1}}
	}
	if voteValue == OUT {
		selfieModif = bson.M{"$inc": bson.M{"number_of_minus": 1}}
	}
	// -> Start of transation
	voteId := bson.NewObjectId()
	ops := []txn.Op{
		{
			// Check the user exists
			C:      usersC,
			Id:     user,
			Assert: txn.DocExists,
		},
		{
			// Check the voter is not the author of the selfie
			C:      selfieC,
			Id:     selfie,
			Assert: bson.M{"author": bson.M{"$ne": user}},
		},
		{
			// Check the has not voted on the pic exists
			C:      usersC,
			Id:     user,
			Assert: bson.M{"selfie-voted": bson.M{"$nin": []bson.ObjectId{selfie}}},
			Update: bson.M{"$addToSet": bson.M{"selfie-voted": selfie}},
		},
		{
			// Increment the selfie stats
			C:      selfieC,
			Id:     selfie,
			Assert: txn.DocExists,
			Update: selfieModif,
		},
		{
			// We insert the vote in the vote history
			C:      votesC,
			Id:     voteId,
			Insert: VoteRecord{ID: voteId, Voter: user, Selfie: selfie, Value: voteValue, CreatedOn: time.Now().UTC()},
		},
	}
	id := bson.NewObjectId() // Optional
	// -> End of the transaction
	log.Printf("Running transaction on DB %s \n", DB_NAME)
	tc := repo.DBConn.DB(DB_NAME).C(transactionC)
	runner := txn.NewRunner(tc)
	err := runner.Run(ops, id, nil)
	if err != nil {
		log.Printf("Vote transaction error : %s \n", err.Error())
	}
	return err

}

// find the vote a user made on a selfie
func (repo *SelfieMGORepo) FindVoteFor(selfie bson.ObjectId, user bson.ObjectId) (VoteRecord, error) {
	var record VoteRecord
	err := repo.DBConn.DB(DB_NAME).C(votesC).Find(
		bson.M{"$and": []bson.M{
			{"voter": user},
			{"selfie": selfie},
		},
		}).One(&record)
	return record, err
}

func fromSelfie(s *data.Selfie) SelfieRecord {
	r := SelfieRecord{}
	if len(s.ID) == 0 {
		r.ID = bson.NewObjectId()
		s.ID = r.ID
	}
	r.Author = s.Author
	r.Image = s.Picture
	r.Tags = s.Tags
	r.Gender = string(s.Gender)
	r.Comment = s.Comment
	r.Keywords = make([]string, 0, len(r.Tags)*3)
	var wordCount int = 0
	for _, t := range r.Tags {
		if len(t.Category) > 0 {
			r.Keywords = append(r.Keywords, t.Category)
			wordCount++
		}
		if len(t.Brand) > 0 {
			r.Keywords = append(r.Keywords, t.Brand)
			wordCount++
		}
		if len(t.Color) > 0 {
			r.Keywords = append(r.Keywords, t.Color)
			wordCount++
		}
	}
	r.CreatedOn = time.Now()
	s.Keywords = r.Keywords
	r.RelatedArticle = s.RelatedArticle
	return r
}

func (sr SelfieRecord) SelfieData() *data.Selfie {
	sd, _ := data.NewSelfie(sr.Author, sr.Comment, sr.Tags, sr.Image)
	sd.Keywords = make([]string, 0, len(sr.Tags)*3)
	var wordCount int = 0
	for _, t := range sr.Tags {
		if len(t.Category) > 0 {
			sd.Keywords = append(sd.Keywords, t.Category)
			wordCount++
		}
		if len(t.Brand) > 0 {
			sd.Keywords = append(sd.Keywords, t.Brand)
			wordCount++
		}
		if len(t.Color) > 0 {
			sd.Keywords = append(sd.Keywords, t.Color)
			wordCount++
		}
	}
	sd.RelatedArticle = sr.RelatedArticle
	return sd
}
