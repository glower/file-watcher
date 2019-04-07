# file-watcher  

Get notifications of file change in specified directory. Works on on linux and windows using cgo bindings.

# Download

`go get github.com/glower/file-watcher`

# Usage

``` go
package main

import (
	"context"
	"log"

	"github.com/glower/file-watcher/notification"
	"github.com/glower/file-watcher/watcher"
)

func main() {
	ctx := context.TODO()

	eventCh, errorCh := watcher.Setup(
		ctx,
		[]string{"/home/igor/Downloads", "C:\\Users\\Igor\\Downloads"},
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
```