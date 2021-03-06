// +build linux,!integration

package watcher

// #include "watch_linux.h"
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

// #define IN_ACCESS		0x00000001	/* File was accessed */
// #define IN_MODIFY		0x00000002	/* File was modified */
// #define IN_ATTRIB		0x00000004	/* Metadata changed */
// #define IN_CLOSE_WRITE	0x00000008	/* Writtable file was closed */
// #define IN_CLOSE_NOWRITE	0x00000010	/* Unwrittable file closed */
// #define IN_OPEN			0x00000020	/* File was opened */
// #define IN_MOVED_FROM	0x00000040	/* File was moved from X */
// #define IN_MOVED_TO		0x00000080	/* File was moved to Y */
// #define IN_CREATE		0x00000100	/* Subfile was created */
// #define IN_DELETE		0x00000200	/* Subfile was deleted */
// #define IN_DELETE_SELF	0x00000400	/* Self was deleted */
func convertMaskToAction(mask int) notification.ActionType {
	switch mask {
	case 2, 8: // File was modified
		return notification.ActionType(notification.FileModified)
	case 256: // Subfile was created
		return notification.ActionType(notification.FileAdded)
	case 512: // Subfile was deleted
		return notification.ActionType(notification.FileRemoved)
	case 64: // File was moved from X
		return notification.ActionType(notification.FileRenamedOldName)
	case 128: // File was moved to Y
		return notification.ActionType(notification.FileRenamedNewName)
	default:
		return notification.ActionType(notification.Invalid)
	}
}

// StartWatching starts a CGO function for getting the notifications
func (i *DirectoryWatcher) StartWatching(root string) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fileError("CRITICAL", fmt.Errorf("cannot start watching [%s]: no such directory", root))
		return
	}
	log.Printf("linux.StartWatching(): for [%s]\n", root)
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			go watchDir(root, path)
		}
		return nil
	})
	if err != nil {
		fileError("ERROR", err)
	}
}

func watchDir(rootDirToWatch string, subDir string) {
	croot := C.CString(rootDirToWatch)
	cdir := C.CString(subDir)
	defer func() {
		C.free(unsafe.Pointer(croot))
		C.free(unsafe.Pointer(cdir))
	}()
	C.WatchDirectory(croot, cdir)
}

//export goCallbackFileChange
func goCallbackFileChange(croot, cpath, cfile *C.char, caction C.int) {
	root := strings.TrimSpace(C.GoString(croot))
	path := strings.TrimSpace(C.GoString(cpath))
	file := strings.TrimSpace(C.GoString(cfile))
	action := convertMaskToAction(int(caction))

	absoluteFilePath := filepath.Join(path, file)
	relativeFilePath, err := filepath.Rel(root, absoluteFilePath)
	if err != nil {
		fileError("ERROR", err)
		return
	}

	fi, err := fileinfo.GetFileInformation(absoluteFilePath)
	if err != nil {
		fileError("WARN", err)
		return
	}

	fileChangeNotifier(root, relativeFilePath, fi, action, nil)
}
