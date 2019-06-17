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
	ch := RegisterCallback(path)

	log.Printf("windows.StartWatching(): for [%s]\n", path)
	cpath := C.CString(path)
	defer func() {
		UnregisterCallback(path)
		C.free(unsafe.Pointer(cpath))
	}()

	go func() {
		for {
			select {
			case p := <-ch:
				fmt.Printf("Income channel message to stop directory watcher\n")
				if p.Stop {
					fmt.Printf(">>> Stoping directory watcher\n")
					C.StopWatching(cpath)
				}
			}
		}
	}()

	C.WatchDirectory(cpath)
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
