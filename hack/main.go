package main

import (
	"context"
	"log"

	watchers "github.com/glower/file-watchers"
	"github.com/glower/file-watchers/types"
)

func main() {
	log.Println("Starting the service ...")
	ctx := context.TODO()

	fileChangeNotificationChan := watchers.Setup(ctx, []string{"/srv/torrent/downloads", "/home/igor/Downloads"}, []types.Action{})

	for {
		select {
		case file := <-fileChangeNotificationChan:
			log.Printf("%#v", file)
		}
	}

}
