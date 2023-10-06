package watcher

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/config"
)

func StartLogWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}

	logFile := common.GetLogfileName(config.GetEnv().LogDir)
	err = watcher.Add(logFile)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				// contents, err := os.ReadFile(logFile)
				// if err != nil {
				// 	fmt.Println("[CRITICAL]", err)
				// }
				fmt.Println("modified file:", ev.Op.String())
			}
		case err := <-watcher.Errors:
			fmt.Println("[CRITICAL]", err)
		}
	}
}
