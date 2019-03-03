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
		[]string{"/srv/torrent/downloads", "/home/igor/Downloads"},
		[]notification.ActionType{},
		[]string{".crdownload", ".lock", ".snapshot"},
		nil)

	for {
		select {
		case file := <-eventCh:
			log.Printf("%#v", file)
		case err := <-errorCh:
			log.Printf("[ERROR] %#v", err)
		}
	}

}
