package main

import (
	"encoding/json"
	"fmt"
	//"errors"
	//"io/ioutil"
	//"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/pborman/uuid"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/repos"
	// "github.com/softdev9/trendee-api-master/uploader"
	// "io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//var categories []string = []string{"dresses"}
var colors []string = []string{"17", "19", "3"}
var limit int = 50
var maxSimultaneousConnection = 4

var colorMap map[string]string = map[string]string{
	"1":  "brown",
	"3":  "orange",
	"4":  "yellow",
	"7":  "red",
	"8":  "purple",
	"10": "blue",
	"13": "green",
	"14": "gray",
	"15": "white",
	"16": "black",
	"17": "pink",
	"18": "gold",
	"19": "silver",
	"20": "beige",
}

var categoriesFemale []string = []string{
	//"robe":
	"dresses",
	//"jupe":
	"skirts",
	//"top":
	"womens-tops",
	//"coats" :
	"womens-coats",
	//"veste" :
	"jacket",
	// combinaison
	"womens-suits",
	// Chemises
	"girls-shirts",
	// pantalons
	"womens-pants",
	//shorts
	"sweaters",
	//
	"shoes",
	"boots",
	"womens-sneakers",
	"handbags",
	"womens-eyewear",
	"womens-intimates",
	"swimsuits",
}

var categoriesMale []string = []string{
	"mens-jackets",
	"mens-outerwear",
	"mens-vests",
	"mens-suits",
	"mens-pants",
	"mens-shorts",
	"mens-shirts",
	"mens-tees-and-tshirts",
	"mens-polo-shirts",
	"mens-sweatshirts",
	"mens-sweaters",
	"mens-shoes",
	"mens-shoes-athletic",
	"mens-bags",
	"mens-watches",
	"mens-swimsuits",
	"mens-underwear-and-socks",
}

type ShopSenseImageSizes struct {
	Sizes ShopSenseImageSizeList `json:"sizes"`
}

type ShopSenseMetaData struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
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

type ShopSenseRetailer struct {
	Name string `json:"name"`
}

type ShopSenseResponse struct {
	Products []ShopSenseArticle `json:"products"`
	Meta     ShopSenseMetaData  `json:"metadata"`
}

func (p ShopSenseArticle) toArticleRecord(gender string, cat string, color string) *data.Article {
	return &data.Article{
		Brand: p.Brand.Name,
		Source: data.ArticleSource{
			URL: "http://api.shopstyle.fr",
			Id:  string(p.Id),
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
			strings.ToLower(colorMap[color]),
			strings.ToLower(p.Retailer.Name),
		},
		Retailer:    p.Retailer.Name,
		ClickUrl:    p.ClickUrl,
		Description: p.Description,
	}
}

func (p ShopSenseArticle) Ok() bool {
	return len(p.Brand.Name) > 0
}

func launchShopSensePolling(db *mgo.Session, gender string, cat string, color string, done chan bool) {
	// Create the query plan ie build the different url we need to query to get the maximum number of article
	urls := buildQueryPlan(cat, color)
	//
	remaining := len(urls)
	urlToVisit := 0
	for remaining > 0 {
		numberOfConnectionNeeded := maxSimultaneousConnection
		if remaining < maxSimultaneousConnection {
			numberOfConnectionNeeded = remaining
		}
		cs := make(chan string, numberOfConnectionNeeded)
		log.Println("Number of connection needed ", numberOfConnectionNeeded)
		for i := 0; i < numberOfConnectionNeeded; i++ {
			go getArticle(urls[urlToVisit], gender, cat, color, db, cs)
			urlToVisit = urlToVisit + 1
		}
		for j := 0; j < numberOfConnectionNeeded; j++ {
			log.Println("Done ", <-cs)
		}
		remaining = remaining - numberOfConnectionNeeded
	}
	done <- true
}

func buildQueryPlan(cat string, color string) []string {
	resp, err := http.Get(buildUrl(cat, color, limit, 0))
	if err != nil {
		log.Fatal("Could not contact shop style")
	}
	defer resp.Body.Close()
	var shopSenseResponse ShopSenseResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&shopSenseResponse)
	if err != nil {
		log.Fatal("Could not decode shop style response")
	}
	//numberOfPages := (shopSenseResponse.Meta.Total / limit) + 1
	//if numberOfPages > 10 {
	//	numberOfPages = 1
	//}
	numberOfPages := 1
	urlToQuery := make([]string, 1, 1)
	for page := 0; page < numberOfPages; page++ {
		offset := limit * page
		urlToQuery[page] = buildUrl(cat, color, limit, offset)
	}
	return urlToQuery
}

func buildUrl(cat string, color string, limit int, offset int) string {
	var baseUrl string = "http://api.shopstyle.fr/api/v2/products?pid=uid2401-32446201-84&cat=%s&fl=c%s&limit=%v&offset=%v"
	url := fmt.Sprintf(baseUrl, cat, color, limit, offset)
	return url
}

