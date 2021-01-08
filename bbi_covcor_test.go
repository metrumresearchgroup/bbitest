package bbitest

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

var CovCorTestMods = []string{
	"acop",
	"example2_itsimp",
	"1001",
}

func TestCovCorHappyPath(t *testing.T) {

	for _, mod := range(CovCorTestMods) {
		commandAndArgs := []string{
			"nonmem",
			"covcor",
			filepath.Join(SUMMARY_TEST_DIR, mod, mod),
		}

		output, err := executeCommand(context.Background(),"bbi", commandAndArgs...)

		require.Nil(t,err)
		require.NotEmpty(t,output)

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

var CovCorErrorMods = []string{
	"12",
	"iovmm",
}

func TestCovCorErrors (t *testing.T) {
	for _, tm := range(CovCorErrorMods) {

		commandAndArgs := []string{
			"nonmem",
			"covcor",
			filepath.Join(SUMMARY_TEST_DIR, tm, tm),
		}

		// try without flag and get error
		output, err := executeCommandNoErrorCheck(context.Background(),"bbi", commandAndArgs...)
		require.NotNil(t,err)
		errorMatch, _ := regexp.MatchString(noFilePresentError, output)
		require.True(t, errorMatch)
	}
}
