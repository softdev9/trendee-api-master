package v0

import (
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	//"github.com/softdev9/trendee-api-master/utils"
)

type UserRegAck struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type Response struct {
	Error handlers.ErrorDescriptor `json:"error"`
}

func NewRegistrationResponse(u data.User, token data.AuthToken) *UserRegistrationResponse {
	return &UserRegistrationResponse{
		UserInfos: UserRegAck{
			Id:    u.ID.Hex(),
			Email: u.Email,
		},
		AuthToken: token,
	}
}

type UserRegistrationResponse struct {
	Error     handlers.ErrorDescriptor `json:"error"`
	UserInfos UserRegAck               `json:"user"`
	AuthToken data.AuthToken           `json:"auth_token"`
}

/*
func NewUserUpdateResponse(u *data.User) *UserUpateResponse {
	userResponse := UserResponseLight{
		Id:        u.ID.Hex(),
		FirstName: u.Firstname,
		LastName:  u.Lastname,
		Gender:    u.Gender,
		Usertype:  u.UserType,
		City:      u.City,
		Country:   u.Country,
		Birthday:  utils.DateToString(u.Birthday),
	}
	return &UserUpateResponse{
		UserResponseLight: userResponse,
	}
}

type Response struct {
	Error handlers.ErrorDescriptor `json:"error"`
}

/// User folllowed response
type UserFollowedResponse struct {
	Response
	User UserResponseLight `json:"followed"`
}

type UserSuggestionResponse struct {
	Response
	Users []UserResponseLight `json:"users"`
}

func NewUserResponseLight(u *data.User) UserResponseLight {
	userResponse := UserResponseLight{
		Id:        u.ID.Hex(),
		FirstName: u.Firstname,
		LastName:  u.Lastname,
		Gender:    u.Gender,
		Usertype:  u.UserType,
		City:      u.City,
		Country:   u.Country,
		Birthday:  utils.DateToString(u.Birthday),
	}
	return userResponse
}
*/
