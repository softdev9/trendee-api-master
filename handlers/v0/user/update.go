package user

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/gateways"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type UpdateResp struct {
	v0.Response
	User *descriptor.UserLight
}

type IncomingUpdate struct {
	Fields string `schema:"fields"`
	Values string `schema:"values"`
}

type Field string

func IsValidFieldListAndValue(updateMap map[string]string) bool {
	for key, value := range updateMap {
		log.Printf("Key %s Value : %s ", key, value)
		f := Field(key)
		if !f.Ok() {
			return false
		}
		if !ValidateValue(key, value) {
			return false
		}
	}
	return true
}

func ValidateEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

func ValidateBirthday(birthday string) bool {
	const shortForm = "02-01-2006"
	date, err := time.Parse(shortForm, birthday)
	if err != nil {
		log.Printf("Error %s ", err.Error())
		return false
	}
	log.Printf("Date parsed %s ", date)
	return true
}

func ValidateValue(key string, value string) bool {
	if key == "email" {
		if !ValidateEmail(value) {
			log.Printf("[DEBUG] email error for email %s", value)
			return false
		}
	}
	if key == "username" {
		if len(value) == 0 {
			log.Printf("[DEBUG] username error for username %s", value)
			return false
		}
	}
	if key == "gender" {
		if value != data.GENDER_MAN && value != data.GENDER_WOMAN {
			log.Printf("[DEBUG] gender error for %s", value)
			return false
		}
	}
	if key == "birthday" {
		if !ValidateBirthday(value) {
			log.Printf("[DEBUG] birthday error for %s", value)
			return false
		}
	}
	return true
}

func (f Field) Ok() bool {
	s := string(f)
	switch s {
	case "email",
		"birthday",
		"description",
		"city",
		"website",
		"country",
		"gender",
		"username":
		return true
	}
	return false
}

func UpdateProfile(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[DEBUG] - UPDATE PROFILE ")
	u := context.Get(req, "user").(*data.User)
	resp := UpdateResp{}
	updateMap, err := ExtractUpdateMap(req)
	log.Printf("[DEBUG] - UPDATE PROFILE %s ", updateMap)
	if err != nil {
		resp.Error = handlers.ErrorDescriptor{
			Error:     "Wrong update infos sent by client : " + err.Error(),
			ErrorCode: 102,
		}
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	userRepo := context.Get(req, repos.UserR).(repos.InfoUpdater)
	u, err = userRepo.UpdateUser(u, updateMap)
	if err != nil {
		resp.Error = handlers.ErrorDescriptor{
			Error:     "Server is unable to update user infos " + err.Error(),
			ErrorCode: 103,
		}
		handlers.Respond(rw, req, http.StatusInternalServerError, &resp)
		return
	}
	// The user has changed is email
	if email, ok := updateMap["email"]; ok {
		if len(email) > 0 {
			mailSender := context.Get(req, "mailSender").(gateways.TrendeeMailSender)
			_, err := mailSender.EmailForEmailChange(email, u.EmailChange.Random)
			if err != nil {
				log.Printf("[EMAIL ERROR] %s, ", err.Error)
			}
		}
	}
	resp.User = descriptor.NewUserLight(u, false, true)
	handlers.Respond(rw, req, http.StatusOK, &resp)
}

func ExtractUpdateMap(req *http.Request) (map[string]string, error) {
	var updateReq IncomingUpdate

	if err := handlers.DecodeForm(req, &updateReq); err != nil {
		return nil, errors.New("not  a valid api request")
	}
	log.Printf("[DEBUG] - UPDATE REQ %s ", updateReq)
	fields := strings.Split(updateReq.Fields, ",")
	values := strings.Split(updateReq.Values, ",")

	if len(fields) != len(values) {
		return nil, errors.New("Malformed requsest")
	}
	updateMap := make(map[string]string, len(fields))
	for i, key := range fields {
		if len(key) > 0 {
			updateMap[key] = values[i]
		}
	}
	if !(IsValidFieldListAndValue(updateMap)) {
		return nil, errors.New("Malformed request")
	}
	return updateMap, nil
}
