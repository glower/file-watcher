package watch

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/glower/file-watchers/file"
	"github.com/glower/file-watchers/notifications"
	"github.com/glower/file-watchers/types"
)

// ActionToString maps Action value to string
func ActionToString(action types.Action) string {
	switch action {
	case types.FileAdded:
		return "added"
	case types.FileRemoved:
		return "removed"
	case types.FileModified:
		return "modified"
	case types.FileRenamedOldName, types.FileRenamedNewName:
		return "renamed"
	default:
		return "invalid"
	}
}

// DirectoryWatcher ...
type DirectoryWatcher struct {
	ActionFilters              []types.Action
	FileChangeNotificationChan chan types.FileChangeNotification
	NotificationWaiter         notifications.NotificationWaiter
}

// DirectoryWatcherImplementer ...
type DirectoryWatcherImplementer interface {
	StartWatching(path string)
}

var watcher *DirectoryWatcher
var once sync.Once

// TODO: add options for filter dirs
// TODO: add filter for file names

// SetupDirectoryWatcher ...
func SetupDirectoryWatcher(callbackChan chan types.FileChangeNotification, filters []types.Action) *DirectoryWatcher {
	once.Do(func() {
		watcher = &DirectoryWatcher{
			ActionFilters:              filters,
			FileChangeNotificationChan: callbackChan,
			NotificationWaiter: notifications.NotificationWaiter{
				FileChangeNotificationChan: callbackChan,
				Timeout:                    time.Duration(5 * time.Second),
				MaxCount:                   10,
			},
		}
	})
	return watcher
}

func fileChangeNotifier(watchDirectoryPath, relativeFilePath string, fileInfo file.ExtendedFileInfoImplementer, action types.Action) error {

	// TODO: add some filter here
	if fileInfo.IsDir() {
		return fmt.Errorf("file change for a directory is not supported")
	}

	// TODO: add filter for file names
	if fileInfo.IsTemporaryFile() {
		return fmt.Errorf("file change for a tmp file is not supported")
	}

	for _, filter := range watcher.ActionFilters {
		if action == filter {
			log.Printf("fileChangeNotifier(): action [%s] is filtered\n", ActionToString(filter))
			return nil
		}
	}

	absoluteFilePath := filepath.Join(watchDirectoryPath, relativeFilePath)
	log.Printf("watch.fileChangeNotifier(): watch directory path [%s], relative file path [%s], action [%s]\n", watchDirectoryPath, relativeFilePath, ActionToString(action))

	wait, exists := watcher.NotificationWaiter.LookupForFileNotification(absoluteFilePath)
	if exists {
		wait <- true
		return nil
	}

	watcher.NotificationWaiter.RegisterFileNotification(absoluteFilePath)

	host, _ := os.Hostname()
	mimeType, err := fileInfo.ContentType()
	if err != nil {
		watcher.NotificationWaiter.UnregisterFileNotification(absoluteFilePath)
		return fmt.Errorf("can't get ContentType from the file [%s]: %v", absoluteFilePath, err)
	}

	data := &types.FileChangeNotification{
		MimeType:           mimeType,
		AbsolutePath:       absoluteFilePath,
		Action:             action,
		DirectoryPath:      watchDirectoryPath,
		Machine:            host,
		FileName:           fileInfo.Name(),
		RelativePath:       relativeFilePath,
		Size:               fileInfo.Size(),
		Timestamp:          fileInfo.ModTime(),
		WatchDirectoryName: filepath.Base(watchDirectoryPath),
	}

	go watcher.NotificationWaiter.Wait(data)

	return nil
}
