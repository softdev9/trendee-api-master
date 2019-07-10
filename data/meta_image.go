package data

import (
	"image/jpeg"
	"log"
	"net/http"
)

type AspectRatio string

type MetaImage struct {
	AspectRatio AspectRatio          `json:"suggested_ratio"`
	ImageSizes  map[string]ImageSize `json:"image_sizes"`
}

type ImageSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type metaPart struct {
	KeySize   string
	ImageSize ImageSize
}

//
func MetaForImage(img map[string]string, metaChan chan MetaImage) {
	meta := MetaImage{ImageSizes: make(map[string]ImageSize)}
	imgParts := 5
	ch := make(chan metaPart, imgParts)
	go MetaForKey("xlarge", img["xlarge"], ch)
	go MetaForKey("large", img["large"], ch)
	go MetaForKey("medium", img["medium"], ch)
	go MetaForKey("small", img["small"], ch)
	go MetaForKey("xsmall", img["xsmall"], ch)
	for imgParts > 0 {
		metaPart := <-ch
		meta.ImageSizes[metaPart.KeySize] = metaPart.ImageSize
		imgParts--
		log.Printf("Img  part remaingin %d ", imgParts)
	}
	ratio := float32(meta.ImageSizes["xlarge"].Height) / float32(meta.ImageSizes["xlarge"].Width)
	meta.AspectRatio = AspectRatio("square")
	if ratio > 1.30 {
		meta.AspectRatio = AspectRatio("portrait")
	}
	if ratio < 0.70 {
		meta.AspectRatio = AspectRatio("landscape")
	}
	metaChan <- meta
}

func MetaForKey(key string, url string, c chan metaPart) {
	log.Printf("Getting meta info for %s at %s ", key, url)
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := check.Get(url) // add a filter to check redirect
	defer resp.Body.Close()
	log.Printf("response status %d, for  %s", resp.StatusCode, url)
	imgJpg, err := jpeg.Decode(resp.Body)

	if err != nil {
		log.Printf("Error while generating meta data ", err.Error())
		//panic(err)
		c <- metaPart{
			KeySize: key,
			ImageSize: ImageSize{
				Width:  800,
				Height: 800,
			},
		}
		return
	}

	c <- metaPart{
		KeySize: key,
		ImageSize: ImageSize{
			Width:  imgJpg.Bounds().Dx(),
			Height: imgJpg.Bounds().Dy(),
		},
	}
}
