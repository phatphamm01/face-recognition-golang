package socketRouter

import (
	"bytes"
	"face-recognition-golang/constants/path"
	commonUtils "face-recognition-golang/utils/common"
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

type DetectFace struct {
	Image     string  `json:"image"`
	FaceTotal int     `json:"faceTotal"`
	FaceImage *string `json:"faceImage"`
}

func MatToString(mat gocv.Mat) (*string, error) {
	buf, err := mat.ToImage()
	if err != nil {
		return nil, err
	}
	resultImage, err := commonUtils.ImageToBase64(buf)
	if err != nil {
		return nil, err
	}
	return &resultImage, nil
}

func DetectAndDrawRectangles(inputImage []byte) (*DetectFace, error) {
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(fmt.Sprintf("%s/public/models/haarcascade_frontalface_alt.xml", path.GetBasepath())) {
		return nil, fmt.Errorf("Error reading cascade file")
	}

	img, _, err := image.Decode(bytes.NewReader(inputImage))
	if err != nil {
		return nil, err
	}

	mat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return nil, err
	}
	defer mat.Close()

	rects := classifier.DetectMultiScale(mat)

	var croppedImageBytes *string
	if len(rects) == 1 {
		rect := rects[0]
		croppedMat := mat.Region(rect)
		defer croppedMat.Close()

		croppedImg, err := croppedMat.ToImage()
		if err != nil {
			return nil, err
		}

		croppedMat, err = gocv.ImageToMatRGB(croppedImg)
		if err != nil {
			return nil, err
		}

		croppedImageBytes, err = MatToString(croppedMat)
		if err != nil {
			return nil, err
		}
	}

	resultImage, err := MatToString(mat)
	if err != nil {
		return nil, err
	}

	return &DetectFace{
		Image:     *resultImage,
		FaceTotal: len(rects),
		FaceImage: croppedImageBytes,
	}, nil
}
