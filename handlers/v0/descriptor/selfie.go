package descriptor

import (
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	"time"
)

type SelfieDescriptor struct {
	Id               string                `json:"id"`
	Images           map[string]string     `json:"image"`
	MetaImg          data.MetaImage        `json:"meta_img"`
	Author           SelfieAuthoDescriptor `json:"author"`
	Tags             []data.Tag            `json:"tags"`
	NumberOfComments int                   `json:"number_of_comments"`
	NumberOfPlus     int                   `json:"number_of_plus"`
	NumberOfNeutral  int                   `json:"number_of_neutral"`
	NumberOfMinus    int                   `json:"number_of_minus"`
	Comment          string                `json:"comment"`
	CreatedOn        time.Time             `json:"created_on"`
	UserVote         repos.VoteValue       `json:"user_vote"`
	Similars         []ArticleDescriptor   `json:"similars"`
}

type SelfieThumbnail struct {
	Id     string            `json:"id"`
	Images map[string]string `json:"image"`
}

func NewSelfieThumbnail(selfie repos.SelfieRecord) SelfieThumbnail {
	return SelfieThumbnail{
		Id:     selfie.ID.Hex(),
		Images: selfie.Image,
	}
}

func DescribeSelfie(selfie repos.SelfieRecord, requesting *data.User, userDatasource repos.UserIdGetter, vote repos.VoteValue, similars []ArticleDescriptor) SelfieDescriptor {
	c := make(chan data.MetaImage)
	go data.MetaForImage(selfie.Image, c)
	s := SelfieDescriptor{
		Id:               selfie.ID.Hex(),
		Images:           selfie.Image,
		Author:           DescribePictureAuthor(requesting, selfie.Author, userDatasource),
		NumberOfComments: selfie.NumberOfComments,
		Tags:             selfie.Tags,
		Comment:          selfie.Comment,
		CreatedOn:        selfie.CreatedOn,
		UserVote:         vote,
		NumberOfPlus:     selfie.NumberOfPlus,
		NumberOfMinus:    selfie.NumberOfMinus,
		NumberOfNeutral:  selfie.NumberOfNeutral,
		Similars:         similars,
	}

	meta := <-c
	s.MetaImg = meta
	if selfie.CreatedOn.Year() == 1 {
		s.CreatedOn = time.Now()
	}

	return s
}
