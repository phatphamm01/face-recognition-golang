package constants

import (
	"face-recognition-golang/libs/faceCustom"
	"face-recognition-golang/libs/recognizer"
)

const MAX_SIZE = 100

var FacesInstance = faceCustom.Faces{
	Users: map[string]*faceCustom.Face{},
}

type FaceDescriptor struct {
	Pending bool
	Data    []recognizer.Data
}

var FacesDescriptor = map[string]*FaceDescriptor{}
