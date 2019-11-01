package selfie

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/gateways/partners"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/handlers/v0/descriptor"
	"github.com/softdev9/trendee-api-master/repos"
	"github.com/softdev9/trendee-api-master/uploader"
)

type UploadResponse struct {
	v0.Response
	// Error
	// Selfie
	// Similar articles
	Selfie   descriptor.SelfieDescriptor    `json:"selfie"`
	Articles []descriptor.ArticleDescriptor `json:"articles"`
}

type TagString string

type IncomingTags struct {
	Comment string     `json:"comment"`
	Tags    []data.Tag `json:"tags"`
}

func (t TagString) Parse() (*IncomingTags, error) {
	// Create a new json decoder
	decoder := json.NewDecoder(strings.NewReader(string(t)))
	var tags IncomingTags
	if err := decoder.Decode(&tags); err != nil {
		log.Println("[ERROR] Not able to create the tag list ", err.Error())
		return nil, err
	}
	return &tags, nil
}

func UploadSelfie(rw http.ResponseWriter, req *http.Request) {
	u := context.Get(req, "user").(*data.User)
	userRepo := context.Get(req, repos.UserR).(repos.UserRepo)
	selfieRepo := context.Get(req, repos.SelfieR).(repos.SelfieRepo)
	articleRepo := context.Get(req, repos.ArticleR).(repos.ArticleRepo)
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	var uploadResponse UploadResponse
	imgs, err := uploadPicture(req)
	if err != nil {
		uploadResponse.Error = handlers.ErrorDescriptor{
			Error:     fmt.Sprintf("60 Error selfie upload imposible Reason %s", err.Error()),
			ErrorCode: 101,
		}
		handlers.Respond(rw, req, http.StatusOK, &uploadResponse)
		return
	}
	tags, err := TagString(req.FormValue("selfie_data")).Parse()
	if err != nil {
		uploadResponse.Error = handlers.ErrorDescriptor{
			Error:     fmt.Sprintf("69 Unable to decode tags list %s The tags list should be sent with paremater name seflie_data", err.Error()),
			ErrorCode: 101,
		}
		handlers.Respond(rw, req, http.StatusOK, &uploadResponse)
		return
	}
	selfie, err := data.NewSelfie(u.ID, tags.Comment, tags.Tags, imgs)
	if err != nil {
		uploadResponse.Error = handlers.ErrorDescriptor{
			Error:     fmt.Sprintf("78 Error selfie upload imposible Reason %s", err.Error()),
			ErrorCode: 101,
		}
		handlers.Respond(rw, req, http.StatusOK, &uploadResponse)
		return
	}
	comStoreIds := partners.GetArticleForTagList(articleRepo, brandRepo, data.LANG_EN, selfie.Tags)
	frStoreIds := partners.GetArticleForTagList(articleRepo, brandRepo, data.LANG_FR, selfie.Tags)
	fullIds := append(frStoreIds, comStoreIds...)
	selfieRepo.SaveSelfie(selfie)
	saveTagBrands(selfie.Tags, brandRepo)
	selfieRepo.Similars(selfie.ID, fullIds)
	sr, _ := selfieRepo.GetSelfieById(selfie.ID)
	selfieD, articles := LoadDataForSelfie(u, sr, selfieRepo, userRepo, articleRepo, brandRepo)
	uploadResponse.Selfie = selfieD
	uploadResponse.Articles = articles
	handlers.Respond(rw, req, http.StatusOK, &uploadResponse)
}

func saveTagBrands(tags []data.Tag, repo repos.BrandRepo) {
	for _, t := range tags {
		repo.SaveBrand(&repos.BrandRecord{
			Name: t.Brand,
		})
	}
}

func uploadPicture(req *http.Request) (map[string]string, error) {
	sourcePart, _, err := req.FormFile("photo")
	log.Println("req From Value for photo : ", len(req.FormValue("photo")))
	if sourcePart == nil {
		return nil, errors.New("No file sent " + err.Error())
	}

	tmpFile, err := copyToTmpDir(sourcePart)
	if err != nil {
		log.Printf("Unable to open the form file %s", err.Error())
		return nil, err
	}
	xlarge := uploader.Image{Name: "xlarge", Width: 2048}
	large := uploader.Image{Name: "large", Width: 1024}
	medium := uploader.Image{Name: "medium", Width: 512}
	small := uploader.Image{Name: "small", Width: 256}
	xsmall := uploader.Image{Name: "xsmall", Width: 128}
	imgs, err := uploader.UploadImage(tmpFile, uploader.SelfieBucket, xlarge, large, medium, small, xsmall)
	if err != nil {
		log.Printf("[ERROR] Not able to upload the image %s", err.Error())
		return nil, err
	}
	log.Printf("Imgs %v", imgs)
	err = os.Remove(tmpFile.Name())
	if err != nil {
		log.Printf("[ERROR] Not able to remove the file %s", err.Error())
	}
	return imgs, nil
}

func copyToTmpDir(src multipart.File) (*os.File, error) {
	tmpFilePath := uploader.GenerateFilePath()
	tmpFile, err := os.Create(tmpFilePath)
	log.Printf("[DEBUB] 139 Path for temp file : %s", tmpFilePath)
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, src)
	if err != nil {
		log.Printf("[ERROR] L144 while creating temp file : %s", err.Error())
		return nil, err
	}
	return tmpFile, nil
}
