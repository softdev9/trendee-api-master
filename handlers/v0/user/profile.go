package user

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
)

type ProfileResponse struct {
	v0.Response
	User          *descriptor.UserLight          `json:"user"`
	LikedArticles []descriptor.ArticleDescriptor `json:"liked_articles"`
	Selfies       []descriptor.SelfieThumbnail   `json:"selfies"`
}

func GetProfile(rw http.ResponseWriter, req *http.Request) {
	asker := context.Get(req, "user").(*data.User)
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	userRepo := context.Get(req, repos.UserR).(repos.UserRepo)
	seflieRepo := context.Get(req, repos.SelfieR).(repos.SelfieRepo)
	resp := ProfileResponse{}
	var target *data.User
	if user, _ := ParseUserId(req, userRepo); user != nil {
		target = user
	} else {
		target = asker
	}
	resp.User = descriptor.NewUserLight(
		target,
		target.IsFollowed(asker.ID),
		target.ID == asker.ID)
	resp.Selfies = buildSelfiesList(seflieRepo, target.ID)
	// Build the article liked list
	resp.LikedArticles = buildArticleLiked(articleRepo, brandRepo, target.ArticleLiked, asker)
	handlers.Respond(rw, req, http.StatusOK, &resp)
}

func ParseUserId(req *http.Request, usrSrc repos.UserRepo) (*data.User, error) {
	params := req.URL.Query()
	if len(params["userid"]) > 0 {
		idStr := params["userid"][0]
		idBson, err := repos.BuildIdFromString(idStr)
		if err != nil {
			return nil, err
		}
		user, err := usrSrc.GetUserById(idBson)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

func buildSelfiesList(selfieSrc repos.SelfieRepo, authorId bson.ObjectId) []descriptor.SelfieThumbnail {
	selfies, err := selfieSrc.GetSelfiePublishedBy(authorId)
	resp := make([]descriptor.SelfieThumbnail, len(selfies), len(selfies))
	if err != nil {
		log.Printf(" Error while getting selfies published by %s", err.Error())
		return resp
	}
	for i, selfie := range selfies {
		resp[i] = descriptor.NewSelfieThumbnail(selfie)
	}
	return resp
}

func buildArticleLiked(artSrc repos.ArticleRepo, brandRepo repos.BrandRepo, idsToLoad []bson.ObjectId, user *data.User) []descriptor.ArticleDescriptor {
	resp := make([]descriptor.ArticleDescriptor, len(idsToLoad), len(idsToLoad))
	if len(idsToLoad) == 0 {
		log.Printf("[DEBUG] No article liked by the user ")
		return resp
	}
	articles, err := artSrc.GetArticleInList(idsToLoad)
	if err != nil {
		log.Printf("[ERROR] Retriving the articles liked by the user %s ", err.Error())
		return resp
	}
	return descriptor.ListDescription(articles, brandRepo, user)
}
