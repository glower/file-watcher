package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cnt := 0
	for {
		cnt++
		WriteRandomFile("C:\\Users\\Igor\\Downloads\\test")
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("total files written: %d\n", cnt)
	}
}

func WriteRandomFile(root string) {
	name := String(15)
	x := String(1)
	path := filepath.Join(root, x, name)
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		log.Printf("MkdirAll %q: %s", path, err)
	}

	d1 := []byte("hello\ngo\n")
	err = ioutil.WriteFile(path, d1, 0644)
	if err != nil {
		log.Printf("WriteFile %q: %s", path, err)
	}
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
