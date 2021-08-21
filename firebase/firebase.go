package fcm

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"medilane-worker/config"
	"sync"

	"google.golang.org/api/option"
)

type FirebaseConn interface {
	Init(cfg *config.Config)
	Messaging(ctx context.Context) (*messaging.Client, error)
}

var once sync.Once

type FireBaseApp struct {
	app *firebase.App
}

// singleton for api
var singletonFireBaseApp FirebaseConn

func GetInstance() FirebaseConn {
	once.Do(func() { // <-- atomic, does not allow repeating
		singletonFireBaseApp = &FireBaseApp{}
	})
	return singletonFireBaseApp
}

func SetInstance(obj FirebaseConn) {
	singletonFireBaseApp = obj
}

//NewClient new client for worker
func NewClient(cfg *config.Config) FirebaseConn {
	app := &FireBaseApp{}
	app.Init(cfg)
	return app
}

func (obj *FireBaseApp) Init(cfg *config.Config) {
	opt := option.WithCredentialsFile(cfg.FcmKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing app firebase: %v", err))
	}
	obj.app = app
}

func (obj *FireBaseApp) Messaging(ctx context.Context) (*messaging.Client, error) {
	return obj.app.Messaging(ctx)
}
