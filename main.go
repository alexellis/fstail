package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

func main() {

	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Printf(`fstail - Copyright Alex Ellis 2023

Usage:

  fstail

  fstail /var/log/actuated

  FS_PREFIX=0 fstail

`)
			return
		}

	}

	var wd string

	if len(os.Args) > 1 {
		wd = os.Args[1]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		wd = cwd
	}

	printers := map[string]*Streamer{}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:

				if event.Op&fsnotify.Write == fsnotify.Write {
					if _, ok := printers[event.Name]; !ok {
						log.Printf("Found: %s", path.Base(event.Name))
						if f, err := os.Open(event.Name); err == nil {
							s := NewStreamer(f)
							go s.Stream()
							printers[event.Name] = s
						} else {
							log.Println(err)
						}

					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Fatalln("Error:", err)
				}
			}
		}
	}()

	err = watcher.Add(wd)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		for _, s := range printers {
			s.Close()
		}
		done <- true
	}()

	<-done

}

type Streamer struct {
	f *os.File
}

func NewStreamer(f *os.File) *Streamer {
	return &Streamer{f: f}
}

func (s *Streamer) Stream() {
	base := path.Base(s.f.Name())

	prefix := fmt.Sprintf("%s| ", base)
	if v, ok := os.LookupEnv("FS_PREFIX"); ok && v == "0" {
		prefix = ""
	}

	reader := bufio.NewReader(s.f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			break
		}

		fmt.Printf("%s%s", prefix, string(line))
	}
}

func (s *Streamer) Close() {
	if s.f != nil {
		s.f.Close()
	}
}
