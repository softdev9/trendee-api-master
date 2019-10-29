package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

const GENDER_MAN = "men"

const GENDER_WOMAN = "woman"

var ValidUserType = [3]string{"designer", "blogger", "fashionista"}

type User struct {
	ID                   bson.ObjectId     `json:"_id"  bson:"_id"`
	Username             string            `json:"username"`
	Email                string            `json:"email"`
	Description          string            `json:"description"`
	Gender               string            `json:"gender"`
	Birthday             time.Time         `json:"birthday"`
	IdFb                 string            `json:"idFacebook"`
	Password             Password          `json:"hashed-password"`
	ForgotPassword       ForgotPassword    `json:"forgot-password"`
	City                 string            `json:"city"`
	Country              string            `json:"country"`
	ProfilePicture       map[string]string `json:"profile-picture"`
	UserType             string            `json:"usertype"`
	ArticleLiked         []bson.ObjectId   `json:"liked-article" bson:"liked-article"`
	ArticleDisliked      []bson.ObjectId   `json:"disliked-article" bson:"disliked-article"`
	NumberOfSelfiePosted int               `json:"selfie-posted" bson:"selfie-posted"`
	Following            []bson.ObjectId   `bson:"following"`
	FollowedBy           []bson.ObjectId   `bson:"followed-by"`
	SelfieVoted          []bson.ObjectId   `bson:"selfie-voted"`
	NumberOfArticleLiked int               `bson:"number-of-article-liked"`
	SuggestionIndex      int               `bson:"suggestion-index"`
	Website              string            `bson:"website"`
	EmailChange          EmailChange       `bson:"email_change"`
}

func (u *User) Shop() string {
	if strings.ToLower(u.Country) == "france" {
		return StoreFR
	}
	return StoreCOM
}

func (u User) IsFollowing(toId bson.ObjectId) bool {
	following := false
	for _, uId := range u.Following {
		if uId == toId {
			following = true
		}
	}
	return following
}

func (u User) IsFollowed(by bson.ObjectId) bool {
	followed := false
	for _, uId := range u.FollowedBy && u.Following{
		if uId == by {
			followed = true
		}
	}
	return followed
}

type Password struct {
	HashPassword string
	Salt         string
}

func (u *User) DefaultProfilePic() {
	if u.ProfilePicture == nil {
		u.ProfilePicture = make(map[string]string)
		log.Println("Proflie picture map initialized")
	}
	u.ProfilePicture["xlarge"] = "https://s3-us-west-2.amazonaws.com/trendee-profiles/default_user_1024.jpg"
	u.ProfilePicture["large"] = "https://s3-us-west-2.amazonaws.com/trendee-profiles/default_user_512.jpg"
	u.ProfilePicture["medium"] = "https://s3-us-west-2.amazonaws.com/trendee-profiles/default_user_256.jpg"
	u.ProfilePicture["small"] = "https://s3-us-west-2.amazonaws.com/trendee-profiles/default_user_128.jpg"
	u.ProfilePicture["xsmall"] = "https://s3-us-west-2.amazonaws.com/trendee-profiles/default_user_128.jpg"

}

type EmailChange struct {
	Email  string `bson:"email"`
	Random string `bson:"random"`
}

type ForgotPassword struct {
	RandomString string    `bson:"random"`
	GeneratedOn  time.Time `bson:"generated_on"`
}

type FBUser struct {
	IdFB      string `json:"id"`
	FristName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
	Email     string `json:"email"`
}

type Comparable interface {
	Compare(interface{}) bool
}

func HashPassword(password string, salt string) string {
	hashed := sha256.Sum256([]byte(salt + password))
	return hex.EncodeToString(hashed[:])
}

func NewUser(email string, password string, salt string) User {
	u := &User{
		Email:  email,
		Gender: "",
		Password: Password{
			Salt:         salt,
			HashPassword: HashPassword(password, salt),
		},
	}
	u.DefaultProfilePic()
	return *u
}

func newForgotPassword() ForgotPassword {
	return ForgotPassword{
		RandomString: utils.RandSeq(32),
		GeneratedOn:  time.Now().UTC(),
	}
}

func (u *User) GeneratePasswordForogoten() {
	u.ForgotPassword = newForgotPassword()
}

func (u User) HasLiked(art bson.ObjectId) bool {
	inList := false
	for _, a := range u.ArticleLiked {
		if a == art {
			inList = true
		}
	}
	return inList
}

func (u User) HasDisliked(art bson.ObjectId) bool {
	inList := false
	for _, a := range u.ArticleDisliked {
		if a == art {
			inList = true
		}
	}
	return inList
}

func (u User) Compare(toCompare User) bool {
	if u.ID.String() != toCompare.ID.String() {
		return false
	}
	if u.Username != toCompare.Username {
		return false
	}
	if u.Email != toCompare.Email {
		return false
	}
	if !u.Birthday.Equal(toCompare.Birthday) {
		return false
	}
	if u.Gender != toCompare.Gender {
		return false
	}
	if u.IdFb != toCompare.IdFb {
		return false
	}
	return true
}

func IsValidUserType(userType string) bool {
	for _, uType := range ValidUserType {
		if userType == uType {
			return true
		}
	}
	return false
}

func ParseFBDate(date string) time.Time {
	components := strings.Split(date, "/")
	if len(components) == 3 {
		day, _ := strconv.Atoi(components[1])
		month, _ := strconv.Atoi(components[0])
		year, _ := strconv.Atoi(components[2])
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	if len(components) == 1 && len(components[0]) == 4 {
		year, _ := strconv.Atoi(components[0])
		return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	if len(components) == 2 {
		day, _ := strconv.Atoi(components[1])
		month, _ := strconv.Atoi(components[0])
		return time.Date(1, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func Convert(fb FBUser) *User {
	convertedBirthday := ParseFBDate(fb.Birthday)
	return &User{
		Username: fb.FristName + " " + fb.LastName,
		Gender:   fb.Gender,
		Email:    fb.Email,
		Birthday: convertedBirthday,
		IdFb:     fb.IdFB,
		ProfilePicture: map[string]string{
			"xlarge": getFBProfilePicture(fb.IdFB, 1024),
			"large":  getFBProfilePicture(fb.IdFB, 512),
			"medium": getFBProfilePicture(fb.IdFB, 256),
			"small":  getFBProfilePicture(fb.IdFB, 128),
			"xsmall": getFBProfilePicture(fb.IdFB, 92),
		},
	}
}

func getFBProfilePicture(fbId string, size int) string {
	return fmt.Sprintf("https://graph.facebook.com/v2.6/%s/picture?height=%d&width=%d", fbId, size, size)
}
