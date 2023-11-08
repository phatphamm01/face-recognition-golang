package recognizeRouter

import (
	"context"
	"face-recognition-golang/db"
	"face-recognition-golang/db/collection"
	"face-recognition-golang/libs/faceCustom"
	"face-recognition-golang/libs/recognizer"
	commonUtils "face-recognition-golang/utils/common"
	"face-recognition-golang/utils/promise"
	"face-recognition-golang/validator"
	"fmt"
	"image"
	"net/http"

	"github.com/samber/lo"
)

func GetImageFromURL(url string) (*image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

type Result struct {
	IsAuth bool
	Image  string
}

func VerifiedHandler(DB *db.MongoDB, input *validator.ValidateRecognize) ([]Result, error) {
	cursor, err := DB.DatasetCollectionName.Find(context.Background(), map[string]interface{}{
		"UserID": input.UserID,
	})

	if err != nil {
		return nil, err
	}

	var datasets []collection.Dataset
	if err = cursor.All(context.TODO(), &datasets); err != nil {
		return nil, err
	}

	facesDescriptor := lo.Map(datasets, func(item collection.Dataset, _ int) recognizer.Data {
		return recognizer.Data{
			Descriptor: item.Descriptor,
			Id:         item.UserID,
		}
	})

	var results []Result
	promise.ParallelWithArr(lo.Map(input.Images, func(item string, _ int) func() {
		return func() {
			imageFromUrl, err := GetImageFromURL(item)
			if err != nil {
				results = append(results, Result{
					Image: item,
				})
				return
			}

			b64data, err := commonUtils.ImageToBase64(*imageFromUrl)
			if err != nil {
				results = append(results, Result{
					Image: item,
				})
				return
			}

			data, err := faceCustom.RecognizeFace(facesDescriptor, b64data)
			if err != nil {
				results = append(results, Result{
					Image: item,
				})
				return
			}

			fmt.Println("data::", *data)

			results = append(results, Result{
				IsAuth: *data,
				Image:  item,
			})
		}
	}))

	return results, nil
}
