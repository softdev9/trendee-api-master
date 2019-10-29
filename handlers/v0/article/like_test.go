package article

/*
import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)



func TestLikeUnLiked(t *testing.T) {
	u, art := initTestData()
	artRepo := newInMemArticleRepo()
	artRepo.SaveArticleToStore("fake", art)
	// Build the form that is going to be sent
	form := buildForm(art.Id.Hex(), "like")
	// Create the http request to test the endpoint
	req := setupReq(form, u, artRepo)
	// Test the like function
	rw := httptest.NewRecorder()
	PostLikeArticle(rw, req)
	// Check the response
	decoder := json.NewDecoder(rw.Body)
	resp := LikeResponse{}
	if err := decoder.Decode(&resp); err != nil {
		t.Errorf("Could not decode the answer ", err)
	}
	inList := u.HasLiked(art.Id)
	if !inList {
		t.Error("The article has not beed added in the list")
	}
	if resp.Art.NumberOfLike != 1 {
		t.Error("The number of article should be 1 but is ", resp.Art.NumberOfLike)
	}
	if resp.Art.Liked != true {
		t.Error("Liked should be true")
	}
	// Send the unlike
	form = buildForm(art.Id.Hex(), "unlike")
	req = setupReq(form, u, artRepo)
	PostLikeArticle(rw, req)
	// Check the response
	decoder = json.NewDecoder(rw.Body)
	resp = LikeResponse{}
	if err := decoder.Decode(&resp); err != nil {
		t.Errorf("Could not decode the answer ", err)
	}
	inList = u.HasLiked(art.Id)
	if resp.Error.ErrorCode != 0 {
		t.Error("The error coode should be 0 but is ", resp.Error.ErrorCode)
	}
	if inList {
		t.Error("The article should not be in the list once it has been unliked")
	}
	if resp.Art.NumberOfLike != 0 {
		t.Error("The number of article should be 0 but is ", resp.Art.NumberOfLike)
	}

}

func initTestData() (*data.User, *repos.ArticleRecord) {
	u := &data.User{
		Firstname: "kevin",
		Lastname:  "le goff",
		Email:     "kev.legoff@gmail.com",
	}
	// Insert an article
	art := &repos.ArticleRecord{Name: "test"}
	return u, art
}

func buildForm(artId string, operation string) url.Values {
	form := url.Values{}
	form.Add("operation", operation)
	form.Add("article", artId)
	return form
}

func setupReq(form url.Values, user *data.User, artRepo *InMemArticleRepo) *http.Request {
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	brandRepo := InMemBrand{}
	userRepo := &InMemoryArticleLiker{}
	context.Set(req, repos.ArticleR, artRepo)
	context.Set(req, repos.BrandR, brandRepo)
	context.Set(req, repos.UserR, userRepo)
	context.Set(req, "user", user)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

type InMemoryArticleLiker struct{}

func (repo *InMemoryArticleLiker) MarkArticle(user *data.User, articleId bson.ObjectId, operation string) (repos.ArticleLikedChanged, error) {
	user.ArticleLiked = append(user.ArticleLiked, articleId)
	if operation == "like" {
		return repos.ArticleLikedChanged{
			ArticleId:     articleId,
			IncNumLike:    1,
			IncNumDislike: 0,
		}, nil
	}
	if operation == "unlike" {
		// Remove the article id of the Article Liked
		user.ArticleLiked = make([]bson.ObjectId, 1, 1)
		return repos.ArticleLikedChanged{
			ArticleId:     articleId,
			IncNumLike:    -1,
			IncNumDislike: 0,
		}, nil
	}
	return repos.ArticleLikedChanged{
		ArticleId:     articleId,
		IncNumLike:    0,
		IncNumDislike: 0,
	}, nil
}

// A in memory Article repository to test the handlers
type InMemArticleRepo struct {
	articles map[string]repos.ArticleRecord
}

type InMemBrand struct {
	brand []repos.BrandRecord
}

func (repo InMemBrand) FindByName(string) (repos.BrandRecord, error) {
	b := repos.BrandRecord{
		Name: "test_brand",
	}
	return b, nil
}

func newInMemArticleRepo() *InMemArticleRepo {
	return &InMemArticleRepo{
		articles: make(map[string]repos.ArticleRecord),
	}
}

func (r *InMemArticleRepo) SaveArticleToStore(storeName string, article *repos.ArticleRecord) error {
	article.Id = bson.NewObjectId()
	r.articles[article.Id.Hex()] = *article
	return nil
}

func (r *InMemArticleRepo) GetArticleWithId(articleId string) (repos.ArticleRecord, error) {
	art, ok := r.articles[articleId]
	if !ok {
		return repos.ArticleRecord{}, errors.New("Not found")
	}
	return art, nil
}

func (r *InMemArticleRepo) GetArticleNotIds(articleIds []bson.ObjectId) ([]repos.ArticleRecord, error) {
	return nil, nil
}
func (r *InMemArticleRepo) UpdateNumberOfLike(change repos.ArticleLikedChanged) (*repos.ArticleRecord, error) {
	a, _ := r.GetArticleWithId(change.ArticleId.Hex())
	a.NumberOfLike = a.NumberOfLike + change.IncNumLike
	a.NumberOfDislike = a.NumberOfLike + change.IncNumDislike
	r.articles[change.ArticleId.Hex()] = a
	return &a, nil
}

func (r *InMemArticleRepo) CountArticleInStore(storeName string) (int, error) {
	return len(r.articles), nil
}

func (r *InMemArticleRepo) ClearStore(storeName string) {
	// DO NOTHING HERE FOR NOEW
}

func (r *InMemArticleRepo) GetArticleInStore(storeName string, idObjectToQuery string) ([]repos.ArticleRecord, error) {
	// NOT USED SO LETS DO NOTHING HERE
	return nil, nil
}

func (r *InMemArticleRepo) GetArticleInList(ids []bson.ObjectId) ([]repos.ArticleRecord, error) {
	// NOT USED SO LETS DO NOTHING HERE
	return nil, nil
}
*/
