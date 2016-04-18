package main

import (
	"errors"
	"flag"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"google.golang.org/api/vision/v1"
)

func generateFeatures(typeStr string) ([]*vision.Feature, error) {
	types := strings.Split(typeStr, ",")
	features := make([]*vision.Feature, len(types))
	var featureType string

	// Possible values:
	//   "TYPE_UNSPECIFIED" - Unspecified feature type.
	//   "FACE_DETECTION" - Run face detection.
	//   "LANDMARK_DETECTION" - Run landmark detection.
	//   "LOGO_DETECTION" - Run logo detection.
	//   "LABEL_DETECTION" - Run label detection.
	//   "TEXT_DETECTION" - Run OCR.
	//   "SAFE_SEARCH_DETECTION" - Run various computer vision models to
	//   "IMAGE_PROPERTIES" - compute image safe-search properties.
	for i := 0; i < len(types); i++ {
		switch types[i] {
		case "face":
			featureType = "FACE_DETECTION"
		case "landmark":
			featureType = "LANDMARK_DETECTION"
		case "logo":
			featureType = "LOGO_DETECTION"
		case "label":
			featureType = "LABEL_DETECTION"
		case "text":
			featureType = "TEXT_DETECTION"
		case "safe_search":
			featureType = "SAFE_SEARCH_DETECTION"
		case "image_properties":
			featureType = "IMAGE_PROPERTIES"
		default:
			errorMsg := "Invalid feature: " + types[i]
			return nil, errors.New(errorMsg)
		}

		features[i] = &vision.Feature{
			Type:       featureType,
			MaxResults: 5,
		}
	}
	return features, nil
}

func main() {
	confFile, err := ioutil.ReadFile("google_credentials.json")
	if err != nil {
		fmt.Println("Failed to read credential file. Please add 'google_credentials.json'", err)
		return
	}

	var typeStr = flag.String("feature", "label,text", "face, landmark, logo, label, text, safe_search, image_properties")
	var image = flag.String("image", "", "image file path")
	flag.Parse()
	if *image == "" {
		fmt.Println("Please specify image file name")
		return
	}

	cfg, err := google.JWTConfigFromJSON([]byte(confFile), vision.CloudPlatformScope)
	client := cfg.Client(context.Background())
	svc, err := vision.New(client)

	features, err := generateFeatures(*typeStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	requests := make([]*vision.AnnotateImageRequest, 1)
	images := []string{*image}

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
