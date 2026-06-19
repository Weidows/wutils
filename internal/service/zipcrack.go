package service

import (
	"github.com/Weidows/wutils/cmd/wutils/zip"
	"github.com/Weidows/wutils/internal/app"
)

// ZipCrackService cracks the password of encrypted archives (.zip/.7z/.rar/.tar.gz/.tar.bz2).
type ZipCrackService struct{}

// NewZipCrackService creates a new ZipCrackService.
func NewZipCrackService() *ZipCrackService {
	return &ZipCrackService{}
}

func (s *ZipCrackService) Name() string        { return "zipcrack" }
func (s *ZipCrackService) Description() string { return "压缩包密码破解" }
func (s *ZipCrackService) Status() app.ServiceStatus { return app.StatusStopped }
func (s *ZipCrackService) Start() error               { return nil }
func (s *ZipCrackService) Stop() error                { return nil }

// CrackPassword attempts to find the archive password using the default dictionary.
func (s *ZipCrackService) CrackPassword(archivePath string) string {
	return zip.CrackPassword(archivePath)
}

// CrackPasswordWithList attempts to find the archive password using a provided password list.
func (s *ZipCrackService) CrackPasswordWithList(archivePath string, passwords []string) string {
	return zip.CrackPasswordWithList(archivePath, passwords)
}
