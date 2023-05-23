package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/deep0ne/grpc-youtube-thumbnail/protos"
	"github.com/deep0ne/grpc-youtube-thumbnail/utils"
)

// Uber Go Style - prefix unexported globals
const (
	_defaultPort    = ":8081"
	_statusSaved    = "Image was saved to image/ folder successfully"
	_statusNotSaved = "Image was not saved"
)

type YouTubeServer struct {
	protos.UnimplementedYouTubeThumbnailServer
	Logger *logrus.Logger
	Redis  *redis.Client
}

func NewYouTubeServer() *YouTubeServer {
	return &YouTubeServer{
		Logger: utils.NewLogger(),
		Redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func (s *YouTubeServer) GetThumbnail(ctx context.Context, video *protos.Video) (*protos.Thumbnail, error) {
	if !strings.Contains(video.URL, "youtu") {
		s.Logger.Errorf("wrong youtube video URL was passed: %v", video.URL)
		return nil, utils.ErrWrongURL
	}

	thumbnail, err := s.Redis.Get(video.URL).Result()
	if err == nil {
		s.Logger.Logln(logrus.InfoLevel, "Thumbnail was found in Redis Cache. Returning...")
		return &protos.Thumbnail{URL: thumbnail}, nil
	} else if err != redis.Nil {
		s.Logger.Errorf("Something went wrong with Redis: %v", err)
		return nil, err
	} else {
		thumbnail, err := utils.FormThumbnailURL(video)
		if err != nil {
			s.Logger.Errorf("error while forming URL for thumbnail image: %v", err)
			return nil, err
		}

		err = s.Redis.Set(video.URL, thumbnail, 0).Err()
		if err != nil {
			s.Logger.Errorf("error caching thumbnail in Redis...")
		}

		s.Logger.Logln(logrus.InfoLevel, "Successfully generated thumbnail URL. Proceeding to saving image...")
		return &protos.Thumbnail{URL: thumbnail}, nil
	}
}

func (s *YouTubeServer) SaveThumbnail(ctx context.Context, thumbnail *protos.Thumbnail) (*protos.SaveStatus, error) {
	if thumbnail == nil {
		return &protos.SaveStatus{Status: _statusNotSaved}, errors.New("Something wrong with Thumbnail")
	}

	r := regexp.MustCompile(`vi/(.*?)/hqdefault.jpg`)
	filename := r.FindStringSubmatch(thumbnail.URL)[1] + ".jpg"

	file, err := os.Create(filepath.Join("../images", filename))
	if err != nil {
		s.Logger.Errorf("error creating image: %v", err)
		return &protos.SaveStatus{Status: _statusNotSaved}, err
	}
	defer file.Close()

	response, err := http.Get(thumbnail.URL)
	if err != nil {
		s.Logger.Errorf("error from GET request: %v", err)
		return &protos.SaveStatus{Status: _statusNotSaved}, err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		s.Logger.Errorf("error while copying body to file: %v", err)
		return &protos.SaveStatus{Status: _statusNotSaved}, err
	}

	s.Logger.Logln(logrus.InfoLevel, "Successfully saved image!")
	return &protos.SaveStatus{Status: _statusSaved}, nil
}

// Uber Go Style: Verify interface compliance at compile time
var _ protos.YouTubeThumbnailServer = (*YouTubeServer)(nil)

// Uber Go Style Exit Once Rule
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	lis, err := net.Listen("tcp", _defaultPort)
	if err != nil {
		return errors.New("failed to listen...")
	}

	os.Mkdir("./../images", os.ModePerm)

	grpcServer := grpc.NewServer()
	youtube := NewYouTubeServer()
	protos.RegisterYouTubeThumbnailServer(grpcServer, youtube)

	err = grpcServer.Serve(lis)
	if err != nil {
		return errors.New("grpc server failed to serve...")
	}

	return nil
}
