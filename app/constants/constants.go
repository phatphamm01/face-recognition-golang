package constants

import (
	"face-recognition-golang/libs/faceCustom"
	"face-recognition-golang/libs/recognizer"
)

var FacesInstance = faceCustom.Faces{
	Users: map[string]*faceCustom.Face{},
}

type FaceDescriptor struct {
	Pending bool
	Data    []recognizer.Data
}

var FacesDescriptor = map[string]*FaceDescriptor{}
