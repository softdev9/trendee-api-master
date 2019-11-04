package handlers

import (
	"bytes"
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/schema"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/gateways"
	"github.com/softdev9/trendee-api-master/repos"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Adapter func(http.Handler) http.Handler

var api_keys = map[string]string{
	"E8g9I5Gq3beJEwFwFOJC2ER22xpy631b": "iOS",
	"tTb9jow0v4DgBJqLBmbiI5nEGVUa6AuO": "android",
	"7a4dFV0cn7zDAs6ri0Aity2z0Z3nhb62": "webtrendee",
}

type ErrorDescriptor struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
}

type OK interface {
	Ok() error
}

func Logging() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Printf("[INFO] Request received on endpoint : %s -> %s  \n", r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
			elapsed := time.Since(start)
			log.Printf("[PERF] Response produced in %s", elapsed)
		})
	}
}

func WithAuth() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := context.Get(r, "version")
			if v == 0 {
				h.ServeHTTP(w, r)
			} else {
				key := r.Header.Get("Authorization")
				if len(key) > 0 {
					client := api_keys[key]
					if len(client) > 0 {
						context.Set(r, "client", client)
						h.ServeHTTP(w, r)
					} else {
						Respond(w, r, http.StatusUnauthorized, &unauthorized_respone{Cause: "no client matching key"})
					}
				} else {
					Respond(w, r, http.StatusUnauthorized, &unauthorized_respone{Cause: "No Authorization header provided"})

				}
			}
		})
	}
}

func WithAPIVersion() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accept := r.Header.Get("Accept")
			var version = 0
			if accept != "" {
				versionString := strings.Split(accept, ".")
				if len(versionString) == 3 {
					v, err := strconv.Atoi(versionString[2])
					if err != nil {
						version = 0
						fmt.Println(err.Error())
					} else {
						fmt.Println("Version detected ", v)
						version = v
					}
				} else {
					fmt.Println("Length ", len(versionString))
				}
			}
			fmt.Println("Accpet value ", accept, "version ", version)
			context.Set(r, "version", version)
			h.ServeHTTP(w, r)
		})
	}
}

type unauthorized_respone struct {
	Cause string `json:"cause"`
}

func WithRepos(session *mgo.Session) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbCopy := session.Copy()
			//defer dbCopy.Close()
			context.Set(r, repos.TokenR, repos.NewTokenRepo(dbCopy.Copy()))
			context.Set(r, repos.SelfieR, repos.NewSelfieRepoMGO(dbCopy.Copy()))
			context.Set(r, repos.BrandR, repos.NewBrandRepo(dbCopy.Copy()))
			context.Set(r, repos.UserR, repos.NewMongoUserRepo(dbCopy.Copy()))
			context.Set(r, repos.ColorsR, repos.NewColorRepo(dbCopy.Copy()))
			context.Set(r, repos.ArticleR, repos.NewArticleRepo(dbCopy.Copy()))
			h.ServeHTTP(w, r)
			context.Clear(r)
		})
	}
}

func WithMailSender(sender gateways.TrendeeMailSender) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Set(r, "mailSender", sender)
			h.ServeHTTP(w, r)
			context.Clear(r)
		})
	}
}

// Check that the token sent together with the request is valid
func WithAuthToken(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenR := context.Get(r, repos.TokenR).(repos.TokenRepo)
		userR := context.Get(r, repos.UserR).(repos.UserIdGetter)
		var token string
		if r.Method == "POST" || r.Method == "PUT" {
			token = r.FormValue("token")
		}
		if r.Method == "GET" {
			params := r.URL.Query()
			if len(params["token"]) == 0 {
				Respond(w, r, http.StatusUnauthorized, &ErrorDescriptor{Error: "Send a token", ErrorCode: 100})
				return
			}
			token = params["token"][0]
		}
		log.Printf("[INFOS] Token from user : %v ", token)
		id, err := tokenR.GetUserWithToken(token)
		u, err := userR.GetUserById(id)
		if err != nil {
			//context.Clear(r)
			log.Print("[INFOS] [ERROR] Token submitted is invalid")
			Respond(w, r, http.StatusUnauthorized, &ErrorDescriptor{Error: "Invalid token ", ErrorCode: 100})
			return
		}
		log.Printf("[INFOS] User ID : %v ", u.ID)
		context.Set(r, "user", u)
		f(w, r)
	}
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func DecodeForm(req *http.Request, o interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	var formDecoder = schema.NewDecoder()
	formDecoder.Decode(o, req.Form)
	if obj, ok := o.(OK); ok {
		if err := obj.Ok(); err != nil {
			return err
		}
	}
	return nil
}

func Respond(rw http.ResponseWriter, req *http.Request, status int, d interface{}) {
	// Create a buffer for the json encodoer
	var buf bytes.Buffer // A Buffer needs no initialization.
	// Check if the data implement the public interface if yes then we change it
	if obj, ok := d.(data.Public); ok {
		d = obj.Public()
	}
	if err := json.NewEncoder(&buf).Encode(d); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	if _, err := io.Copy(rw, &buf); err != nil {

	}

}
