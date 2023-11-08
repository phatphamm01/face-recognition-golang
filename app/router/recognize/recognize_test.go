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

	images, err := VerifiedHandler(DB, &validator.ValidateRecognize{
		UserID: "652933ebbdce4daaad99d5af",
		Images: []string{
			"http://res.cloudinary.com/finder-next/image/upload/v1699450821/finder/ums1wqt5fymyu2tadc8l.jpg",
		},
	})

	commonUtils.PrintJSON(images)

}
