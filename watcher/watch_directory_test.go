// +build integration

package watcher

import (
	"testing"

	"github.com/glower/file-watcher/types"
)

func TestSetupDirectoryWatcher(t *testing.T) {
	type args struct {
		callbackChan chan types.FileChangeNotification
		filters      []types.Action
	}

	fileChangeNotificationChan := make(chan types.FileChangeNotification)

	tests := []struct {
		name string
		args args
		dir  string
		want *types.FileChangeNotification
	}{
		{
			name: "test 1: file change notification",
			args: args{
				callbackChan: fileChangeNotificationChan,
				filters:      []types.Action{},
			},
			dir: "/test1",
			want: &types.FileChangeNotification{
				Action:             1,
				BackupToStorages:   []string(nil),
				MimeType:           "image/jpeg",
				Machine:            "tokyo",
				FileName:           "file1.txt",
				AbsolutePath:       "\\foo\\bar\\test\\file1.txt",
				RelativePath:       "test/file1.txt",
				DirectoryPath:      "/test1",
				WatchDirectoryName: "foo",
				Size:               12345,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := SetupDirectoryWatcher(tt.args.callbackChan, tt.args.filters)
			w.StartWatching(tt.dir)
			action := <-tt.args.callbackChan

			if action.Action != tt.want.Action {
				t.Errorf("action.Action = %v, want %v", action.Action, tt.want.Action)
			}

			if action.MimeType != tt.want.MimeType {
				t.Errorf("action.MimeType = %v, want %v", action.MimeType, tt.want.MimeType)
			}
		})
	}
}