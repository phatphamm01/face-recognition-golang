package socketRouter

import (
	"context"
	"encoding/json"
	"face-recognition-golang/constants"
	"face-recognition-golang/db"
	"face-recognition-golang/db/collection"
	"face-recognition-golang/libs/faceCustom"
	commonUtils "face-recognition-golang/utils/common"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageQueue struct {
	Images  []string
	Mutex   sync.Mutex
	MaxSize int // Số lượng tối đa ảnh được lưu trữ
}

func NewImageQueue(maxSize int) *ImageQueue {
	return &ImageQueue{
		MaxSize: maxSize,
	}
}

func (q *ImageQueue) AddImage(imagePath string) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	q.Images = append(q.Images, imagePath)
	if len(q.Images) > q.MaxSize {
		// Nếu số lượng ảnh vượt quá giới hạn, loại bỏ ảnh cũ nhất
		q.Images = q.Images[1:]
	}
}

func (q *ImageQueue) GetNextImage() (string, bool) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if len(q.Images) == 0 {
		return "", false
	}
	lenImages := len(q.Images)
	image := q.Images[lenImages-1]
	q.Images = q.Images[:lenImages-1]
	return image, true
}
func processImages(ctx context.Context, c *websocket.Conn, imageQueue *ImageQueue, userId *string) {
	for {
		select {
		case <-ctx.Done():
			// User disconnect, stop processing images
			fmt.Println("User disconnected, stopping image processing")
			return
		default:
			base64Image, ok := imageQueue.GetNextImage()
			if !ok {
				time.Sleep(time.Second)
				continue
			}

			bytes, err := commonUtils.Base64ToBytes(base64Image)
			if err != nil {
				continue
			}

			detectAndDrawRectangles, err := DetectAndDrawRectangles(bytes)
			if err != nil {
				fmt.Println(err)
				return
			}

			if detectAndDrawRectangles.FaceImage != nil {
				isDone := constants.FacesInstance.AddFace(*userId, *detectAndDrawRectangles.FaceImage)

				if isDone {
					stopMessage := map[string]string{
						"event": "DONE",
					}
					stopMessageJSON, err := json.Marshal(stopMessage)
					if err != nil {
						log.Println("marshal:", err)
						break
					}

					err = c.WriteMessage(websocket.TextMessage, stopMessageJSON)
					if err != nil {
						log.Println("write:", err)
						break
					}

					data, err := faceCustom.FaceClassify(constants.FacesInstance.Users[*userId].Data, *userId)
					if err != nil {
						fmt.Println(err)
						break
					}

					newData := make([]interface{}, len(data))
					for i, item := range data {
						newData[i] = collection.Dataset{
							Descriptor: item.Descriptor,
							UserID:     *userId,
						}
					}

					_, err = db.Client.DatasetCollectionName.InsertMany(context.TODO(), newData)
					if err != nil {
						fmt.Println(err)
						break
					}
					objID, _ := primitive.ObjectIDFromHex(*userId)
					filter := bson.M{"_id": objID}
					update := bson.M{"$set": bson.M{"isVerifiedFace": true}}
					res, err := db.Client.DB.Collection("users").UpdateOne(context.Background(), filter, update)
					if err != nil {
						fmt.Println(err)
						break
					}

					fmt.Println(res)
					constants.FacesInstance.Users[*userId] = nil
				}
			}

			imageAfterRecognize := map[string]interface{}{
				"event":     "IMAGE_AFTER_RECOGNIZE",
				"image":     detectAndDrawRectangles.Image,
				"faceTotal": detectAndDrawRectangles.FaceTotal,
			}

			imageAfterRecognizeJSON, err := json.Marshal(imageAfterRecognize)
			if err != nil {
				log.Println("marshal:", err)

			}

			err = c.WriteMessage(websocket.TextMessage, imageAfterRecognizeJSON)
			if err != nil {
				log.Println("write:", err)
			}
		}
	}
}

var WsRecognizeRecord = websocket.New(func(c *websocket.Conn) {
	log.Println("Kết nối WebSocket đã được thiết lập")
	var userId string
	imageQueue := NewImageQueue(10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go processImages(ctx, c, imageQueue, &userId)

	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if messageType == websocket.TextMessage {
			var data map[string]string
			err := json.Unmarshal(p, &data)
			if err != nil {
				fmt.Println(err)
				return
			}
			commonUtils.PrintJSON(data)
			switch data["event"] {
			case "USER":
				userId = data["data"]
			}
		} else if messageType == websocket.BinaryMessage && userId != "" {
			imageBase64, err := commonUtils.BytesToBase64Image(p)
			if err != nil {
				fmt.Println(err)
				continue
			}
			imageQueue.AddImage(imageBase64)
		}
	}
	log.Println("Kết nối WebSocket đã đóng")
})
