package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deep0ne/grpc-youtube-thumbnail/protos"
)

const (
	thumbnailURL string = "https://img.youtube.com/vi/"
	image        string = "/hqdefault.jpg"
)

var ErrWrongURL = errors.New("wrong youtube video URL")

/*
Youtube links can be of two types:

1. https://www.youtube.com/watch?v=some_id_here
2. https://youtu.be/some_id_here

FormThumbnailURL gets the id of the video (some_id_here) and forms a thumbnail link
*/

func FormThumbnailURL(video *protos.Video) (string, error) {
	var id string

	if strings.Contains(video.URL, "watch") {
		splitted := strings.Split(video.URL, "=")
		fmt.Println(splitted)
		if len(splitted) != 2 {
			return "", ErrWrongURL
		}
		id = splitted[len(splitted)-1]
	} else {
		idx := strings.LastIndex(video.URL, "/")
		if idx == -1 {
			return "", ErrWrongURL
		}
		id = video.URL[idx+1:]
	}

	thumbnail := thumbnailURL + id + image
	return thumbnail, nil
}
