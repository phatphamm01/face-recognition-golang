package recognizeRouter

import (
	"context"
	"face-recognition-golang/constants/path"
	"face-recognition-golang/db"
	"face-recognition-golang/db/collection"
	"face-recognition-golang/libs/faceCustom"
	"face-recognition-golang/libs/recognizer"
	commonUtils "face-recognition-golang/utils/common"
	"face-recognition-golang/validator"
	"fmt"
	"image"
	"strings"

	"github.com/fogleman/gg"
	"github.com/samber/lo"
)

func VerifiedHandler(DB *db.MongoDB, input *validator.ValidateRecognize) *string {
	cursor, err := DB.DatasetCollectionName.Find(context.Background(), map[string]interface{}{
		"UserID": input.UserID,
	})

	if err != nil {
		return &input.Image
	}

	var results []collection.Dataset
	if err = cursor.All(context.TODO(), &results); err != nil {
		return &input.Image
	}

	facesDescriptor := lo.Map(results, func(item collection.Dataset, _ int) recognizer.Data {
		return recognizer.Data{
			Descriptor: item.Descriptor,
			Id:         item.UserID,
		}
	})
	b64data := strings.Split(input.Image, "base64,")[1]
	isAuth, err := faceCustom.RecognizeFace(facesDescriptor, b64data)

	fmt.Println(isAuth)
	if err != nil {
		return &input.Image
	}

	if isAuth != nil && *isAuth {
		imageVerified, err := AddVerifiedImage(b64data)
		if err != nil {
			return &input.Image
		}

		base64Image, err := commonUtils.ImageToBase64(imageVerified)
		if err != nil {
			return &input.Image
		}

		result := "data:image/jpeg;base64," + base64Image
		return &result
	}

	return &input.Image
}

func AddVerifiedImage(images string) (image.Image, error) {
	var (
		VerifiedImagePath = path.GetBasepath() + "/public/verified.jpg"
	)

	img, err := commonUtils.Base64ToImage(images)
	if err != nil {
		return nil, err
	}
	Height := img.Bounds().Max.Y
	Width := img.Bounds().Max.X

	dc := gg.NewContext(Width, Height)
	dc.DrawImage(img, 0, 0)

	verified, err := GetImageFromPath(VerifiedImagePath)
	if err != nil {
		panic(err)
	}

	verifiedHeight := (*verified).Bounds().Max.Y
	verifiedWidth := (*verified).Bounds().Max.X

	verifiedDc := gg.NewContext(verifiedWidth, verifiedHeight)
	verifiedDc.Scale(0.86, 0.86)
	verifiedDc.DrawImage(*verified, 0, 0)

	dc.DrawImage(verifiedDc.Image(), Width-340, Height-80)

	return dc.Image(), nil
}
