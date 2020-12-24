package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

var changedFlag uint32

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go fileEventWatcher(watcher)
	go testRunner()

	err = discoverFiles(".", watcher)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	<-done
}

// searches the given path (not following symlinks) and watches all directories and go related files
func discoverFiles(path string, watcher *fsnotify.Watcher) error {
	return filepath.Walk(
		path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				// ignore hidden directories like .git
				base := filepath.Base(path)
				if base != "." && strings.HasPrefix(base, ".") {
					return filepath.SkipDir
				}
				err = watcher.Add(path)
				if err != nil {
					log.Fatal(err)
				}
			}

			return nil
		},
	)
}

func fileEventWatcher(watcher *fsnotify.Watcher) {
	validFilePattern := regexp.MustCompile(`(go\.mod)|(.*\.go)`)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			info, err := os.Stat(event.Name)
			if err != nil {
				// most likely no access. lets ignore it
				continue
			}

			if info.IsDir() {
				if event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Rename == fsnotify.Rename {
					err := discoverFiles(event.Name, watcher)
					if err != nil {
						log.Fatal(err)
					}
				}
				continue
			}

			if validFilePattern.MatchString(filepath.Base(event.Name)) {
				atomic.StoreUint32(&changedFlag, 1)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Fatal("error:", err)
		}
	}
}

func runTest() {
	args := append([]string{"test"}, os.Args[1:]...)
	cmd := exec.Command("go", args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	log.Println("running ", cmd.String())
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err)
	}
}

func testRunner() {
	buildTriggerd := false
	for {
		changed := atomic.LoadUint32(&changedFlag) == 1
		atomic.StoreUint32(&changedFlag, 0)

		if !buildTriggerd && !changed {
			runTest()
			buildTriggerd = true
			continue
		}
		if buildTriggerd && changed {
			buildTriggerd = false
			time.Sleep(10 * time.Millisecond)
			continue
		}
		time.Sleep(500 * time.Millisecond)
	}
}
