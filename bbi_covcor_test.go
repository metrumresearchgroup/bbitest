package babylontest

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestCovCorHappyPath(t *testing.T) {

	mod := `acop`


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

	if *update_summary {
		UpdateGoldenFile(gtd)
	}

	RequireOutputMatchesGoldenFile(gtd)
}



