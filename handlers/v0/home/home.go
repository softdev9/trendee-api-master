package home

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/selfie"
	//"github.com/softdev9/trendee-api-master/handlers/v0/article"
	"log"
	"net/http"
	"strconv"

	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
)

type HomeResponse struct {
	v0.Response
	Page HomePage `json:"home_page"` // One page of the home
}

type HomePage struct {
	Users  []descriptor.UserSuggestionDescriptor `json:"users_suggestion"` // The list of users
	Selfie []descriptor.SelfieDescriptor         `json:"selfies"`          // The list of selfies that the user can vote on
}

var numberOfArticlePerPage int = 20
var numberOfSuggestionPerPage int = 10

func GetHome(rw http.ResponseWriter, req *http.Request) {
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	userRepo := context.Get(req, repos.UserR).(*repos.MongoUserRepo)
	selfieRepo := context.Get(req, repos.SelfieR).(*repos.SelfieMGORepo)
	u := context.Get(req, "user").(*data.User)
	// Verify the user token
	users := getMostSimilarUser(u, userRepo, selfieRepo, articleRepo, "")

	page := 0
	pageParam := req.URL.Query().Get("page")
	if pageParam != "" {
		n, err := strconv.Atoi(pageParam)
		if err == nil && n >= 0 {
			page = n
		}
	}

	listLen := 20
	sizeParam := req.URL.Query().Get("size")
	if sizeParam != "" {
		n, err := strconv.Atoi(sizeParam)
		if err == nil && n > 0 {
			listLen = n
		}
	}

	selfies, _ := getSelfies(u, selfieRepo, userRepo, articleRepo, brandRepo, listLen, page)
	homePage := HomePage{
		Users:  users,
		Selfie: selfies,
	}
	// Feed the articles
	handlers.Respond(rw, req, http.StatusOK, &HomeResponse{Page: homePage})
}

type ProducedSelfie struct {
	Index  int
	Selfie descriptor.SelfieDescriptor
}

func getSelfies(user *data.User, selfieRepo *repos.SelfieMGORepo,
	userRepo repos.UserIdGetter, artSrc repos.ArticleRepo,
	brandSrc repos.BrandRepo, listLen int, page int) ([]descriptor.SelfieDescriptor, error) {
	selfies, err := selfieRepo.GetSelfieFromOther(user.ID, nil, listLen, page)
	if err != nil {
		return nil, err
	}
	numberOfSelfies := len(selfies)
	resp := make([]descriptor.SelfieDescriptor, numberOfSelfies, numberOfSelfies)
	// Laucnh each stuff in go routine
	produced := make(chan ProducedSelfie)
	for i, s := range selfies {
		go produceSelfieRepsonse(i, s, user, selfieRepo, userRepo, artSrc, brandSrc, produced)
	}
	numberOfSelfieDescribed := 0
	for numberOfSelfies-numberOfSelfieDescribed > 0 {
		sd := <-produced
		//log.Printf("Response from %d  : %#v ", numberOfSelfieDescribed, sd.Similars)
		resp[sd.Index] = sd.Selfie
		numberOfSelfieDescribed = numberOfSelfieDescribed + 1
	}
	log.Printf("Number of selfie done  %d", numberOfSelfieDescribed)
	return resp, nil

}

func produceSelfieRepsonse(index int, s repos.SelfieRecord, user *data.User, selfieRepo *repos.SelfieMGORepo,
	userRepo repos.UserIdGetter, artSrc repos.ArticleRepo,
	brandSrc repos.BrandRepo, produced chan ProducedSelfie) {
	selfieDesc, _ := selfie.LoadDataForSelfie(user, s, selfieRepo, userRepo, artSrc, brandSrc)
	produced <- ProducedSelfie{
		Index:  index,
		Selfie: selfieDesc,
	}
	return

}

func getMostSimilarUser(user *data.User, userRepo *repos.MongoUserRepo, selfieRepo *repos.SelfieMGORepo, artRepo repos.ArticleRepo, userCursor string) []descriptor.UserSuggestionDescriptor {
	// Get the list of user id with three seflie posted who have posted
	users, err := userRepo.GetUserSuggestions(user.ID, userCursor)
	if err != nil {
		return nil
	}
	size := numberOfSuggestionPerPage
	if len(users) == 0 {
		return nil
	}
	if len(users) < numberOfSuggestionPerPage {
		size = len(users)
	}
	resp := make([]descriptor.UserSuggestionDescriptor, size, size)
	for i, u := range users {
		// Build the thumbnail with selfies or with articles
		thumbnails := make([]map[string]string, 3, 3)
		// We want the selfie first
		selfies, _ := selfieRepo.GetSelfiePublishedBy(u.ID)
		log.Printf(" Number of selfie from user  %d", len(selfies))
		numberOfThumnailAdded := 0
		maxSelfieToAdd := 3
		maxArticleToAdd := 0
		if len(selfies) < 3 {
			maxSelfieToAdd = len(selfies)
			maxArticleToAdd = 3 - maxSelfieToAdd
		}
		for numberOfThumnailAdded < maxSelfieToAdd {
			thumbnails[numberOfThumnailAdded] = selfies[numberOfThumnailAdded].Image
			numberOfThumnailAdded++
		}
		if maxArticleToAdd > 0 {
			numberOfArticleAdded := 0
			articles, _ := artRepo.GetArticleInList(u.ArticleLiked)
			if len(articles) < maxArticleToAdd {
				break
			}
			log.Printf("Size of the article list %d ->  Number of article to add %d", len(articles), maxArticleToAdd)
			for numberOfThumnailAdded < 3 && numberOfArticleAdded < maxArticleToAdd {
				log.Printf("Number of thumbnails %d", numberOfThumnailAdded)
				log.Printf("Number of size of the article list %d", len(articles))
				img := articles[numberOfArticleAdded].Image
				thumbnails[numberOfThumnailAdded] = img
				numberOfThumnailAdded++
				numberOfArticleAdded++
			}
		}
		//log.Printf("\t inserting user %q %v \n ", u.Firstname, thumbnails)
		resp[i] = descriptor.DescribeUser(u, thumbnails)
	}
	return resp

}
