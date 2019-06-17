// +build windows,!integration

package watcher

// #include "watch_windows.h"
import "C"
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/glower/file-watcher/notification"
	fileinfo "github.com/glower/file-watcher/util"
)

func init() {
	C.Setup()
}

// StartWatching starts a CGO function for getting the notifications
func (w *DirectoryWatcher) StartWatching(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fileError("CRITICAL", err)
		return
	}

	log.Printf("windows.StartWatching(): for [%s]\n", path)
	cpath := C.CString(path)
	defer func() {
		C.free(unsafe.Pointer(cpath))
	}()

	// TODO: refactor me
	go func() {
		for {
			select {
			case p := <-w.StopWatchCh:
				fmt.Printf("Income channel message to stop [%s] directory watcher\n", p)
				if p == path {
					fmt.Printf(">>> Stoping [%s] directory watcher\n", p)
					C.StopWatching(cpath)
				} else {
					w.StopWatchCh <- p
				}
			}
		}
	}()

	C.WatchDirectory(cpath)
	fmt.Printf(">>> Stoping DONE for [%s]\n", path)
}

//export goCallbackFileChange
func goCallbackFileChange(cpath, cfile *C.char, caction C.int) {
	path := strings.TrimSpace(C.GoString(cpath))
	file := strings.TrimSpace(C.GoString(cfile))
	action := notification.ActionType(int(caction))

	fmt.Printf(">>> goCallbackFileChange(): path=%s, file=%s, action=%s\n", path, file, ActionToString(action))

	absoluteFilePath := filepath.Join(path, file)
	fi, err := fileinfo.GetFileInformation(absoluteFilePath)

	if err != nil {
		fileError("WARN", err)
		return
	}

	fileChangeNotifier(path, file, fi, action)
}
