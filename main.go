package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/spiderdev86/trendee-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/spiderdev86/trendee-api/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/spiderdev86/trendee-api/gateways"
	"github.com/spiderdev86/trendee-api/handlers"
	"github.com/spiderdev86/trendee-api/handlers/v0"
	"github.com/spiderdev86/trendee-api/handlers/v0/article"
	"github.com/spiderdev86/trendee-api/handlers/v0/facebook"
	"github.com/spiderdev86/trendee-api/handlers/v0/follow"
	"github.com/spiderdev86/trendee-api/handlers/v0/home"
	"github.com/spiderdev86/trendee-api/handlers/v0/login"
	"github.com/spiderdev86/trendee-api/handlers/v0/profilepicture"
	"github.com/spiderdev86/trendee-api/handlers/v0/register"
	"github.com/spiderdev86/trendee-api/handlers/v0/selfie"
	"github.com/spiderdev86/trendee-api/handlers/v0/user"
	"github.com/spiderdev86/trendee-api/handlers/v0/vote"
	"github.com/spiderdev86/trendee-api/handlers/web/appinvite"
)

const (
	usage = "-prod to launch in production mode"
)

func main() {
	// Check if the test flag is perseent
	var dbHost string
	prodEnv := flag.Bool("prod", false, usage)
	flag.Parse()
	if *prodEnv == false {
		log.Println("[INFO] In Test env")
		dbHost = os.Getenv("MONGODB_TEST")
	} else {
		log.Println("[INFO] In Prod env")
		dbHost = os.Getenv("MONGODB_PROD")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("[ERROR] Web server herror : No Port specified")
		os.Exit(1)
	}
	if dbHost == "" {
		log.Println("[ERROR] Mongo db error : MONGODB_PROD not set OR MONGODB_TEST if you set the test flag ")
		os.Exit(1)
	}
	dbSession, err := mgo.Dial(dbHost)
	if err != nil {
		log.Println("[ERROR] not able to connect to db host ", err)
		os.Exit(1)
	}
	defer dbSession.Close()
	sender := gateways.InitMailGun()
	r := mux.NewRouter()

	r.HandleFunc("/home", handlers.WithAuthToken(home.GetHome)).Methods("GET")
	r.HandleFunc("/article/mark", handlers.WithAuthToken(article.PostLikeArticle)).Methods("POST")
	r.HandleFunc("/article", handlers.WithAuthToken(article.ArticleDetails)).Methods("GET")
	r.HandleFunc("/colors", v0.GetColorList).Methods("GET")
	r.HandleFunc("/brands", v0.GetBrandList).Methods("GET")
	r.HandleFunc("/user/follow", handlers.WithAuthToken(follow.PostFollow)).Methods("POST")
	r.HandleFunc("/user/register", register.RegisterPost).Methods("POST")
	r.HandleFunc("/user/login", login.Login).Methods("POST")
	r.HandleFunc("/user/facebook", facebook.PostFBToken).Methods("POST")
	r.HandleFunc("/user/profilepicture", handlers.WithAuthToken(profilepicture.UploadProfilePic)).Methods("POST")
	r.HandleFunc("/user", v0.GetProfile).Methods("GET")
	r.HandleFunc("/user/follower", handlers.WithAuthToken(user.GetFollowers)).Methods("GET")
	r.HandleFunc("/user/following", handlers.WithAuthToken(user.GetFollowing)).Methods("GET")
	r.HandleFunc("/user/profile", handlers.WithAuthToken(user.GetProfile)).Methods("GET")
	r.HandleFunc("/user", handlers.WithAuthToken(user.UpdateProfile)).Methods("PUT")
	r.HandleFunc("/selfie", handlers.WithAuthToken(selfie.UploadSelfie)).Methods("POST")
	r.HandleFunc("/selfie", handlers.WithAuthToken(selfie.GetDetails)).Methods("GET")
	r.HandleFunc("/selfie/vote", handlers.WithAuthToken(vote.PostVote)).Methods("POST")
	r.HandleFunc("/appinvite", appinvite.GetAppInviteHTML).Methods("GET")
	r.HandleFunc("/anonymous/selfie", selfie.GetDetailsAnonymous).Methods("GET")
	r.HandleFunc("/anonymous/article", article.GetDetailsAnonymous).Methods("GET")

	http.Handle("/", handlers.Adapt(
		r,
		handlers.WithMailSender(sender),
		handlers.WithRepos(dbSession),
		handlers.Logging(),
		handlers.WithAPIVersion(),
	))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("[FATAL] ListenAndServe: ", err)
	}
}
