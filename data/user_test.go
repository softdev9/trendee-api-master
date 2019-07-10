package data

import (
	//"fmt"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	u := User{
		ID:       "test_id",
		Username: "kevin le goff",
		Email:    "kev.legoff@gmail.com",
		Birthday: time.Date(1988, time.August, 11, 0, 0, 0, 0, time.UTC),
	}
	if u.ID == "" {
		t.Error("Id should have been retreived")
	}
}

func TestIsValidUserType(t *testing.T) {
	var cases = []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "wrong input",
			Expected: false,
		},
		{
			Input:    "designer",
			Expected: true,
		},
	}
	for _, c := range cases {
		output := IsValidUserType(c.Input)
		if c.Expected != output {
			t.Error("Should have been ", c.Expected, " for ", c.Input)
		}
	}
}

func TestGetFBProfile(t *testing.T) {
	expected := "https://graph.facebook.com/v2.6/587678864/picture?height=256&width=256"
	got := getFBProfilePicture("587678864", 256)
	if expected != got {
		t.Errorf("Expected %s but got %s", expected, got)
	}
}

func TestShopFromCountry(t *testing.T) {
	u := User{
		ID:       "test_id",
		Username: "kevin le goff",
		Email:    "kev.legoff@gmail.com",
		Birthday: time.Date(1988, time.August, 11, 0, 0, 0, 0, time.UTC),
		Country:  "france",
	}
	shop := u.Shop()
	if shop != StoreFR {
		t.Errorf("The store should be %s but is %s", StoreFR, shop)
	}
	u.Country = "england"
	shop = u.Shop()
	if shop != StoreCOM {
		t.Errorf("The store should be %s but is %s", StoreCOM, shop)
	}
}
