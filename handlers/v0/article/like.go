package article

import (
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

type IncomingLike struct {
	ArticleId string `schema:"article"`
	Operation string `schema:"operation"`
}

type LikeResponse struct {
	v0.Response
	Art *descriptor.ArticleDescriptor `json:"article"`
}

func (request *IncomingLike) Ok() error {
	if len(request.ArticleId) < 1 {
		return errors.New("No object id")
	}
	if !repos.IsValidBsonObjectId(request.ArticleId) {
		return errors.New("Invalid article object id sent")
	}
	if request.Operation != "like" && request.Operation != "unlike" && request.Operation != "dislike" && request.Operation != "undislike" {
		return errors.New("Invalid operation parameter should be like or dislike")
	}
	return nil
}

func PostLikeArticle(rw http.ResponseWriter, req *http.Request) {
	// Build the response
	var resp LikeResponse
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandNameFinder)
	userRepo := context.Get(req, repos.UserR).(repos.ArticleMarker)
	user := context.Get(req, "user").(*data.User)
	like := IncomingLike{}
	// Decode the incoming like article request
	if err := handlers.DecodeForm(req, &like); err != nil {
		resp.Error = handlers.ErrorDescriptor{ErrorCode: 101, Error: "Not able to decode the incoming request" + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, resp)
		return
	}
	// Check The article exists
	art, err := articleRepo.GetArticleWithId(like.ArticleId)
	if err != nil {
		resp.Response.Error = handlers.ErrorDescriptor{ErrorCode: 102, Error: "Article not found " + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	artChange, err := userRepo.MarkArticle(user, art.Id, like.Operation)
	artLiked, err := articleRepo.UpdateNumberOfLike(artChange)
	if err != nil {
		resp.Response.Error = handlers.ErrorDescriptor{ErrorCode: 103, Error: "Not able to like article" + err.Error()}
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	log.Println("Like operation ", like.Operation, "liked ", user.HasLiked(art.Id), "disliked", user.HasDisliked(art.Id))
	artD := descriptor.Article(*artLiked, brandRepo, user.HasLiked(art.Id), user.HasDisliked(art.Id))
	handlers.Respond(rw, req, http.StatusOK, &LikeResponse{Art: &artD})
}
