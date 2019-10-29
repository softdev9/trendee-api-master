package partners

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
)

const ShopStyleApiKey string = "uid2401-32446201-84"

const ShopStyleBaseUrlFr string = "api.shopstyle.fr"

const ShopStyleBaseUrlCom string = "api.shopstyle.com"

const ShopStyleSearchPath string = "/api/v2/products"

type ShopSenseResponse struct {
	Products []ShopSenseArticle `json:"products"`
	Meta     ShopSenseMetaData  `json:"metadata"`
}

type ShopSenseImageSizes struct {
	Sizes ShopSenseImageSizeList `json:"sizes"`
}

type ShopSenseMetaData struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

var subcatNeedingParent map[string]bool = map[string]bool{
	"sleeveless":           true,
	"longsleeve":           true,
	"printed":              true,
	"cashmere":             true,
	"halter":               true,
	"midi":                 true,
	"maxi":                 true,
	"denim":                true,
	"leather":              true,
	"faux fur":             true,
	"suede":                true,
	"velvet":               true,
	"fur & shearling":      true,
	"leather & suede":      true,
	"wool":                 true,
	"crewneck & scoopneck": true,
	"turtleneck":           true,
	"v-neck":               true,
	"athletic":             true,
	"dress":                true,
}

type ShopSenseImageSizeList struct {
	Small       ShopSenseImageSize `json:"small"`
	XLarge      ShopSenseImageSize `json:"XLarge"`
	Medium      ShopSenseImageSize `json:"Medium"`
	Large       ShopSenseImageSize `json:"Large"`
	IPhoneSmall ShopSenseImageSize `json:"IPhoneSmall"`
	Best        ShopSenseImageSize `json:"Best"`
	Original    ShopSenseImageSize `json:"Original"`
	IPhone      ShopSenseImageSize `json:"IPhone"`
}

