package repos

import (
	"errors"
	"log"
	"time"

	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/pborman/uuid"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/txn"
	"github.com/softdev9/trendee-api-master/data"
)

type MongoUserRepo struct {
	DBConn *mgo.Session
}

type UserRepo interface {
	SaveUser(art *data.User) error
	GetUserById(id bson.ObjectId) (*data.User, error)
	EstablishFollow(from bson.ObjectId, to bson.ObjectId) (*data.User, *data.User, error)
	UpdateUser(user *data.User, updateInfo map[string]string) (*data.User, error)
}

type UserSaver interface {
	SaveUser(art *data.User) error
}

type InfoUpdater interface {
	UpdateUser(user *data.User, updateInfo map[string]string) (*data.User, error)
}

type UserUpdater interface {
	UpdateUser(user *data.User, updateInfo map[string]string) (*data.User, error)
}

type UserIdGetter interface {
	GetUserById(id bson.ObjectId) (*data.User, error)
}

type UserUpdaterFor interface {
	UpdateFor(selector string, u *data.User) error
}

type ProfilePictureUpdater interface {
	UpdateProfilePicture(id bson.ObjectId, imgs map[string]string) error
}

type BasicInfosUpdater interface {
	UpdateBasicInfo(id bson.ObjectId, gender string, username string,
		birthdate time.Time, country string, city string) (*data.User, error)
}

type FBUserGetter interface {
	FindFBUser(string) *data.User
}

type EmailGetter interface {
	GetUserByEmail(email string) (*data.User, error)
}

type EmailUserVerifier interface {
	FindByEmail(string) bool
}

type CredentialChecker interface {
	FindByEmail(email string) bool
	CheckEmailPassword(string, string) (*data.User, *data.AuthToken, error)
}

type FollowEstablisher interface {
	EstablishFollow(from bson.ObjectId, to bson.ObjectId) (*data.User, *data.User, error)
}

type ListFinder interface {
	FindFromList(in []bson.ObjectId) ([]data.User, error)
}

func (repo *MongoUserRepo) FindFromList(in []bson.ObjectId) ([]data.User, error) {
	var users []data.User
	findCriteria := bson.M{"_id": bson.M{"$in": in}}
	if err := repo.DBConn.DB(DB_NAME).C(usersC).Find(findCriteria).All(&users); err != nil {
		log.Println(" [ERROR] ArticleRepoMGO_147:= ", err.Error())
		return nil, err
	}
	return users, nil
}

func (repo *MongoUserRepo) EstablishFollow(from bson.ObjectId, to bson.ObjectId) (*data.User, *data.User, error) {
	if from == to {
		return nil, nil, errors.New("Trying to follow your self")
	}
	ops := []txn.Op{
		{
			// Check the user exists
			C:      usersC,
			Id:     from,
			Assert: txn.DocExists,
		},
		{
			// Check the voter is not the author of the selfie
			C:      usersC,
			Id:     to,
			Assert: txn.DocExists,
		},
		{
			// Check the has not voted on the pic exists
			C:      usersC,
			Id:     from,
			Assert: bson.M{"following": bson.M{"$nin": []bson.ObjectId{to}}},
			Update: bson.M{"$addToSet": bson.M{"following": to}},
		},
		{
			// Check the has not voted on the pic exists
			C:      usersC,
			Id:     to,
			Assert: bson.M{"followed-by": bson.M{"$nin": []bson.ObjectId{from}}},
			Update: bson.M{"$addToSet": bson.M{"followed-by": from}},
		},
	}
	// Id for the transaction
	id := bson.NewObjectId()
	// -> End of the transaction
	log.Printf("Running Follow transaction on DB %s \n", DB_NAME)
	tc := repo.DBConn.DB(DB_NAME).C(transactionC)
	runner := txn.NewRunner(tc)
	err := runner.Run(ops, id, nil)
	if err != nil {
		log.Printf("Follow transaction error : %s \n", err.Error())
	}
	f, _ := repo.GetUserById(from)
	t, _ := repo.GetUserById(to)
	log.Printf("Ending Follow transaction on DB %s \n", DB_NAME)
	return f, t, err
}

func NewMongoUserRepo(db *mgo.Session) *MongoUserRepo {
	return &MongoUserRepo{DBConn: db}
}

