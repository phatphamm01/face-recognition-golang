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

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var WsRecognizeRecord = websocket.New(func(c *websocket.Conn) {
	log.Println("Kết nối WebSocket đã được thiết lập")
	var userId string

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
			base64Image, err := DetectAndDrawRectangles(p)
			if err != nil {
				fmt.Println(err)
				return
			}

			if base64Image.FaceImage != nil {
				isDone := constants.FacesInstance.AddFace(userId, *base64Image.FaceImage)

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

					data, err := faceCustom.FaceClassify(constants.FacesInstance.Users[userId].Data, userId)
					if err != nil {
						fmt.Println(err)
						return
					}

					newData := make([]interface{}, len(data))
					for i, item := range data {
						newData[i] = collection.Dataset{
							Descriptor: item.Descriptor,
							UserID:     userId,
						}
					}

					_, err = db.Client.DatasetCollectionName.InsertMany(context.TODO(), newData)
					if err != nil {
						fmt.Println(err)
						return
					}
					objID, _ := primitive.ObjectIDFromHex(userId)
					filter := bson.M{"_id": objID}
					update := bson.M{"$set": bson.M{"isVerifiedFace": true}}
					res, err := db.Client.DB.Collection("users").UpdateOne(context.Background(), filter, update)
					if err != nil {
						fmt.Println(err)
						return
					}

					fmt.Println(res)
					constants.FacesInstance.Users[userId] = nil
				}
			}

			imageAfterRecognize := map[string]interface{}{
				"event":     "IMAGE_AFTER_RECOGNIZE",
				"image":     base64Image.Image,
				"faceTotal": base64Image.FaceTotal,
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
	log.Println("Kết nối WebSocket đã đóng")
})
