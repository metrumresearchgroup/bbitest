package babylontest

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var CovCorTestMods = []string{
	"acop",
}

func TestCovCorHappyPath(t *testing.T) {

	for _, mod := range(CovCorTestMods) {
		commandAndArgs := []string{
			"nonmem",
			"covcor",
			filepath.Join(SUMMARY_TEST_DIR, mod, mod),
		}

		output, err := executeCommand(context.Background(),"bbi", commandAndArgs...)

		assert.Nil(t,err)
		assert.NotEmpty(t,output)

		gtd := GoldenFileTestingDetails{
			t:               t,
			outputString:    output,
			goldenFilePath:  filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden.covcor.json"),
		}

		if os.Getenv("UPDATE_SUMMARY") == "true" {
			UpdateGoldenFile(gtd)
		}

		RequireOutputMatchesGoldenFile(gtd)
	}
}