func (repo *MongoUserRepo) CheckEmailPassword(email string, password string) (*data.User, *data.AuthToken, error) {
	c := repo.DBConn.DB(DB_NAME).C(usersC)
	var u *data.User
	if err := c.Find(bson.M{"email": email}).One(&u); err != nil {
		return nil, nil, errors.New("User Not Found")
	}
	if u != nil {
		// We check the password
		hashedPassword := data.HashPassword(password, u.Password.Salt)
		if hashedPassword != u.Password.HashPassword {
			return nil, nil, errors.New("Password Invalid")
		} else {
			token := data.NewAuthToken()
			inserted := TokenAssociated{
				UserId:    u.ID,
				Token:     token.Token,
				ExpiresOn: token.ExpiresOn,
			}
			_, err := repo.DBConn.DB(DB_NAME).C(tokensC).Upsert(bson.M{"id_user": u.ID}, inserted)
			return u, &token, err
		}
	}
	return nil, nil, errors.New("User Not Found")
}

func (repo *MongoUserRepo) FindByEmail(email string) bool {
	c := repo.DBConn.DB(DB_NAME).C(usersC)
	//Log.Info("Looking for user with email ", email)
	if count, err := c.Find(bson.M{"email": email}).Count(); err == nil {
		found := count > 0
		return found
	} else {
		return false
	}
}

func (repo *MongoUserRepo) GetUserByEmail(email string) (*data.User, error) {
	c := repo.DBConn.DB(DB_NAME).C(usersC)
	var u *data.User
	err := c.Find(bson.M{"email": email}).One(&u)
	return u, err
}

func (repo *MongoUserRepo) GetUserById(id bson.ObjectId) (*data.User, error) {
	c := repo.DBConn.DB(DB_NAME).C(usersC)
	u := &data.User{}
	err := c.FindId(id).One(&u)
	return u, err
}

func (repo *MongoUserRepo) FindFBUser(fbId string) *data.User {
	c := repo.DBConn.DB(DB_NAME).C(usersC)
	var u *data.User
	c.Find(bson.M{"idfb": fbId}).One(&u)
	return u
}

func (repo *MongoUserRepo) UpdateProfilePicture(id bson.ObjectId, imgs map[string]string) error {
	c := repo.DBConn.DB(DB_NAME).C(usersC)
	return c.UpdateId(id, bson.M{"$set": bson.M{"profilepicture": imgs}})
}

func (repo *MongoUserRepo) Update(u *data.User) error {
	if err := repo.DBConn.DB(DB_NAME).C(usersC).UpdateId(u.ID, u); err != nil {
		return err
	}
	return nil
}

func (repo *MongoUserRepo) UpdateBasicInfo(id bson.ObjectId, gender string, username string,
	birthday time.Time, country string, city string) (*data.User, error) {
	change := bson.M{
		"$set": bson.M{
			"username": username,
			"gender":   gender,
			"birthday": birthday,
			"country":  country,
			"city":     city,
		},
	}
	if err := repo.DBConn.DB(DB_NAME).C(usersC).UpdateId(id, change); err != nil {
		return nil, err
	}
	user := &data.User{}
	err := repo.DBConn.DB(DB_NAME).C(usersC).FindId(id).One(user)
	return user, err
}

func (repo *MongoUserRepo) SaveUser(u *data.User) error {
	u.ID = bson.NewObjectId()
	if err := repo.DBConn.DB(DB_NAME).C(usersC).Insert(u); err != nil {
		return err
	}
	return nil
}

