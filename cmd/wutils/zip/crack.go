package zip

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/Weidows/wutils/pkg/zip"
)

func CrackPassword(archivePath string) string {
	passwords, err := ioutil.ReadFile("password.txt")
	if err != nil {
		fmt.Errorf("failed to read password file: %v", err)
		return ""
	}

	passwordsString := string(passwords)
	sep := "\n"
	if strings.Contains(passwordsString, "\r\n") {
		sep = "\r\n"
	}

	passwordList := strings.Split(string(passwords), sep)
	var wg sync.WaitGroup
	success := make(chan string)

	for _, password := range passwordList {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if zip.NewArchive(archivePath, p).TryUnzip() {
				success <- p
			} else {
				fmt.Printf("failed to unzip with password: %s\n", p)
			}
		}(password)
	}

	go func() {
		wg.Wait()
		close(success)
	}()

	for p := range success {
		return p
	}
	return ""
}
