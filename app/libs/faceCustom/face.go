package faceCustom

import (
	"face-recognition-golang/constants/path"
	"face-recognition-golang/libs/recognizer"
	"fmt"
)

var DataDir = path.GetBasepath() + "/public/models"

type Face struct {
	Data     []string
	UserName string
}

type Faces struct {
	Users map[string]*Face
}

const MAX_SIZE = 100

func (f *Faces) AddFace(username string, image string) bool {
	if f.Users[username] == nil {
		f.Users[username] = &Face{
			Data:     []string{},
			UserName: username,
		}
	}
	if len(f.Users[username].Data) == MAX_SIZE {
		return false
	}

	f.Users[username].Data = append(f.Users[username].Data, image)
	return len(f.Users[username].Data) == MAX_SIZE
}

func AddFile(rec *recognizer.Recognizer, Path, Id string) {
	err := rec.AddImageToDataset(Path, Id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func FaceClassify(images []string, username string) ([]recognizer.Data, error) {
	rec := recognizer.Recognizer{}
	err := rec.Init(DataDir)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rec.Tolerance = 0.3
	rec.UseGray = true
	rec.UseCNN = false
	defer rec.Close()

	for _, image := range images {
		AddFile(&rec, image, username)
	}

	rec.SetSamples()
	return rec.Dataset, nil
}

func RecognizeFace(recognizerData []recognizer.Data, image string) (*bool, error) {
	rec := recognizer.Recognizer{}
	err := rec.Init(DataDir)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rec.Tolerance = 0.2
	rec.UseGray = true
	rec.UseCNN = false
	defer rec.Close()

	rec.Dataset = recognizerData

	rec.SetSamples()

	faces, err := rec.ClassifyMultiples(image)

	if err != nil {
		if err.Error() == "Not a single face on the image" {
			result := false
			return &result, nil
		}
		fmt.Println("err::", err)
		return nil, err
	}

	var result bool
	result = len(faces) == 1 && faces[0].Id != ""
	return &result, nil
}
