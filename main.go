package main

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/api/vision/v1"
)

func main() {
	confFile, err := ioutil.ReadFile("google_credentials.json")
	if err != nil {
		log.Fatalln("failed to read configuration file", err)
	}

	cfg, err := google.JWTConfigFromJSON([]byte(confFile), vision.CloudPlatformScope)
	client := cfg.Client(context.Background())
	svc, err := vision.New(client)

	// Possible values:
	//   "TYPE_UNSPECIFIED" - Unspecified feature type.
	//   "FACE_DETECTION" - Run face detection.
	//   "LANDMARK_DETECTION" - Run landmark detection.
	//   "LOGO_DETECTION" - Run logo detection.
	//   "LABEL_DETECTION" - Run label detection.
	//   "TEXT_DETECTION" - Run OCR.
	//   "SAFE_SEARCH_DETECTION" - Run various computer vision models to
	//   "IMAGE_PROPERTIES" - compute image safe-search properties.
	features := make([]*vision.Feature, 2)
	features[0] = &vision.Feature{
		Type:       "LABEL_DETECTION",
		MaxResults: 5,
	}
	features[1] = &vision.Feature{
		Type:       "TEXT_DETECTION",
		MaxResults: 5,
	}

	requests := make([]*vision.AnnotateImageRequest, 2)
	images := []string{"images/IMGP8896.JPG", "images/IMGP9629.JPG"}

	for i := 0; i < len(images); i++ {
		imgData, _ := ioutil.ReadFile(images[i])
		enc := base64.StdEncoding.EncodeToString([]byte(imgData))
		img := &vision.Image{Content: enc}
		requests[i] = &vision.AnnotateImageRequest{
			Image:    img,
			Features: features,
		}
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: requests,
	}
	res, err := svc.Images.Annotate(batch).Do()

	var result string

	for i := 0; i < len(images); i++ {
		body, err := json.MarshalIndent(res.Responses[i], "", "\t")
		if err != nil {
			fmt.Println("Responses error")
			fmt.Println(err)
		}
		result += string(body)
	}
	fmt.Println(result)
}
