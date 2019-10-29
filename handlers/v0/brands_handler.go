package v0

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/repos"
	"log"
	"net/http"
)

type BrandResp struct {
	Response
	Brands []data.PublicBrand `json:"brands"`
}

func GetBrandList(rw http.ResponseWriter, req *http.Request) {
	brandRepo := context.Get(req, repos.BrandR).(repos.BrandRepo)
	brands, err := brandRepo.GetAllBrands()
	response := BrandResp{}
	if err != nil {
		response.Response.Error = handlers.ErrorDescriptor{
			ErrorCode: 100,
			Error:     "Not able to get the list brand" + err.Error(),
		}
		handlers.Respond(rw, req, http.StatusInternalServerError, &response)
		return
	}
	response.Brands = make([]data.PublicBrand, len(brands), len(brands))
	log.Printf("%d brands found", len(brands))
	for i, b := range brands {
		response.Brands[i] = b.Public()
	}
	handlers.Respond(rw, req, http.StatusOK, &response)
}
