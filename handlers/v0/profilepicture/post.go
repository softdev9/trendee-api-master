package profilepicture

import (
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/softdev9/trendee-api-master/data"
	"github.com/softdev9/trendee-api-master/handlers"
	"github.com/softdev9/trendee-api-master/handlers/v0"
	"github.com/softdev9/trendee-api-master/repos"
	"github.com/softdev9/trendee-api-master/uploader"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type UpadtedProfile struct {
	v0.Response
	ProfilePic map[string]string `json:"profile_picture"`
}

func UploadProfilePic(rw http.ResponseWriter, req *http.Request) {
	u := context.Get(req, "user").(*data.User)
	log.Printf("User : %v ", u)

	resp := &UpadtedProfile{}
	userRepo := context.Get(req, repos.UserR).(repos.ProfilePictureUpdater)
	xlarge := uploader.Image{Name: "xlarge", Width: 1024}
	large := uploader.Image{Name: "large", Width: 768}
	medium := uploader.Image{Name: "medium", Width: 512}
	small := uploader.Image{Name: "small", Width: 256}
	xsmall := uploader.Image{Name: "xsmall", Width: 64}
	// Upload the actual picture
	sourcePart, _, err := req.FormFile("photo")
	if sourcePart == nil {
		resp.Error = handlers.ErrorDescriptor{
			Error:     "No photo sent in the multipart",
			ErrorCode: 101,
		}
		handlers.Respond(rw, req, http.StatusOK, &resp)
		return
	}
	tmpFile, err := copyToTmpDir(sourcePart)
	imgs, err := uploader.UploadImage(tmpFile, uploader.ProfilePicBucket, xlarge, large, medium, small, xsmall)
	err = userRepo.UpdateProfilePicture(u.ID, imgs)
	if err != nil {
		resp.Error = handlers.ErrorDescriptor{
			Error:     "Not able to update the user profile",
			ErrorCode: 102,
		}
		handlers.Respond(rw, req, http.StatusInternalServerError, &resp)
		return
	}
	resp.ProfilePic = imgs
	handlers.Respond(rw, req, http.StatusOK, &resp)
}

func copyToTmpDir(src multipart.File) (*os.File, error) {
	tmpFilePath := uploader.GenerateFilePath()
	tmpFile, err := os.Create(tmpFilePath)
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, src)
	if err != nil {
		log.Printf("[ERROR] L63 while creating file %s", err.Error())
		return nil, err
	}
	return tmpFile, nil
}
