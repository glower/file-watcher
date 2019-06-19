package watcher

import (
	"context"

	"github.com/glower/file-watcher/notification"
)

// Watch ...
type Watch struct {
	ctx     context.Context
	EventCh chan notification.Event
	ErrorCh chan notification.Error

	watcher *DirectoryWatcher
}

// Setup adds a watcher for a file changes in specified directories and returns a channel for notifications
func Setup(ctx context.Context, options *Options) *Watch {
	eventCh := make(chan notification.Event)
	errorCh := make(chan notification.Error)
	if options == nil {
		options = &Options{IgnoreDirectoies: true}
	}

	watcher := Create(ctx, eventCh, errorCh, options)
	w := &Watch{
		ctx:     ctx,
		ErrorCh: errorCh,
		EventCh: eventCh,
		watcher: watcher,
	}

	return w
}

// StopWatching ...
func (w *Watch) StopWatching(dir string) {
	w.watcher.StopWatching(dir)
}

// StartWatching ...
func (w *Watch) StartWatching(dir string) {
	go w.watcher.StartWatching(dir)
}

// CreateFileAddedNotification creates fileAdded event for given path and file
func (w *Watch) CreateFileAddedNotification(watchDirectoryPath, relativeFilePath string) {
	w.watcher.CreateFileAddedNotification(watchDirectoryPath, relativeFilePath)
}
