package data

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const invalidTags = "invalid_tags"     // Error message for invalid tags

type Tag struct {
	Gender      Gender   `json:"gender"`
	SubCategory string   `json:"sub_category"`
	Category    string   `json:"category"` // Id of the category
	Brand       string   `json:"brand"`    // Id of the Brand
	Color       string   `json:"color"`    // Id of the color
	Position    Location `json:"position"` // Position of the tag

}

func (t Tag) Ok() bool {
	if !t.Gender.Ok() {
		return false
	}
	if len(t.Category) == 0 {
		return false
	}
	if !t.Position.Ok() {
		return false
	}
	return true
}

func (g Gender) Ok() bool {
	if g != "woman" && g != "men" {
		log.Printf("Invalid gender")
		return false
	}
	return true
}

type Tags []Tag

func (tags Tags) Ok() bool {
	if tags == nil {
		log.Println("[WARN] Nil tag received")
		return false
	}
	for _, t := range tags {
		if !t.Ok() {
			log.Printf("[WARN] Invalid tags %s \n", t)
			return false
		}
	}
	return true
}

type Percent float32

func (p Percent) Ok() bool {
	if p < 0.0 {
		return false
	}
	if p > 1.0 {
		return false
	}
	return true
}

// Store the position of a tag on the picture
// (0,0) = top left corner
// (100, 100) = bottom right corner
type Location struct {
	PosX Percent `json:"posX" bson:"posX"`
	PosY Percent `json:"posY" bson:"posY"`
}

// Check the posX and posY are valid
func (l Location) Ok() bool {
	if !l.PosX.Ok() {
		return false
	}
	if !l.PosY.Ok() {
		return false
	}
	return true
}

// Store data about a selfie
type Selfie struct {
	ID               bson.ObjectId
	Author           bson.ObjectId
	Tags             []Tag
	Picture          map[string]string
	MetaImage        MetaImage
	Keywords         []string
	Gender           Gender
	Comment          string
	NumberOfComments int
	NumberOfPlus     int
	NumberOfNeutral  int
	NumberOfMinus    int
	CreatedOn        time.Time
	RelatedArticle   []bson.ObjectId
}

func (s Selfie) Attributes() []string {
	return s.Keywords
}

func GenderForTags(tags []Tag) Gender {
	maleCount := 0
	femaleCount := 0
	for _, g := range tags {
		if g.Gender == "men" {
			maleCount++
		} else {
			femaleCount++
		}
	}
	if maleCount > femaleCount {
		return Gender("men")
	}
	return Gender("woman")
}

func NewSelfie(u bson.ObjectId, comment string, tags []Tag, img map[string]string) (*Selfie, error) {
	//g := Gender(gender)
	gender := GenderForTags(tags)
	t := Tags(tags)
	log.Printf("[DEBUG] Tags %s", t)
	if !t.Ok() {
		return nil, errors.New("invalid_tags")
	}
	return &Selfie{
		Author:  u,
		Tags:    tags,
		Picture: img,
		Gender:  gender,
		Comment: comment,
	}, nil
}
