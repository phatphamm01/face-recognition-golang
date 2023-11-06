package recognizeRouter

import (
	"face-recognition-golang/db"
	commonUtils "face-recognition-golang/utils/common"
	"face-recognition-golang/validator"
	"testing"
)

func Test(t *testing.T) {
	DB, err := db.NewMongoDB()
	if err != nil {
		t.Error(err)
	}

	imageBase64, err := commonUtils.ReadImageFromFile("./image.jpg")
	if err != nil {
		t.Error(err)
	}

	image := VerifiedHandler(DB, &validator.ValidateRecognize{
		UserID: "65039bc7c0691e1ea988dd05",
		Image:  *imageBase64,
	})

	if err = commonUtils.WriteImageToFile("./image_verified.jpg", *image); err != nil {
		t.Error(err)
	}
}
