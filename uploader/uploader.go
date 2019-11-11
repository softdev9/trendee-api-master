package uploader

import (
	"errors"
	"fmt"

	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/nfnt/resize"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/github.com/pborman/uuid"
	//"github.com/softdev9/trendee-api-master/data"
	"image/jpeg"

	"github.com/softdev9/trendee-api-master/gateways/simple_uploader"
	// "io"
	"log"
	"os"
)

var tmpPath = os.Getenv("TMP_FOLDER")
var SelfieBucket = "trendee-selfies"
var ProfilePicBucket = "trendee-profile-eu"
var BrandsBucket = "trendee-brands"
var ArticleBucket = "trendee-articles"

type Image struct {
	Name  string
	Width uint
	Url   string
}

func UploadImage(source *os.File, bucket string, sizes ...Image) (map[string]string, error) {
	uploaded := make(map[string]string)
	c := make(chan interface{})
	if source == nil {
		fmt.Println("FILE IS nil")
		return nil, errors.New("not able to get the file")
	}
	for _, size := range sizes {
		go downSampleBy(source.Name(), size, bucket, c)
	}
	err := waitUpload(c, uploaded, len(sizes))
	return uploaded, err
}

func downSampleBy(path string, target Image, bucket string, c chan interface{}) {
	// Open the path
	log.Printf("[DEBUG] Opening file at %s", path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("err 85: ", err.Error())
		c <- err
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("err 57: ", err.Error())
		c <- err
		return
	}
	filePath := GenerateFilePath()
	m := resize.Resize(target.Width, 0, img, resize.Lanczos3)
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println("err 98: ", err.Error())
		c <- err
		return
	}
	defer out.Close()
	// write new image to file
	jpeg.Encode(out, m, nil)
	// Upload the file
	url, err := simple_uploader.UploadFile(out, bucket)
	if err != nil {
		fmt.Println("err 107: ", err.Error())
		c <- err
		return
	}
	target.Url = url
	os.Remove(out.Name())
	c <- target
}

func GenerateFilePath() string {
	// Create a  file with a uuid.jpg
	fileName := uuid.New() + ".jpg"
	filePath := fmt.Sprintf("%s/%s", tmpPath, fileName)
	return filePath
}

func waitUpload(c chan interface{}, uploaded map[string]string, uploadeToWait int) error {
	for uploadeToWait > 0 {
		switch r := <-c; r.(type) {
		case Image:
			img := r.(Image)
			uploaded[img.Name] = img.Url
			uploadeToWait--
		case error:
			err := r.(error)
			uploadeToWait = 0
			return err
		default:
			panic("Unexpected type")
		}
	}
	return nil

}
