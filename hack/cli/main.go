package main

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/glower/file-watcher/notification"
	"github.com/glower/file-watcher/watcher"
)

func main() {
	log.Println("Starting the service ...")
	ctx := context.TODO()

	w := watcher.Setup(
		ctx,
		[]string{"C:\\Users\\Igor\\Files", "C:\\Users\\Igor\\Downloads"},
		// []string{"C:\\Users\\Igor\\Downloads"},
		[]notification.ActionType{},
		[]string{".crdownload", ".lock", ".snapshot"},
		nil)

	go func() {
		time.Sleep(30 * time.Second)
		w.StopWatching("C:\\Users\\Igor\\Files")
		time.Sleep(50 * time.Second)
		w.StopWatching("C:\\Users\\Igor\\Downloads")
	}()

	for {
		select {
		case file := <-w.EventCh:
			log.Printf("[EVENT] %#v", file)
		case err := <-w.ErrorCh:
			log.Printf("[%s] %s\n", err.Level, err.Message)
		}
	}
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
// https://golangcode.com/print-the-current-memory-usage/
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats

	if m.Alloc > 1024*1024 {
		log.Printf("Alloc = %d MiB", bToMb(m.Alloc))
	} else {
		log.Printf("Alloc = %d KiB (%d b)", bToKb(m.Alloc), m.Alloc)
	}
	log.Printf("\tTotalAlloc = %d MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %d MiB", bToMb(m.Sys))
	log.Printf("\tNumGC = %d\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
