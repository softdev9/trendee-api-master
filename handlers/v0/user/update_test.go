package user

import (
	"bytes"
	"errors"
	//"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestUpdateHandler_ExtractUpdateMap(t *testing.T) {

	cases := []struct {
		Fields        string
		Values        string
		Expected      map[string]string
		ErrorExpected error
	}{
		{
			Fields: "email",
			Values: "kev.legoff@gmail.com",
			Expected: map[string]string{
				"email": "kev.legoff@gmail.com",
			},
			ErrorExpected: nil,
		},
		{
			Fields:        "",
			Values:        "",
			Expected:      map[string]string{},
			ErrorExpected: nil,
		},
		{
			Fields:        "username",
			Values:        ", test,",
			Expected:      nil,
			ErrorExpected: errors.New("test is ok"),
		},
	}
	for _, tc := range cases {
		form := url.Values{}
		form.Add("fields", tc.Fields)
		form.Add("values", tc.Values)
		req, err := http.NewRequest("PUT", "", bytes.NewBufferString(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Errorf("Not able to create the request")
		}
		result, err := ExtractUpdateMap(req)
		if tc.ErrorExpected != nil && err == nil {
			t.Errorf("An error should have been thrown but was not")
		}
		equals := reflect.DeepEqual(tc.Expected, result)
		if !equals {
			t.Errorf("the maps are different got %#v , \n expected %#v \n ", result, tc.Expected)
		}
	}
}

func TestUpdateHandler_ValidateFields(t *testing.T) {

	cases := []struct {
		Input    map[string]string
		Expected bool
	}{
		{
			Input: map[string]string{
				"email":       "kev.legoff@gmail.com",
				"username":    "kevinlegoff",
				"description": "A little description about the user",
				"birthday":    "11-08-1988",
				"country":     "France",
				"city":        "Paris",
			},
			Expected: true,
		},
		{
			Input: map[string]string{
				"mail": "kev.legoff.com",
			},
			Expected: false,
		},
		{
			Input: map[string]string{
				"email": "email_invalid",
			},
			Expected: false,
		},
	}
	for _, c := range cases {
		valid := IsValidFieldListAndValue(c.Input)
		if valid != c.Expected {
			t.Errorf("Valid shoud be %t for map %#v", c.Expected, c.Input)
		}
	}
}

func TestUpdateHandler_ValidField(t *testing.T) {
	cases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "email",
			Expected: true,
		},
		{
			Input:    "gender",
			Expected: true,
		},
		{
			Input:    "username",
			Expected: true,
		},
		{
			Input:    "description",
			Expected: true,
		},
		{
			Input:    "birthday",
			Expected: true,
		},
		{
			Input:    "city",
			Expected: true,
		},
		{
			Input:    "country",
			Expected: true,
		},
		{
			Input:    "website",
			Expected: true,
		},
		{
			Input:    "invalid",
			Expected: false,
		},
	}
	for _, c := range cases {
		fieldName := c.Input
		valid := Field(fieldName).Ok()
		if valid != c.Expected {
			t.Errorf("Ok should be %t for %s ", c.Expected, fieldName)
		}
	}
}

func TestUpdateHandler_ValidEmail(t *testing.T) {
	cases := []struct {
		Input    string
		Expected bool
	}{
		{
			"kevin.legoff@trendee.co",
			true,
		},
	}
	for _, c := range cases {
		result := ValidateEmail(c.Input)
		if result != c.Expected {
			t.Errorf("Email %s should be %t ", c.Input, c.Expected)
		}
	}
}

func TestUpdateHandler_ValidBirthday(t *testing.T) {
	cases := []struct {
		Input    string
		Expected bool
	}{
		{
			"12-12-1988",
			true,
		},
		{
			"32-32-1988",
			false,
		},
	}
	for _, c := range cases {
		result := ValidateBirthday(c.Input)
		if result != c.Expected {
			t.Errorf("Birthday %s should be %t ", c.Input, c.Expected)
		}
	}
}

func TestUpdateHandler_ValidValue(t *testing.T) {
	cases := []struct {
		InputType  string
		InputValue string
		Expected   bool
	}{
		{
			"email",
			"kevin.legoff@trendee.co",
			true,
		},
		{
			"email",
			"kever",
			false,
		},
		{
			"username",
			"",
			false,
		},
		{
			"gender",
			"fdfsfd",
			false,
		},
		{
			"gender",
			"male",
			false,
		},
		{
			"gender",
			"female",
			false,
		},
		{
			"gender",
			"woman",
			true,
		},
		{
			"gender",
			"men",
			true,
		},
		{
			"birthday",
			"12-12-1988",
			true,
		},
		{
			"birthday",
			"32-32-1988",
			false,
		},
	}
	for _, c := range cases {
		result := ValidateValue(c.InputType, c.InputValue)
		if result != c.Expected {
			t.Errorf("%s : %s should be %t ", c.InputValue, c.InputType, c.Expected)
		}
	}
}
