// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	"os/exec"
	"os"
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {

					fi, err := os.Stat(filePath);
					if err != nil {
						log.Println("error:", err)
						return
					}
					
					log.Println("modified file:", event.Name, "file size: ", fi.Size())

					if fi.Size() != 9047626 && fi.Size() != 0 {
						cmd := exec.Command(bashAlias, bashParameterAlias, chwp)
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
					
						err = cmd.Run()
						if err != nil {
							log.Println("error: ", err)
							return
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}