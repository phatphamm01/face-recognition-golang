package commonUtils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func ReadImageFromFile(fileName string) (*string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := fileInfo.Size()
	buffer := make([]byte, size)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	base64Image := base64.StdEncoding.EncodeToString(buffer)

	return &base64Image, nil
}

func WriteImageToFile(fileName string, content string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decodedImage, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}

	if _, err := file.Write(decodedImage); err != nil {
		return err
	}

	return nil

}

func PrintJSON(data interface{}) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(dataJSON))
}

func BytesToBase64Image(imageBytes []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	if err := jpeg.Encode(encoder, img, nil); err != nil {
		return "", err
	}

	encoder.Close()
	base64Image := buffer.String()

	return base64Image, nil
}

func ImageToBase64(image image.Image) (string, error) {
	var buffer bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	if err := jpeg.Encode(encoder, image, nil); err != nil {
		return "", err
	}

	encoder.Close()
	base64Image := buffer.String()

	return base64Image, nil
}

func Base64ToImage(base64Image string) (image.Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(base64Image)))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func SaveImage(base64Image string, fileName string) error {
	image, err := Base64ToImage(base64Image)
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	if err := jpeg.Encode(file, image, nil); err != nil {
		return err
	}

	return nil
}

func Base64ToBytes(base64Image string) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(base64Image)))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func BytesToImage(imageBytes []byte) (*image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	return &img, nil
}
