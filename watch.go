package watchers

import (
	"context"
	"log"

	"github.com/glower/file-watchers/types"
	"github.com/glower/file-watchers/watch"
)

// Setup adds a watcher for a file changes in specified directories and returns a channel for notifications
func Setup(ctx context.Context, dirs []string, filters []types.Action) chan types.FileChangeNotification {
	log.Printf("watchers.SetupFSWatchers(): for %v\n", dirs)

	fileChangeNotificationChan := make(chan types.FileChangeNotification)

	for _, dir := range dirs {
		w := watch.SetupDirectoryWatcher(fileChangeNotificationChan, filters)
		go w.StartWatching(dir)
	}

	return fileChangeNotificationChan
}
