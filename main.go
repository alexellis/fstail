package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

func main() {

	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Printf(`fstail - Copyright Alex Ellis 2023

Usage:

  # Work in current directory
  fstail

  # Work in /var/log/actuated
  fstail /var/log/actuated

  # Work in /var/log/actuated and match strings with "server error"
  fstail /var/log/actuated "server error"

  # Disable prefix printing
  FS_PREFIX=0 fstail

  # Print a prefix with the calculated Pod/container name
  FS_PREFIX=k8s fstail /var/log/containers

`)
			return
		}
	}

	var (
		wd            string
		match         string
		k8sPrefix     bool
		disablePrefix bool
	)
	log.Println(len(os.Args))

	if v, ok := os.LookupEnv("FS_PREFIX"); ok {
		if v == "0" {
			disablePrefix = true
		} else if v == "k8s" {
			k8sPrefix = true
		}
	}

	if len(os.Args) == 2 {
		wd = os.Args[1]
	} else if len(os.Args) == 3 {
		wd = os.Args[1]
		match = os.Args[2]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		wd = cwd
	}

	prefixStyle := "none"
	if k8sPrefix {
		prefixStyle = "k8s"
	} else if !disablePrefix {
		prefixStyle = "filename"
	}

	matchStyle := "*"
	if len(match) > 0 {
		matchStyle = match
	}

	fmt.Printf("Watching: %s match: %s, prefix: %s\n", wd, matchStyle, prefixStyle)

	printers := map[string]*Streamer{}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if len(match) > 0 {
		files, err := os.ReadDir(wd)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {

			if !strings.Contains(file.Name(), match) {
				continue
			}

			log.Printf("Attaching to: %s", file.Name())

			if f, err := os.Open(path.Join(wd, file.Name())); err == nil {
				s := NewStreamer(f, k8sPrefix, disablePrefix)
				go s.Stream()
				printers[path.Join(wd, file.Name())] = s
			} else {
				log.Println(err)
			}

		}
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Printf("Event: %s %s", event.Name, event.Op.String())

				if len(match) > 0 && !strings.Contains(event.Name, match) {
					log.Printf("Skipping: %s", event.Name)
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					if _, ok := printers[event.Name]; !ok {
						if f, err := os.Open(event.Name); err == nil {
							s := NewStreamer(f, k8sPrefix, disablePrefix)
							go s.Stream()
							printers[event.Name] = s
						} else {
							log.Println(err)
						}
					}
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					if _, ok := printers[event.Name]; !ok {
						if f, err := os.Open(event.Name); err == nil {
							s := NewStreamer(f, k8sPrefix, disablePrefix)
							go s.Stream()
							printers[event.Name] = s
						} else {
							log.Println(err)
						}
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if _, ok := printers[event.Name]; ok {
						printers[event.Name].Close()
						delete(printers, event.Name)
					}
				}

			case err := <-watcher.Errors:
				if err != nil {
					log.Fatalln("Error:", err)
				}
			}
		}
	}()

	log.Printf("Adding watch for: %s", wd)
	if err = watcher.Add(wd); err != nil {
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

	k8sPrefix     bool
	disablePrefix bool
}

func NewStreamer(f *os.File, k8sPrefix bool, disablePrefix bool) *Streamer {
	return &Streamer{f: f, k8sPrefix: k8sPrefix, disablePrefix: disablePrefix}
}

func (s *Streamer) Stream() {
	base := path.Base(s.f.Name())

	var prefix string

	if !s.k8sPrefix && !s.disablePrefix {
		prefix = fmt.Sprintf("%s| ", base)
	} else if s.k8sPrefix {
		podSt, _, ok := strings.Cut(base, "_")
		if ok {
			prefix = fmt.Sprintf("%s| ", podSt)
		}
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
