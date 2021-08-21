package task

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/panjf2000/ants"
	log "github.com/sirupsen/logrus"
	"medilane-worker/database"
	fcm "medilane-worker/firebase"
	"medilane-worker/models"
	"medilane-worker/queue"
	"os"
	"time"
)

const (
	serverKey = "YOUR-KEY"
)

type NotificationWorker struct {
}

func NewNotificationWorker() *NotificationWorker {
	return &NotificationWorker{}
}

func (n *NotificationWorker) HandleNotification(i interface{}) {
	str := fmt.Sprintf("%v", i)
	var notification models.NotificationQueue
	_ = jsoniter.Unmarshal([]byte(str), &notification)
	fmt.Printf("run with %v\n", i)

	// save to database
	err := n.SaveNotification(notification)
	if err != nil {
		log.Errorf("Error when save notification to database: %v", err.Error())
	}
	// send notify to firebase
	n.SendToToken(&notification)
}

func (n *NotificationWorker) GetFcmTokens(userIds []uint) []string {
	var tokens []models.FcmToken
	err := database.GetInstance().GetFcmToken(userIds, &tokens)
	tokenList := make([]string, 0)
	if err != nil {
		return tokenList
	}
	for _, item := range tokens {
		tokenList = append(tokenList, item.Token)
	}
	return tokenList
}

func (n *NotificationWorker) SendToToken(notification *models.NotificationQueue) {
	// [START send_to_token_golang]
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := fcm.GetInstance().Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationTokens := n.GetFcmTokens(notification.UserId)

	// See documentation on defining a message payload.
	notiStr, _ := jsoniter.Marshal(notification)
	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"body":  notification.Message,
			"title": notification.Title,
			"data":  string(notiStr),
		},
		Tokens: registrationTokens,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	br, err := client.SendMulticast(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_to_token_golang]
}

func (n *NotificationWorker) SaveNotification(data models.NotificationQueue) error {
	var notifications []models.Notification
	for _, item := range data.UserId {
		notiItem := models.Notification{
			EntityId: data.EntityId,
			Action:   data.Action,
			Entity:   data.Entity,
			Status:   data.Status,
			Message:  data.Message,
			UserId:   item,
			Title:    data.Title,
		}
		notifications = append(notifications, notiItem)
	}
	return database.GetInstance().InsertManyNotification(notifications)
}

func (n *NotificationWorker) GetNotificationFromQueue(queueName string) string {
	ctx := context.Background()
	rs, err := queue.GetInstance().RPop(ctx, queueName)
	if err != nil {
		if err.Error() == "redis: nil" {
		} else {
			log.Errorf("Error when pop data from queue: %v", err.Error())
		}
		return ""
	}
	return rs
}

func (n *NotificationWorker) Run(signalChan chan os.Signal) {
	// Use the pool with a method,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		n.HandleNotification(i)
	})

	for {
		select {
		case <-signalChan:
			p.Release()
			return
		default:
			data := n.GetNotificationFromQueue("notification")
			if data != "" {
				_ = p.Invoke(data)
			}
			fmt.Printf("running goroutines: %d\n", p.Running())
			time.Sleep(1 * time.Second)
		}
	}
}
