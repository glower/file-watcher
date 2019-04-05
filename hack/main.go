package main

import (
	"context"
	"log"

	"github.com/glower/file-watcher/notification"
	"github.com/glower/file-watcher/watcher"
)

func main() {
	log.Println("Starting the service ...")
	ctx := context.TODO()

	eventCh, errorCh := watcher.Setup(
		ctx,
		[]string{"/home/igor/Download"},
		[]notification.ActionType{},
		[]string{".crdownload", ".lock", ".snapshot"},
		nil)

	for {
		select {
		case file := <-eventCh:
			log.Printf("[EVENT] %#v", file)
		case err := <-errorCh:
			log.Printf("[%s] %s\n", err.Level, err.Message)
			if err.Level == "CRITICAL" {
				log.Printf("%s\n", err.Stack)
			}
		}
	}

}