func getArticle(url string, gender string, cat string, color string, session *mgo.Session, done chan string) {
	log.Printf(" Fetcher launched for %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error 209", err.Error())
		done <- url
		return
	}
	defer resp.Body.Close()
	// Decode the json received
	var shopSenseResponse ShopSenseResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&shopSenseResponse)
	if err != nil {
		log.Println("Error 215", err.Error())
		done <- url
		return
	}
	copyDBSess := session.Copy()
	defer copyDBSess.Close()
	artRepo := repos.NewArticleRepo(copyDBSess)
	brandRepo := repos.NewBrandRepo(copyDBSess)
	// Build the article records
	var arts []interface{}
	for _, p := range shopSenseResponse.Products {
		//if artRepo.IsInShopFromProvider("http://api.shopstyle.fr", p.Id) {
		//	break
		//}
		if p.Ok() {
			// log.Printf("Saving product from url : %s %i, %s \n", url, p.Id, p.Name)
			articleRecord := p.toArticleRecord(gender, cat, color)
			articleRecord.Id = bson.NewObjectId()
			articleRecord.Category = cat
			articleRecord.Store = data.StoreFR
			start := time.Now()
			defineLogoAndStoreBrand(p.Brand.Name, brandRepo)
			//storeDefaultBrand(p.Brand.Name, brandRepo)
			elapsed := time.Since(start)
			log.Printf("[PERF] logo and brand took %s", elapsed)
			articleRecord.Image = buildImageDescriptor(p.ShopSenseImage)
			arts = append(arts, articleRecord)
		}
	}
	if err := artRepo.SaveMultiToStore(arts); err != nil {
		log.Fatal("[ERROR] Could not insert article", err.Error())
		done <- url
		return
	}
	done <- url
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

func storeDefaultBrand(brandN string, brandRepo repos.BrandRepo) {
	if brandRepo.HasBrand(brandN) {
		return
	}
	logo := map[string]string{
		"xsmall": "http://trendee.co/assets/img/favicon/favicon-96x96.png",
		"small":  "http://trendee.co/assets/img/favicon/favicon-96x96.png",
		"medium": "http://trendee.co/assets/img/favicon/favicon-96x96.png",
		"large":  "http://trendee.co/assets/img/favicon/favicon-96x96.png",
		"xlarge": "http://trendee.co/assets/img/favicon/favicon-96x96.png",
	}
	newBrand := &repos.BrandRecord{
		Name: brandN,
		Logo: logo,
	}
	err := brandRepo.SaveBrand(newBrand)
	if err != nil {
		log.Println("Could not save the brand ", err.Error())
	}
}

func defineLogoAndStoreBrand(brandN string, brandRepo repos.BrandRepo) {
	if brandRepo.HasBrand(brandN) {
		return
	}
	name := strings.Replace(strings.ToLower(brandN), " ", "", -1)
	url := fmt.Sprintf("https://logo.clearbit.com/%s.com?format=jpg&size=96", name)
	respImg, err := http.Get(url)
	if err != nil {
		return
	}
	defer respImg.Body.Close()
	var logo map[string]string
	if respImg.StatusCode != http.StatusOK {
		logo = map[string]string{
			"xsmall": "http://trendee.co/assets/img/favicon/favicon-96x96.png",
			"small":  "http://trendee.co/assets/img/favicon/favicon-96x96.png",
			"medium": "http://trendee.co/assets/img/favicon/favicon-96x96.png",
			"large":  "http://trendee.co/assets/img/favicon/favicon-96x96.png",
			"xlarge": "http://trendee.co/assets/img/favicon/favicon-96x96.png",
		}
	} else {
		logo = map[string]string{
			"xsmall": fmt.Sprintf("https://logo.clearbit.com/%s.com?format=jpg&size=32", name),
			"small":  fmt.Sprintf("https://logo.clearbit.com/%s.com?format=jpg&size=96", name),
			"medium": fmt.Sprintf("https://logo.clearbit.com/%s.com?format=jpg&size=128", name),
			"large":  fmt.Sprintf("https://logo.clearbit.com/%s.com?format=jpg&size=256", name),
			"xlarge": fmt.Sprintf("https://logo.clearbit.com/%s.com?format=jpg&size=512", name),
		}
	}
	newBrand := &repos.BrandRecord{
		Name: brandN,
		Logo: logo,
	}
	err = brandRepo.SaveBrand(newBrand)
	if err != nil {
		log.Println("Could not save the brand ", err.Error())
	}
}

func getArticlesForGenderCat(gender string, db *mgo.Session, done chan bool) {
	var cats []string
	if gender == "male" {
		cats = categoriesMale
	} else {
		cats = categoriesFemale
	}
	c := make(chan bool)
	numberOfConnLaunched := 5
	for _, cat := range cats {
		for kCol, _ := range colorMap {
			log.Printf(" Downloading %s %s \n", cat, kCol)
			go launchShopSensePolling(db, gender, cat, kCol, c)
			numberOfConnLaunched = numberOfConnLaunched - 1
			if numberOfConnLaunched == 0 {
				log.Printf("Waiting for one chanel to finish \n")
				<-c
				numberOfConnLaunched = numberOfConnLaunched + 1
			}
		}
	}
	done <- true
}

func main() {
	start := time.Now()
	// Initializing db connection
	//dbUrl := os.Getenv("MONGODB_TEST")
	dbUrl := os.Getenv("MONGODB_PROD")
	log.Println("dbURL : ", dbUrl)
	session, err := mgo.Dial(dbUrl)
	if err != nil {
		log.Panic("[ERROR] Not able to connect to the db", err)
		return
	}
	defer session.Close()
	log.Println("Connected to ", dbUrl)
	// Send a request every hour to get the latest article
	//pollEvery(2*time.Second, session, getArticle
	c := make(chan bool, 2)
	go getArticlesForGenderCat("male", session, c)
	go getArticlesForGenderCat("female", session, c)
	<-c
	<-c
	elapsed := time.Since(start)
	log.Printf("[PERF] feeder took %s", elapsed)
}
