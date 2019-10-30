package descriptor

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	"time"
)

type UserSuggestionDescriptor struct {
	Id          string              `json:"id"`
	Username    string              `json:"username"`
	Usertype    string              `json:"usertype"`
	Gender      string              `json:"gender"`
	Followed    bool                `json:"followed"`
	Avatar      map[string]string   `json:"avatar"`
	Thumnbnails []map[string]string `json:"thumbnails"`
}

type UserLight struct {
	Id                string            `json:"id"`
	Username          string            `json:"username"`
	Email             string            `json:"email"`
	Gender            string            `json:"gender"`
	City              string            `json:"city"`
	Country           string            `json:"country"`
	Birthday          time.Time         `json:"birthday"`
	Avatar            map[string]string `json:"avatar"`
	Description       string            `json:"description"`
	Website           string            `json:"website"`
	NumberOfFollowers int               `json:"number_of_followers"`
	NumberOfFollowing int               `json:"number_of_following"`
	Followed          bool              `json:"followed"`
}

type FolloweeDescriptor struct {
	Id       string `json:"id"`
	Followed bool   `json:"followed"`
}

func DescribeFollowee(u *data.User) *FolloweeDescriptor {
	return &FolloweeDescriptor{
		Id:       u.ID.Hex(),
		Followed: true,
	}
}

func DescribeUser(u data.User, thumbnails []map[string]string) UserSuggestionDescriptor {
	return UserSuggestionDescriptor{
		Id:       u.ID.Hex(),
		Username: u.Username,
		Usertype: u.UserType,
		Avatar:   u.ProfilePicture,
		Gender:   u.Gender,
		//Articles:  aritcles,
		Thumnbnails: thumbnails,
	}
}

//
type SelfieAuthoDescriptor struct {
	Id       string            `json:"id"`
	Username string            `json:"username"`
	Gender   string            `json:"gender"`
	Followed bool              `json:"followed"`
	Avatar   map[string]string `json:"avatar"`
}

func DescribePictureAuthor(user *data.User, toDescribe bson.ObjectId, datasource repos.UserIdGetter) SelfieAuthoDescriptor {
	describee, _ := datasource.GetUserById(toDescribe)
	followed := false
	if user != nil {
		followed = user.IsFollowing(describee.ID)
	}
	return SelfieAuthoDescriptor{
		Id:       describee.ID.Hex(),
		Username: describee.Username,
		Gender:   describee.Gender,
		Followed: followed,
		Avatar:   describee.ProfilePicture,
	}
}

// Return a pointer a user light
// This return a pointer because in most case where this function
// is used the UserLight might be nil
func NewUserLight(user *data.User, followed bool, showEmail bool) *UserLight {
	email := user.Email
	if user.EmailChange.Email != "" {
		email = user.EmailChange.Email
	}
	if !showEmail {
		email = ""
	}
	return &UserLight{
		Id:                user.ID.Hex(),
		Username:          user.Username,
		Email:             email,
		Gender:            user.Gender,
		City:              user.City,
		Country:           user.Country,
		Birthday:          user.Birthday,
		Avatar:            user.ProfilePicture,
		Description:       user.Description,
		Website:           user.Website,
		NumberOfFollowers: len(user.FollowedBy),
		NumberOfFollowing: len(user.Following),
		Followed:          followed,
	}
}
