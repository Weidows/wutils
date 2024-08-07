package keep_runner

import (
	"testing"

	"github.com/Weidows/wutils/utils/log"
)

func Test_ol(t *testing.T) {
	NewKeepRunner(log.GetLogger()).OlList()
}
