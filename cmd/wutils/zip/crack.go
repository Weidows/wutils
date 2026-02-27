package zip

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Weidows/wutils/pkg/zip"
)

const (
	defaultMaxWorkers = 50
	configDir         = ".config"
	wutilsDir         = "wutils"
	passwordFile      = "password-dict.txt"
)

func getHomePasswordFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	dir := filepath.Join(home, configDir, wutilsDir)
	path := filepath.Join(dir, passwordFile)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		defaultContent := []byte("# Add your password dictionary here, one password per line\n" +
			"test\n123456\npassword\n" +
			"admin\n12345678\n")
		if err := os.WriteFile(path, defaultContent, 0644); err != nil {
			return "", fmt.Errorf("failed to create password file: %w", err)
		}
		fmt.Printf("Created default password dictionary at: %s\n", path)
	}
	return path, nil
}

func CrackPassword(archivePath string) string {
	passwordFilePath, err := getHomePasswordFilePath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}

	passwords, err := os.ReadFile(passwordFilePath)
	if err != nil {
		fmt.Printf("Error: failed to read password file: %v\n", err)
		return ""
	}

	passwordList := parsePasswords(string(passwords))
	if len(passwordList) == 0 {
		fmt.Println("Error: no valid passwords found in dictionary")
		fmt.Printf("Please add passwords to: %s\n", passwordFilePath)
		return ""
	}

	result := crackWithWorkers(archivePath, passwordList, defaultMaxWorkers)
	if result != "" {
		fmt.Printf("Password found: %s\n", result)
	} else {
		fmt.Println("Password not found in dictionary")
	}
	return result
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

			archive := zip.NewArchive(archivePath, p)
			if archive.TryUnzip() {
				select {
				case success <- p:
				default:
				}
			}

			count := tried.Add(1)
			percent := float64(count) / float64(total) * 100
			fmt.Printf("\rProgress: %d/%d (%.1f%%)", count, total, percent)
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
