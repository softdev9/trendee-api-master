package article

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/algos/jaccard"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
	"time"
)

type InFicheArticleReq struct {
	ArticleId string
}

type OutFicheArt struct {
	v0.Response
	Art     *descriptor.ArticleDescriptor  `json:"article"`
	Similar []descriptor.ArticleDescriptor `json:"similars"`
}

func newInFicheArticleReq(id string) *InFicheArticleReq {
	return &InFicheArticleReq{ArticleId: id}
}

func (req *InFicheArticleReq) Ok() error {
	if len(req.ArticleId) == 0 {
		return errors.New("no_article_id")
	}
	return nil
}

func ArticleDetails(rw http.ResponseWriter, req *http.Request) {
	// We can get the user from the context
	// Get the user
	user := context.Get(req, "user").(*data.User)
	incArt := newInFicheArticleReq(req.URL.Query().Get("article"))
	// We check the the article id is ok
	if err := incArt.Ok(); err != nil {
		handlers.Respond(
			rw,
			req,
			http.StatusOK,
			OutFicheArt{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{
						ErrorCode: 101, Error: err.Error()},
				},
			})
		return
	}
	// Get the article repo
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	art, err := articleRepo.GetArticleWithId(incArt.ArticleId)
	if err != nil {
		handlers.Respond(
			rw,
			req,
			http.StatusOK,
			OutFicheArt{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{
						ErrorCode: 102, Error: err.Error()},
				},
			})
		return
	}
	// We have the article
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	hasLiked := false
	hasDisliked := false
	idArt, _ := repos.BuildIdFromString(incArt.ArticleId)
	for _, idLiked := range user.ArticleLiked {
		if idLiked == idArt {
			hasLiked = true
		}
	}
	for _, idDisliked := range user.ArticleDisliked {
		if idDisliked == idArt {
			hasDisliked = true
		}
	}
	likedAndDisliked := append(user.ArticleLiked, user.ArticleDisliked...)
	likedAndDisliked = append(likedAndDisliked, art.Id)
	arts := LoadSimilarAndNotSeen(
		likedAndDisliked,
		art.Gender,
		art.Category,
		art, articleRepo,
	)
	artD := descriptor.Article(art, brandRepo, hasLiked, hasDisliked)
	handlers.Respond(
		rw,
		req,
		http.StatusOK,
		OutFicheArt{
			Art:     &artD,
			Similar: descriptor.ListDescriptor(arts, brandRepo, 20),
		})
}

func LoadSimilarAndNotSeen(excluded []bson.ObjectId, gender string, category string, art jaccard.Attributor, repo repos.ArticleRepo) []*data.Article {
	start := time.Now()
	arts, _ := repo.GetArticleNotIds(excluded, gender, category)
	log.Printf("Selfies Keyword %s", art.Attributes())
	loadingTime := time.Since(start)
	log.Printf("[PERF]  Loading %v articles took %v", len(arts), loadingTime)
	log.Printf("[DEBUG] Running jaccard on %d articles for gender %s and cat %s", len(arts), gender, category)
	log.Printf("[DEBUG] We got %d to analyze", len(arts))
	res := jaccard.JaccardOrder(art, arts)
	return res
}

func GetDetailsAnonymous(rw http.ResponseWriter, req *http.Request) {
	incArt := newInFicheArticleReq(req.URL.Query().Get("article"))
	if err := incArt.Ok(); err != nil {
		handlers.Respond(
			rw,
			req,
			http.StatusOK,
			OutFicheArt{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{
						ErrorCode: 101, Error: err.Error()},
				},
			})
		return
	}
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	art, err := articleRepo.GetArticleWithId(incArt.ArticleId)
	if err != nil {
		handlers.Respond(
			rw,
			req,
			http.StatusOK,
			OutFicheArt{
				Response: v0.Response{
					Error: handlers.ErrorDescriptor{
						ErrorCode: 102, Error: err.Error()},
				},
			})
		return
	}

	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	hasLiked := false
	hasDisliked := false
	arts := LoadSimilarAndNotSeen(
		nil,
		art.Gender,
		art.Category,
		art, articleRepo,
	)
	artD := descriptor.Article(art, brandRepo, hasLiked, hasDisliked)
	handlers.Respond(
		rw,
		req,
		http.StatusOK,
		OutFicheArt{
			Art:     &artD,
			Similar: descriptor.ListDescriptor(arts, brandRepo, 20),
		},
	)
}
