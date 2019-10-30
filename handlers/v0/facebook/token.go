package facebook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
)

type IncomingToken struct {
	Token string `schema:"token"`
}

const (
	fbAppId         = "****************"
	fbAppSecret     = "********************************"
	fbAPPToken      = "****************|+++++++++++++++++++++"
	fbDebugTokenURL = "https://graph.facebook.com/debug_token"
	fbGetUserInfos  = "https://graph.facebook.com/me"
)

type UserAuthenticationResponse struct {
	v0.Response
	User  *descriptor.UserLight `json:"user"`
	Token *data.AuthToken       `json:"auth_token"`
}

type FBCheckTokenResponse struct {
	Data FBCheckTokenData `json:"data"`
}

type FBCheckTokenData struct {
	IsValid bool   `json:"is_valid"`
	AppId   string `json:"app_id"`
}

func (token *IncomingToken) Ok() error {
	if len(token.Token) == 0 {
		return errors.New("no token has been sent")
	}
	return nil
}

func PostFBToken(rw http.ResponseWriter, req *http.Request) {
	incoming := IncomingToken{}
	resp := UserAuthenticationResponse{}
	if err := handlers.DecodeForm(req, &incoming); err != nil {
		resp.Error.ErrorCode = 402
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	if err := checkFBToken(incoming); err != nil {
		resp.Error.ErrorCode = 100
		resp.Error.Error = "Invalid fb token sent"
		handlers.Respond(rw, req, http.StatusUnauthorized, &resp)
		return
	}
	fbUser, err := getFBProfile(incoming)
	if err != nil {
		resp.Error.ErrorCode = 101
		resp.Error.Error = "not able to get the user infos"
		handlers.Respond(rw, req, http.StatusUnauthorized, &resp)
		return
	}
	// User finder
	userFinder := context.Get(req, repos.UserR).(repos.FBUserGetter)
	userSaver := context.Get(req, repos.UserR).(repos.UserSaver)
	user := userFinder.FindFBUser(fbUser.IdFB)
	if user == nil {
		user = data.Convert(*fbUser)
		if err := userSaver.SaveUser(user); err != nil {
			resp.Error.ErrorCode = 103
			resp.Error.Error = "Unable to save the user"
			handlers.Respond(rw, req, http.StatusInternalServerError, &resp)
			return
		}
	}

	tokenRepo := context.Get(req, repos.TokenR).(repos.TokenRepo)
	err, info := tokenRepo.GetAuthToken(*user)
	if info.OK() {
		// We return the token and the user id
		resp.Token = &info
		log.Print("Token found")
	} else {
		log.Print("New token created")
		// We need to create or update a token for this user
		token := data.NewAuthToken()
		tokenRepo.SaveAuthToken(token, *user)
		resp.Token = &token
	}
	resp.User = descriptor.NewUserLight(user, false, true)
	handlers.Respond(rw, req, http.StatusOK, &resp)
}

func checkFBToken(token IncomingToken) (err error) {
	debugTokenPramas := "?input_token=" + token.Token + "&access_token=" + fbAPPToken
	res, err := http.Get(fbDebugTokenURL + debugTokenPramas)
	// fmt.Print(fbDebugTokenURL + debugTokenPramas + "\n \n")
	var facebookTokenCheck FBCheckTokenResponse
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&facebookTokenCheck)
	if facebookTokenCheck.Data.IsValid != true {
		return errors.New("Invalid token")
	}
	if facebookTokenCheck.Data.AppId != fbAppId {
		return errors.New("Invalid AppId")
	}
	return nil
}

func getFBProfile(token IncomingToken) (*data.FBUser, error) {
	appSecretProof := GenerateHMAC(token.Token, fbAppSecret)
	accessTokenParams := "?access_token=" + token.Token + "&appsecret_proof=" + appSecretProof
	res, err := http.Get(fbGetUserInfos + accessTokenParams)
	if err != nil {
		return nil, errors.New("Not able to contact facebook")
	}
	var fbUser data.FBUser
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&fbUser); err != nil {
		return nil, errors.New("Not able to parse facebook return")
	}
	return &fbUser, nil
}

/**
** Generaetes the appsecret_proof
** for more infos see https://developers.facebook.com/docs/graph-api/securing-requests#appsecret_proof
**  should return 9b78e54000bc8186e7ae4704742553aa84f9d1eedab7874122b70681276797b2
**/
func GenerateHMAC(message string, key string) string {
	keyByte := []byte(key)
	h := hmac.New(sha256.New, keyByte)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
