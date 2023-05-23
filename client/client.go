package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"

	"github.com/deep0ne/grpc-youtube-thumbnail/protos"
)

const address string = "localhost:8081"

// Uber Go Style error naming

type URLFlags []string

// interface implementation for flag.Var
func (u *URLFlags) String() string {
	return "Just to implement interface"
}

func (u *URLFlags) Set(value string) error {
	values := strings.Split(value, ",")
	for _, val := range values {
		*u = append(*u, val)
	}
	return nil
}

// function to form thumbnails from "--urls" flag
func FormThumbnails(urls []string) []*protos.Video {
	videos := make([]*protos.Video, len(urls))
	for i := 0; i < len(urls); i++ {
		videos[i] = &protos.Video{URL: urls[i]}
	}
	return videos
}

// function to form thumbnails from links in a file
func FormThumbnailsFromFile(fileToParse string) ([]*protos.Video, error) {
	videos := make([]*protos.Video, 0)
	file, err := os.Open("./../" + fileToParse)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		link := scanner.Text()
		videos = append(videos, &protos.Video{URL: link})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return videos, nil
}

// Uber Go Style: Exit Once Rule
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		URLs        URLFlags
		Async       bool
		FileToParse string
		videos      []*protos.Video
	)

	flag.Var(&URLs, "urls", "Usage: pass youtube urls separated by comma")
	flag.BoolVar(&Async, "async", false, "Usage: flag for downloading large amount of files asynchronously")
	flag.StringVar(&FileToParse, "f", "", "Usage: pass -f flag if you want to get thumbnails from links in a file")
	flag.Parse()

	if len(URLs) == 0 && len(FileToParse) == 0 {
		return errors.New("you must pass at least one URL to parse or a file with links. See usage of \"urls\" & \"f\" flags")
	}

	if len(URLs) > 0 && len(FileToParse) > 0 {
		return errors.New("you must pass either urls or file with links. See usage of \"urls\" & \"f\" flags")
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return errors.New("error while grpc dialing to localhost...")
	}

	defer conn.Close()
	c := protos.NewYouTubeThumbnailClient(conn)

	if len(FileToParse) == 0 {
		videos = FormThumbnails(URLs)
	} else {
		videos, err = FormThumbnailsFromFile(FileToParse)
		if err != nil {
			return err
		}
	}

	if Async {
		var (
			wg      sync.WaitGroup
			errChan = make(chan error, len(videos))
		)
		for _, video := range videos {
			wg.Add(1)
			go func(video *protos.Video) {
				defer wg.Done()
				thumbnail, err := c.GetThumbnail(context.Background(), video)
				if err != nil {
					errChan <- err
				}

				status, err := c.SaveThumbnail(context.Background(), thumbnail)
				if err != nil {
					errChan <- err
				} else {
					log.Println("Thumbnail Status--->", status.Status)
				}
			}(video)
		}

		go func() {
			wg.Wait()
			close(errChan)
		}()

		for err := range errChan {
			if err != nil {
				return err
			}
		}

	} else {
		for _, video := range videos {
			thumbnail, err := c.GetThumbnail(context.Background(), video)
			if err != nil {
				return err
			}

			status, err := c.SaveThumbnail(context.Background(), thumbnail)
			if err != nil {
				return err
			}
			log.Println(status.Status)
		}
	}
	return nil
}
