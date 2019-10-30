package descriptor

import (
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"time"
)

type BrandDescriptor struct {
	Name string            `json:"name"`
	Logo map[string]string `json:"image"`
}

type ArticleThumbnail struct {
	Id       string            `json:"id"`
	Image    map[string]string `json:"image"`
	Name     string            `json:"name"`
	Price    float32           `json:"price"`
	Currency string            `json:"currency"`
}

func NewArticlThumbnail(art data.Article) ArticleThumbnail {
	return ArticleThumbnail{
		Id:       art.Id.Hex(),
		Name:     art.Name,
		Image:    art.Image,
		Price:    art.Price,
		Currency: art.Currency,
	}
}

type ArticleDescriptor struct {
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	Brand           BrandDescriptor   `json:"brand"`
	Retailer        string            `json:"retailer"`
	Price           float32           `json:"price"`
	Currency        string            `json:"currency"`
	ClickUrl        string            `json:"clickUrl"`
	Image           map[string]string `json:"image"`
	InStock         bool              `json:"inStock"`
	RegistredOn     string            `json:"registred_on"`
	Liked           bool              `json:"liked"`
	Disliked        bool              `json:"disliked"`
	Descritpion     string            `json:"description"`
	NumberOfLike    int               `json:"numberOfLike"`
	NumberOfDislike int               `json:"numberOfDislike"`
	NumberOfShare   int               `json:"numberOfShare"`
	Category        string            `json:"category"`
}

func Article(art data.Article, brandRepo repos.BrandNameFinder, like bool, disliked bool) ArticleDescriptor {
	now := time.Now()
	b, _ := brandRepo.FindByName(art.Brand)
	artD := ArticleDescriptor{
		Id:              art.Id.Hex(),
		Name:            art.Name,
		Brand:           Brand(b),
		Retailer:        art.Retailer,
		Price:           art.Price,
		Currency:        art.Currency,
		ClickUrl:        art.ClickUrl,
		Image:           art.Image,
		InStock:         art.InStock,
		Descritpion:     art.Description,
		Liked:           like,
		Disliked:        disliked,
		NumberOfDislike: art.NumberOfDislike,
		NumberOfLike:    art.NumberOfLike,
		NumberOfShare:   art.NumberOfShare,
		RegistredOn:     art.RegistredOn.UTC().String(),
		Category:        art.Category,
	}
	elapsed := time.Since(now)
	log.Printf("Converting article took %s", elapsed)
	return artD
}

func ListDescriptor(arts []*data.Article, brandRepo repos.BrandNameFinder, maxArticles int) []ArticleDescriptor {
	limit := maxArticles
	if len(arts) < maxArticles {
		limit = len(arts)
	}
	list := make([]ArticleDescriptor, limit, limit)
	log.Printf("Converting limit  %d arts size %d ", limit, len(arts) )
	for i:=0; i < limit; i++ {
		list[i] = Article(*arts[i], brandRepo, false, false)
	}
	return list
}

func ListDescription(arts []data.Article, brandRepo repos.BrandNameFinder, user *data.User) []ArticleDescriptor {
	articles := make([]ArticleDescriptor, len(arts), len(arts))
	hasLiked := false
	hasDisliked := false
	for i, art := range arts {
		if user != nil {
			hasLiked = user.HasLiked(art.Id)
			hasDisliked = user.HasDisliked(art.Id)
		}
		articles[i] = Article(
			art,
			brandRepo,
			hasLiked,
			hasDisliked,
		)
	}
	return articles
}

func Brand(b repos.BrandRecord) BrandDescriptor {
	return BrandDescriptor{
		Name: b.Name,
		Logo: b.Logo,
	}
}
