package selfie

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	//"github.com/softdev9/trendee-api-master/handlers/v0/article"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
)

type SelfieDestailsResponse struct {
	v0.Response
	// Error
	// Selfie
	// Similar articles
	Selfie   descriptor.SelfieDescriptor    `json:"selfie"`
	Articles []descriptor.ArticleDescriptor `json:"articles"`
}

func GetDetails(rw http.ResponseWriter, req *http.Request) {
	var resp SelfieDestailsResponse
	selfie := req.URL.Query().Get("selfie")
	u := context.Get(req, "user").(*data.User)
	selfieRepo := context.Get(req, repos.SelfieR).(repos.SelfieRepo)
	userRepo := context.Get(req, repos.UserR).(repos.UserRepo)
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	selfieId, err := repos.BuildIdFromString(selfie)
	if err != nil {
		resp.Error = handlers.ErrorDescriptor{ErrorCode: 101, Error: " Invalid selfie id " + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, resp)
		return
	}
	sr, err := selfieRepo.GetSelfieById(selfieId)
	if err != nil {
		resp.Error = handlers.ErrorDescriptor{ErrorCode: 101, Error: " Selfie not found " + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, resp)
		return
	}

	log.Printf("Selfie %#v \n", sr)
	// LoadSimilarAndNotSeen(excluded []bson.ObjectId, gender string, art jaccard.Attributor, repo repos.ArticleRepo)
	sd, ad := LoadDataForSelfie(u, sr, selfieRepo, userRepo, articleRepo, brandRepo)
	resp.Articles = ad
	resp.Selfie = sd
	handlers.Respond(rw, req, http.StatusOK, resp)
}

func LoadDataForSelfie(user *data.User, sr repos.SelfieRecord, selfieSrc repos.SelfieRepo, userSrc repos.UserIdGetter, artSrc repos.ArticleRepo, brandSrc repos.BrandRepo) (descriptor.SelfieDescriptor, []descriptor.ArticleDescriptor) {
	// Get the vote var voteValue repos.VoteValue = -1
	voteValue := repos.VoteValue(-1)
	shop := ""
	if user != nil {
		if vr, err := selfieSrc.FindVoteFor(sr.ID, user.ID); err == nil {
			voteValue = vr.Value
		}
		shop = user.Shop()
	}
	log.Printf("shop %s ", shop)
	arts, err := artSrc.GetArticlesInStoreIds(shop, sr.RelatedArticle)
	log.Printf("[DEBUG] len of arts found for selfies  arts : %d  : related articles %d", len(arts), len(sr.RelatedArticle))
	if err != nil {
		log.Printf("[ERROR] Retreiving articles for selfie %s \n", err.Error())
	}
	artsD := []descriptor.ArticleDescriptor{}
	artsD = descriptor.ListDescription(arts, brandSrc, user)
	upperLimit := 4
	if len(arts) < 4 {
		upperLimit = len(arts)
	}
	log.Print("number of article do display %d from len arts %d", upperLimit, len(arts))
	return descriptor.DescribeSelfie(
			sr,
			user,
			userSrc,
			voteValue,
			artsD[:upperLimit],
		),
		artsD[upperLimit:]

}

func GetDetailsAnonymous(rw http.ResponseWriter, req *http.Request) {
	var resp SelfieDestailsResponse
	selfie := req.URL.Query().Get("selfie")
	// Dependencies
	selfieRepo := context.Get(req, repos.SelfieR).(repos.SelfieRepo)
	userRepo := context.Get(req, repos.UserR).(repos.UserRepo)
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	selfieId, err := repos.BuildIdFromString(selfie)
	if err != nil {
		resp.Error = handlers.ErrorDescriptor{ErrorCode: 1001, Error: " Invalid selfie id " + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, resp)
		return
	}
	selfieRecord, err := selfieRepo.GetSelfieById(selfieId)
	if err != nil {
		resp.Error = handlers.ErrorDescriptor{ErrorCode: 1002, Error: " Selfie not found " + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, resp)
		return
	}
	resp.Selfie.Id = selfie
	sd, ad := LoadDataForSelfie(nil, selfieRecord, selfieRepo, userRepo, articleRepo, brandRepo)
	resp.Articles = ad
	resp.Selfie = sd
	handlers.Respond(rw, req, http.StatusOK, &resp)
}
