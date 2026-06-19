package service

import (
	"github.com/Weidows/wutils/cmd/wutils/media"
	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
)

// MediaService groups photos and videos by time and GPS proximity.
type MediaService struct{}

// NewMediaService creates a new MediaService.
func NewMediaService() *MediaService {
	return &MediaService{}
}

func (s *MediaService) Name() string        { return "media" }
func (s *MediaService) Description() string { return i18n.G("media.description") }
func (s *MediaService) Status() app.ServiceStatus { return app.StatusStopped }
func (s *MediaService) Start() error               { return nil }
func (s *MediaService) Stop() error                { return nil }

// ClusterAndCopy clusters media files in inputDir and copies them to an output directory.
func (s *MediaService) ClusterAndCopy(inputDir string) {
	media.ClusterAndCopy(inputDir)
}
