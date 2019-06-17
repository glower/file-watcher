package watcher

import (
	"context"

	"github.com/glower/file-watcher/notification"
)

type Watch struct {
	EventCh chan notification.Event
	ErrorCh chan notification.Error

	watcher *DirectoryWatcher
}

// Setup adds a watcher for a file changes in specified directories and returns a channel for notifications
func Setup(ctx context.Context, dirsToWatch []string, actionFilters []notification.ActionType, fileFilters []string, options *Options) *Watch {
	eventCh := make(chan notification.Event)
	errorCh := make(chan notification.Error)

	if options == nil {
		options = &Options{IgnoreDirectoies: true}
	}

	watcher := Create(eventCh, errorCh, actionFilters, fileFilters, options)
	w := &Watch{
		ErrorCh: errorCh,
		EventCh: eventCh,
		watcher: watcher,
	}

	for _, dir := range dirsToWatch {
		go watcher.StartWatching(dir)
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
