package repos

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"time"
)

type TokenAssociated struct {
	UserId    bson.ObjectId `bson:"id_user"`
	Token     string        `bson:"token"`
	ExpiresOn time.Time     `bson:"expires_on"`
}

type TokenRepoMGO struct {
	DBConn *mgo.Session
}

type TokenRepo interface {
	SaveAuthToken(data.AuthToken, data.User) error
	GetAuthToken(u data.User) (error, data.AuthToken)
	GetUserWithToken(token string) (bson.ObjectId, error)
}

func NewTokenRepo(db *mgo.Session) TokenRepo {
	return &TokenRepoMGO{DBConn: db}
}

func (repo TokenRepoMGO) getTokenC() *mgo.Collection {
	c := repo.DBConn.DB(DB_NAME).C(tokensC)
	return c
}

func (repo TokenRepoMGO) SaveAuthToken(token data.AuthToken, u data.User) error {
	inserted := TokenAssociated{
		UserId:    u.ID,
		Token:     token.Token,
		ExpiresOn: token.ExpiresOn,
	}
	_, err := repo.DBConn.DB(DB_NAME).C(tokensC).Upsert(bson.M{"id_user": u.ID}, inserted)
	return err
}

func (repo *TokenRepoMGO) GetAuthToken(u data.User) (error, data.AuthToken) {
	c := repo.DBConn.DB(DB_NAME).C(tokensC)
	var t TokenAssociated
	err := c.Find(bson.M{"id_user": u.ID}).One(&t)
	if err != nil {
		return err, data.AuthToken{}
	}
	return nil, data.AuthToken{Token: t.Token, ExpiresOn: t.ExpiresOn}
}

func (repo *TokenRepoMGO) GetUserWithToken(token string) (bson.ObjectId, error) {
	var t TokenAssociated
	err := repo.getTokenC().Find(bson.M{"token": token}).One(&t)
	if err != nil {
		return "", err
	}
	return t.UserId, nil
}