func (repo *MongoUserRepo) GetUserSuggestions(id bson.ObjectId, cursorId string) ([]data.User, error) {
	u, err := repo.GetUserById(id)
	if err != nil {
		return nil, err
	}
	followedId := append(u.Following, id)
	c := repo.DBConn.DB(DB_NAME).C(usersC)
	findCriteria := bson.M{
		"$and": []interface{}{
			bson.M{"_id": bson.M{"$nin": followedId}},
			bson.M{"suggestion-index": bson.M{"$gt": 5}},
		},
	}

	if len(cursorId) > 0 {
		cursor, err := BuildIdFromString(cursorId)
		if err != nil {
			return nil, err
		}
		findCriteria = bson.M{"$and": []interface{}{
			bson.M{"_id": bson.M{"$nin": followedId}},
			bson.M{"_id": bson.M{"$lt": cursor}},
			bson.M{"suggestion-index": bson.M{"$gt": 5}},
		},
		}
	}
	var results []data.User
	err = c.Find(findCriteria).Limit(10).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

type ArticleMarker interface {
	MarkArticle(user *data.User, articleId bson.ObjectId, value string) (ArticleLikedChanged, error)
}

// Add the article id to the list of valid article
// At this point we are sure that the article is valid
// Return the value of the liked change
func (repo *MongoUserRepo) MarkArticle(user *data.User, articleId bson.ObjectId, value string) (ArticleLikedChanged, error) {
	// Prepare the return
	liked := ArticleLikedChanged{ArticleId: articleId}
	// Check if the  article already exists in the liked list
	inLikedArticle := checkInList(user.ArticleLiked, articleId)
	// Check if the article is in the dislike list
	inDisikedArticle := checkInList(user.ArticleDisliked, articleId)
	var update bson.M
	if value == "like" {
		if inLikedArticle {
			return liked, errors.New("The article has already been liked")
		}
		liked.IncNumLike = 1
		if inDisikedArticle {
			liked.IncNumDislike = -1
		}
		update = bson.M{
			"$addToSet": bson.M{"liked-article": articleId},
			"$pull":     bson.M{"disliked-article": articleId},
			"$inc":      bson.M{"number-of-article-liked": 1, "suggestion-index": 1},
		}
	}
	if value == "unlike" {
		if !inLikedArticle {
			return liked, errors.New("The article has not been liked")
		}
		liked.IncNumLike = -1
		update = bson.M{
			"$pull": bson.M{"liked-article": articleId},
			"$inc":  bson.M{"number-of-article-liked": -1, "suggestion-index": -1},
		}
	}
	if value == "dislike" {
		if inDisikedArticle {
			return liked, errors.New("The article has already been disliked")
		}
		liked.IncNumDislike = 1
		if inLikedArticle {
			liked.IncNumLike = -1
		}
		update = bson.M{
			"$pull":     bson.M{"liked-article": articleId},
			"$addToSet": bson.M{"disliked-article": articleId},
		}
	}
	if value == "undislike" {
		if !inDisikedArticle {
			return liked, errors.New("The article has not been liked")
		}
		liked.IncNumDislike = -1
		update = bson.M{
			"$pull": bson.M{"disliked-article": articleId},
		}
	}
	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	_, err := repo.DBConn.DB(DB_NAME).C(usersC).Find(bson.M{"_id": user.ID}).Apply(change, user)
	return liked, err
}

func checkInList(arts []bson.ObjectId, articleId bson.ObjectId) bool {
	for _, idArt := range arts {
		if articleId == idArt {
			return true
		}
	}
	return false
}
func (repo *MongoUserRepo) UpdateUser(u *data.User, updateInfo map[string]string) (*data.User, error) {
	description := u.Description
	if desc, ok := updateInfo["description"]; ok {
		description = desc
	}
	username := u.Username
	if un, ok := updateInfo["username"]; ok {
		username = un
	}
	birthday := u.Birthday
	if bd, ok := updateInfo["birthday"]; ok {
		const shortForm = "02-01-2006"
		birthday, _ = time.Parse(shortForm, bd)
	}
	website := u.Website
	if wb, ok := updateInfo["website"]; ok {
		website = wb
	}
	gender := u.Gender
	if g, ok := updateInfo["gender"]; ok {
		gender = g
	}
	city := u.City
	if c, ok := updateInfo["city"]; ok {
		city = c
	}
	country := u.Country
	if l, ok := updateInfo["country"]; ok {
		country = l
	}

	var mailUpdate data.EmailChange
	if email, ok := updateInfo["email"]; ok {
		mailUpdate.Email = email
		mailUpdate.Random = uuid.New()
	}

	change := bson.M{
		"$set": bson.M{
			"username":     username,
			"gender":       gender,
			"website":      website,
			"birthday":     birthday,
			"description":  description,
			"city":         city,
			"country":      country,
			"email_change": mailUpdate,
		},
	}

	if err := repo.DBConn.DB(DB_NAME).C(usersC).UpdateId(u.ID, change); err != nil {
		return nil, err
	}
	user := &data.User{}
	err := repo.DBConn.DB(DB_NAME).C(usersC).FindId(u.ID).One(user)
	return user, err
}
