package main

import (
	"context"
	"os"
	"testing"

	"github.com/deep0ne/grpc-youtube-thumbnail/protos"
)

func TestGetThumbnail(t *testing.T) {
	server := NewYouTubeServer()
	tests := []struct {
		name                 string
		video                *protos.Video
		expectedThumbnailURL string
	}{
		{
			name:                 "valid video URL",
			video:                &protos.Video{URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
			expectedThumbnailURL: "https://img.youtube.com/vi/dQw4w9WgXcQ/hqdefault.jpg",
		},
		{
			name:                 "invalid video URL",
			video:                &protos.Video{URL: "https://bullshit_URL"},
			expectedThumbnailURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thumbnail, _ := server.GetThumbnail(context.Background(), tt.video)
			if thumbnail != nil && thumbnail.URL != tt.expectedThumbnailURL {
				t.Errorf("expected thumbnail: %v, got %v", tt.expectedThumbnailURL, thumbnail.URL)
			}
		})
	}
}

func TestSaveThumbnail(t *testing.T) {
	server := NewYouTubeServer()
	os.Mkdir("./../images", os.ModePerm)
	tests := []struct {
		name               string
		video              *protos.Video
		expectedSaveStatus string
	}{
		{
			name:               "valid URL",
			video:              &protos.Video{URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
			expectedSaveStatus: "Image was saved to image/ folder successfully",
		},
		{
			name:               "invalid URL",
			video:              &protos.Video{URL: "https://bullshit_URL"},
			expectedSaveStatus: "Image was not saved",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thumbnail, _ := server.GetThumbnail(context.Background(), tt.video)
			saveStatus, _ := server.SaveThumbnail(context.Background(), thumbnail)
			if saveStatus != nil && saveStatus.Status != tt.expectedSaveStatus {
				t.Errorf("Saving thumbnail works wrong. Expected status: %v. Got: %v", tt.expectedSaveStatus, saveStatus.Status)
			}
		})
	}
}
