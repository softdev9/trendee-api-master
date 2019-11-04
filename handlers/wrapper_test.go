package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

type TestDecoded struct {
	Name string `schema:"name"`
}

func (t TestDecoded) Ok() error {
	return errors.New("Test Error")
}

func TestDecodeForm(t *testing.T) {

	cases := []struct {
		Incoming string
		Expected error
	}{
		{
			Incoming: "kevin",
			Expected: errors.New("Test Error"),
		},
	}

	for _, test := range cases {
		form := url.Values{}
		form.Add("name", test.Incoming)
		// Build a request
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Error("Fail to create the the request")
		}
		// Create a request recorde
		result := TestDecoded{}
		err = DecodeForm(req, &result)
		// Check we have updated
		if test.Expected != nil {
			if test.Expected.Error() != err.Error() {
				t.Error("The test is not returning the good error ")
			}
		} else {
			if err != nil {
				t.Error("should not return an error ")
			} else {
				if result.Name != test.Incoming {
					t.Error("Should return the same stuff but we got ",
						result, " for expected ", test.Incoming)
				}
			}
		}

	}

}