type ShopSenseBrand struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ShopSenseImageSize struct {
	SizeName string `json:"sizeName"`
	Url      string `json:"url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type ShopSenseArticle struct {
	Id             int                 `json:"id"`
	Name           string              `json:"name"`
	Currency       string              `json:"currency"`
	Price          float32             `json:"price"`
	InStock        bool                `json:"inStock"`
	ShopSenseImage ShopSenseImageSizes `json:"image"`
	Brand          ShopSenseBrand      `json:"brand"`
	Retailer       ShopSenseRetailer   `json:"retailer"`
	ClickUrl       string              `json:"clickUrl"`
	Description    string              `json:"description"`
}

func buildImageDescriptor(image ShopSenseImageSizes) map[string]string {
	return map[string]string{
		"xlarge": image.Sizes.XLarge.Url,
		"large":  image.Sizes.Large.Url,
		"medium": image.Sizes.Medium.Url,
		"small":  image.Sizes.IPhoneSmall.Url,
		"xsmall": image.Sizes.Small.Url,
	}
}

func (p ShopSenseArticle) toArticleRecord(url string, gender string, cat string, subCat string, color string) *data.Article {
	log.Printf("Brand for article : %s \n", p.Brand.Name)
	return &data.Article{
		Brand: p.Brand.Name,
		Source: data.ArticleSource{
			URL: url,
			Id:  strconv.Itoa(p.Id),
		},
		Gender:   gender,
		Name:     p.Name,
		InStock:  p.InStock,
		Price:    p.Price,
		Currency: p.Currency,
		Keywords: []string{
			strings.ToLower(gender),
			strings.ToLower(p.Brand.Name),
			strings.ToLower(cat),
			strings.ToLower(subCat),
			strings.ToLower(color),
			strings.ToLower(p.Retailer.Name),
		},
		Image:       buildImageDescriptor(p.ShopSenseImage),
		Retailer:    p.Retailer.Name,
		ClickUrl:    p.ClickUrl,
		Description: p.Description,
	}
}

type ShopSenseRetailer struct {
	Name string `json:"name"`
}

func buildQueryUrl(baseUrl string, gender string, cat string, subCat string, color string) string {
	url := &url.URL{}
	url.Scheme = "http"
	url.Host = baseUrl
	url.Path = ShopStyleSearchPath
	q := url.Query()
	q.Set("fts", gender+" "+subCat+" "+color)
	q.Set("pid", ShopStyleApiKey)
	url.RawQuery = q.Encode()
	return url.String()
}

func buildQueryUrl2(baseUrl string, gender string, cat string, subCat string, color string) string {
	url := &url.URL{}
	url.Scheme = "http"
	url.Host = baseUrl
	url.Path = ShopStyleSearchPath
	q := url.Query()
	q.Set("fts", gender+" "+cat+" "+color)
	q.Set("pid", ShopStyleApiKey)
	url.RawQuery = q.Encode()
	return url.String()
}

func buildQueryUrl3(baseUrl string, gender string, cat string, subCat string, color string) string {
	url := &url.URL{}
	url.Scheme = "http"
	url.Host = baseUrl
	url.Path = ShopStyleSearchPath
	q := url.Query()
	q.Set("fts", gender+" "+cat)
	q.Set("pid", ShopStyleApiKey)
	url.RawQuery = q.Encode()
	return url.String()
}

func buildComUrl(t data.Tag) string {
	return buildQueryUrl(
		ShopStyleBaseUrlCom,
		string(t.Gender),
		t.Category,
		t.SubCategory,
		t.Color,
	)
}

func buildUrl(formulaNumber int, lang string, t data.Tag) string {
	gender := string(t.Gender)
	cat := t.Category
	subCat := t.SubCategory
	color := t.Color
	baseUrl := ShopStyleBaseUrlCom
	if lang == data.LANG_FR {
		baseUrl = ShopStyleBaseUrlFr
		gender = data.Gender(gender).TranslateIn(data.LANG_FR)
		cat = data.Category(t.Category).TranslateIn(data.LANG_FR)
		subCat = data.SubCategory(t.SubCategory).TranslateIn(data.LANG_FR)
		color = data.Color(t.Color).TranslateIn(data.LANG_FR)
	}
	if formulaNumber == 0 {
		// Need to specify the cat for some subcat
		val, ok := subcatNeedingParent[t.SubCategory]
		log.Printf("Subcat : %s val %v -> found %v", t.SubCategory, val, ok)
		if ok {
			subCat = cat + " " + subCat
		}

		return buildQueryUrl(baseUrl, string(gender), cat, subCat, color)
	}
	if formulaNumber == 1 {
		return buildQueryUrl2(baseUrl, string(gender), cat, subCat, color)
	}
	return buildQueryUrl3(baseUrl, string(gender), cat, subCat, color)
}

func GetArticleForTagList(artRepo repos.ArticleRepo, brandRepo repos.BrandRepo, lang string, tags []data.Tag) []bson.ObjectId {
	var idList []bson.ObjectId
	for i, tag := range tags {
		formulaTried := 0
		for len(idList) < 10*(i+1) && formulaTried < 2 {
			url := buildUrl(formulaTried, lang, tag)
			log.Printf("[DEBUG] url = %s", url)
			shopSenseResponse := launchShopStyleQuery(url)
			list := extractIdsToList(tag, shopSenseResponse.Products, lang, artRepo, brandRepo)
			idList = append(idList, list...)
			formulaTried++
		}
	}
	return idList
}

func extractIdsToList(tag data.Tag, articles []ShopSenseArticle, lang string, artRepo repos.ArticleRepo, brandRepo repos.BrandRepo) []bson.ObjectId {
	var idList []bson.ObjectId
	for _, shopSenseArticle := range articles {
		// Convert the article to a trendee article
		baseUrl := ShopStyleBaseUrlCom
		store := data.StoreCOM
		brandRepo.SaveBrand(&repos.BrandRecord{Name: shopSenseArticle.Brand.Name})
		if lang == data.LANG_FR {
			baseUrl = ShopStyleBaseUrlFr
			store = data.StoreFR
		}
		art := shopSenseArticle.toArticleRecord(
			baseUrl,
			string(tag.Gender),
			tag.Category,
			tag.SubCategory,
			tag.Color,
		)
		art.Store = store
		// Insert the article to the good store or retreive the id of the article
		err := artRepo.SaveArticleToStore(art)
		if err != nil {
			log.Printf("[ERROR] inserting %s ", err.Error())
		}
		// Add the article id to the list of article
		idList = append(idList, art.Id)
	}
	return idList
}

// Launch the query to shopstyle with the given url
func launchShopStyleQuery(url string) *ShopSenseResponse {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var shopSenseResponse ShopSenseResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&shopSenseResponse)
	if err != nil {
		log.Printf("[FATAL] 224.Could not decode shop style response from url %s,\n\n%s ", url, err.Error())
		log.Printf("[DEBUG] %s", resp.Body)
		return nil
	}
	return &shopSenseResponse
}

// Launch a get request to the url
// the url should be a shop style api url
func launchQuery(url string) ShopSenseResponse {
	resp, err := http.Get(url)
	var shopSenseResponse ShopSenseResponse
	if err != nil {
		log.Fatal("[ERROR] Could not get in touch with shopstyle")
		return shopSenseResponse
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&shopSenseResponse)
	if err != nil {
		log.Fatal("[ERROR] Could not decode shop style response")
		return shopSenseResponse
	}
	return shopSenseResponse
}
