package data

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"testing"
	"log"
)

var validTagList = []Tag{
	{
		Gender:      "men",
		Category:    "dresses",
		Brand:       "h&m",
		SubCategory: "",
		Color:       "green",
		Position: Location{
			PosX: 0.25,
			PosY: 0.65,
		},
	},
	{
		Gender:   "men",
		Category: "dresses",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.55,
			PosY: 0.95,
		},
	},
	{
		Gender:   "woman",
		Category: "dresses",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.33,
			PosY: 0.57,
		},
	},
}

var validTagListFemale = []Tag{
	{
		Gender:   "woman",
		Category: "dresses",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.33,
			PosY: 0.57,
		},
	},
	{
		Gender:   "male",
		Category: "dresses",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.33,
			PosY: 0.57,
		},
	},
}

var invalidTagListNoCat = []Tag{
	{
		Category: "",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.33,
			PosY: 0.57,
		},
	},
}

var invalidTagListNoGender = []Tag{
	{
		Category: "dresses",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.33,
			PosY: 0.57,
		},
	},
}

var invalidTagListPosX = []Tag{
	{
		Category: "category",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 1.01,
			PosY: 0.34,
		},
	},
}

var invalidTagListPosX2 = []Tag{
	{
		Category: "category",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: -0.01,
			PosY: 0.34,
		},
	},
}

var invalidTagListPosY = []Tag{
	{
		Category: "category",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.21,
			PosY: -0.01,
		},
	},
}

var invalidTagListPosY2 = []Tag{
	{
		Category: "category",
		Brand:    "h&m",
		Color:    "green",
		Position: Location{
			PosX: 0.21,
			PosY: 1.01,
		},
	},
}

func TestGenderValidation(t *testing.T) {
	var genderTests = []struct {
		gender   string // input
		expected bool   // expected result
	}{
		{"f", false},
		{"woman", true},
		{"men", true},
	}
	for _, test := range genderTests {
		// NewSelfie(u bson.ObjectId, gender string, comment string, tags []Tag, img map[string]string)
		ok := Gender(test.gender).Ok()
		if test.expected != ok {
			t.Errorf("Test should have returned %t for %s but has returned %t", test.expected, test.gender, test.expected)
		}
	}
}

func TestGenderDetermindation(t *testing.T) {
	expected := Gender("men")
	g := GenderForTags(validTagList)
	if g != expected {
		t.Errorf("the gender expected was %s, but we got %s for the taglist %v ", expected, g, validTagList)
	}
	g = GenderForTags(validTagListFemale)
	expected = Gender("woman")
	if g != expected {
		t.Errorf("the gender expected was %s, but we got %s for the taglist %v ", expected, g, validTagList)
	}
}

func TestTagsValidation(t *testing.T) {
	var genderTests = []struct {
		tags     []Tag // input
		expected error // expected result
	}{
		{validTagList, nil},
		{nil, errors.New(invalidTags)},
		{invalidTagListNoGender, errors.New(invalidTags)},
		{invalidTagListNoCat, errors.New(invalidTags)},
		{invalidTagListPosX, errors.New(invalidTags)},
		{invalidTagListPosX2, errors.New(invalidTags)},
		{invalidTagListPosY, errors.New(invalidTags)},
		{invalidTagListPosY2, errors.New(invalidTags)},
	}
	for _, test := range genderTests {
		_, err := NewSelfie(bson.NewObjectId(), "blue test red", test.tags, nil)
		if test.expected != nil {
			if err == nil {
				t.Errorf("No error returned for %v", test.tags)
			} else {
				if err.Error() != test.expected.Error() {
					t.Error("Invalid error returned")
					return
				}
			}
		} else {
			if err != nil {
				t.Error("Error should not be thrwon")
				return
			}
		}
	}
}
