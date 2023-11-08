package socketRouter

// import (
// 	"context"
// 	"encoding/json"
// 	"face-recognition-golang/constants"
// 	"face-recognition-golang/db"
// 	"face-recognition-golang/db/collection"
// 	"face-recognition-golang/libs/faceCustom"
// 	"face-recognition-golang/libs/recognizer"
// 	commonUtils "face-recognition-golang/utils/common"
// 	"fmt"
// 	"log"

// 	"sync"
// 	"time"

// 	"github.com/gofiber/contrib/websocket"
// 	"github.com/samber/lo"
// )

// type ImageQueue struct {
// 	Images  []string
// 	Mutex   sync.Mutex
// 	MaxSize int // Số lượng tối đa ảnh được lưu trữ
// }

// func NewImageQueue(maxSize int) *ImageQueue {
// 	return &ImageQueue{
// 		MaxSize: maxSize,
// 	}
// }

// func (q *ImageQueue) AddImage(imagePath string) {
// 	q.Mutex.Lock()
// 	defer q.Mutex.Unlock()
// 	q.Images = append(q.Images, imagePath)
// 	if len(q.Images) > q.MaxSize {
// 		// Nếu số lượng ảnh vượt quá giới hạn, loại bỏ ảnh cũ nhất
// 		q.Images = q.Images[1:]
// 	}
// }

// func (q *ImageQueue) GetNextImage() (string, bool) {
// 	q.Mutex.Lock()
// 	defer q.Mutex.Unlock()
// 	if len(q.Images) == 0 {
// 		return "", false
// 	}
// 	image := q.Images[0]
// 	q.Images = q.Images[1:]
// 	return image, true
// }
// func processImages(ctx context.Context, c *websocket.Conn, imageQueue *ImageQueue, username *string) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			// User disconnect, stop processing images
// 			fmt.Println("User disconnected, stopping image processing")
// 			return
// 		default:
// 			base64Image, ok := imageQueue.GetNextImage()
// 			if !ok {
// 				// No images in the queue, sleep and check again
// 				time.Sleep(time.Second)
// 				continue
// 			}

// 			bytes, err := commonUtils.Base64ToBytes(base64Image)
// 			image, err := DetectAndDrawRectangles(bytes)
// 			if err != nil {
// 				fmt.Println(err)

// 			}

// 			if image.FaceTotal != 1 {
// 				fmt.Println("Face not found", image.FaceTotal)
// 				continue
// 			}

// 			isAuth, err := faceCustom.RecognizeFace(constants.FacesDescriptor[*username].Data, image.Image)

// 			if err != nil {
// 				fmt.Println(err)
// 				continue
// 			}

// 			if isAuth != nil && *isAuth {
// 				fmt.Println("isAuth === ", *isAuth)
// 				authMessage := map[string]interface{}{
// 					"event":  "auth",
// 					"isAuth": true,
// 				}

// 				authMessageJSON, err := json.Marshal(authMessage)
// 				if err != nil {
// 					log.Println("marshal:", err)
// 					continue
// 				}

// 				err = c.WriteMessage(websocket.TextMessage, authMessageJSON)
// 				if err != nil {
// 					log.Println("write:", err)
// 					continue
// 				}
// 			}

// 		}
// 	}
// }

// var WsRecognize = websocket.New(func(c *websocket.Conn) {
// 	log.Println("Kết nối WebSocket đã được thiết lập")
// 	var username string
// 	imageQueue := NewImageQueue(10)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	go processImages(ctx, c, imageQueue, &username)

// 	for {
// 		messageType, p, err := c.ReadMessage()
// 		if err != nil {
// 			log.Println("read:", err)
// 			break
// 		}

// 		if messageType == websocket.TextMessage {
// 			username = string(p)
// 			if constants.FacesDescriptor[username] == nil {
// 				constants.FacesDescriptor[username] = &constants.FaceDescriptor{}
// 				constants.FacesDescriptor[username].Pending = true
// 				go func() {
// 					cursor, err := db.Client.DatasetCollectionName.Find(context.Background(), map[string]interface{}{
// 						"Username": username,
// 					})

// 					if err != nil {
// 						fmt.Println(err)
// 						return
// 					}

// 					var results []collection.Dataset
// 					if err = cursor.All(context.TODO(), &results); err != nil {
// 						panic(err)
// 					}

// 					constants.FacesDescriptor[username].Data = lo.Map(results, func(item collection.Dataset, _ int) recognizer.Data {
// 						return recognizer.Data{
// 							Descriptor: item.Descriptor,
// 							Id:         item.UserID,
// 						}
// 					})

// 					constants.FacesDescriptor[username].Pending = false
// 				}()
// 			}
// 		} else if messageType == websocket.BinaryMessage && username != "" {
// 			if constants.FacesDescriptor[username] == nil || constants.FacesDescriptor[username].Pending {
// 				fmt.Println("Face descriptor not found")
// 			} else {
// 				fmt.Println("Face descriptor found")
// 				imageBase64, err := commonUtils.BytesToBase64Image(p)
// 				if err != nil {
// 					fmt.Println(err)
// 					continue
// 				}
// 				imageQueue.AddImage(imageBase64)
// 			}

// 		}
// 	}

// 	log.Println("Kết nối WebSocket đã đóng")
// })
