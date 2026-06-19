package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
	pkgzip "github.com/Weidows/wutils/pkg/zip"
)

const (
	defaultMaxWorkers = 50
	passwordFile      = "password-dict.txt"
)

// ZipCrackService cracks the password of encrypted archives.
type ZipCrackService struct{}

// NewZipCrackService creates a new ZipCrackService.
func NewZipCrackService() *ZipCrackService {
	return &ZipCrackService{}
}

func (s *ZipCrackService) Name() string               { return "zipcrack" }
func (s *ZipCrackService) Description() string        { return i18n.G("zip.description") }
func (s *ZipCrackService) Status() app.ServiceStatus  { return app.StatusStopped }
func (s *ZipCrackService) Start() error               { return nil }
func (s *ZipCrackService) Stop() error                { return nil }

// CrackPassword attempts to find the archive password using the default dictionary.
func (s *ZipCrackService) CrackPassword(archivePath string) string {
	return s.CrackPasswordWithList(archivePath, nil)
}

// CrackPasswordWithList attempts to find the archive password using a provided password list.
func (s *ZipCrackService) CrackPasswordWithList(archivePath string, passwords []string) string {
	if passwords == nil {
		pwPath, err := passwordFilePath()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return ""
		}
		data, err := os.ReadFile(pwPath)
		if err != nil {
			fmt.Printf("Error: failed to read password file: %v\n", err)
			return ""
		}
		passwords = parsePasswords(string(data))
		if len(passwords) == 0 {
			fmt.Println("Error: no valid passwords found in dictionary")
			return ""
		}
	}

	result := crackWithWorkers(archivePath, passwords, defaultMaxWorkers)
	if result != "" {
		fmt.Printf("Password found: %s\n", result)
	} else {
		fmt.Println("Password not found in dictionary")
	}
	return result
}

func passwordFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "wutils")
	path := filepath.Join(dir, passwordFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		defaultContent := []byte("# Add your password dictionary here, one password per line\ntest\n123456\npassword\nadmin\n12345678\n")
		if err := os.WriteFile(path, defaultContent, 0644); err != nil {
			return "", fmt.Errorf("failed to create password file: %w", err)
		}
		fmt.Printf("Created default password dictionary at: %s\n", path)
	}
	return path, nil
}

func parsePasswords(content string) []string {
	sep := "\n"
	if strings.Contains(content, "\r\n") {
		sep = "\r\n"
	}
	lines := strings.Split(content, sep)
	passwords := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			passwords = append(passwords, line)
		}
	}
	return passwords
}

func crackWithWorkers(archivePath string, passwords []string, maxWorkers int) string {
	archivePath = pkgzip.ResolveArchivePath(archivePath)

	var wg sync.WaitGroup
	pool := make(chan struct{}, maxWorkers)
	success := make(chan string, 1)
	var tried atomic.Int64
	total := int64(len(passwords))

	for _, password := range passwords {
		wg.Add(1)
		pool <- struct{}{}
		go func(p string) {
			defer wg.Done()
			defer func() { <-pool }()
			select {
			case <-success:
				return
			default:
			}
			archive := pkgzip.NewArchive(archivePath, p)
			if archive.TryUnzip() {
				select {
				case success <- p:
				default:
				}
			}
			count := tried.Add(1)
			fmt.Printf("\rProgress: %d/%d (%.1f%%)", count, total, float64(count)/float64(total)*100)
		}(password)
	}
	wg.Wait()
	close(pool)
	close(success)
	fmt.Println()
	select {
	case result := <-success:
		return result
	default:
		return ""
	}
}
