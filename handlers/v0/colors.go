package v0

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/repos"
	"net/http"
)

type ColorListResponse struct {
	Response
	Colors []ColorDescriptor `json:"colors"`
}

type ColorDescriptor struct {
	Id      string `json:"id"`
	HexCode string `json:"hexcode"`
}

func GetColorList(rw http.ResponseWriter, req *http.Request) {
	resp := ColorListResponse{}
	colorRepo := context.Get(req, repos.ColorsR).(*repos.ColorRepoMGO)
	colors, err := colorRepo.GetAllColors()
	if err != nil {
		errDesc := handlers.ErrorDescriptor{Error: "Not Able to get the color list" + err.Error(), ErrorCode: 100}
		resp.Response.Error = errDesc
	}
	if colors != nil {
		resp.Colors = make([]ColorDescriptor, len(colors), len(colors))
		for i, c := range colors {
			resp.Colors[i] = ColorDescriptor{HexCode: c.HexCode, Id: c.Name}
		}
	}
	handlers.Respond(rw, req, http.StatusOK, &resp)
}
