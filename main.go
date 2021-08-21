package main

import (
	"github.com/panjf2000/ants"
	"medilane-worker/config"
	"medilane-worker/database"
	fcm "medilane-worker/firebase"
	"medilane-worker/queue"
	"medilane-worker/task"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer ants.Release()

	cfg := config.NewConfig()

	// init redis
	queue.GetInstance().Init(cfg)
	defer queue.GetInstance().Close()

	// init mysql connection
	database.GetInstance().Init(cfg)

	// init firebase connection
	fcm.GetInstance().Init(cfg)

	// notification
	notification := task.NewNotificationWorker()
	go notification.Run(signalChan)

	go func() {
		<-signalChan
		os.Exit(1)
	}()
	select {
	case <-signalChan:
		time.Sleep(1 * time.Second)
	}
}
