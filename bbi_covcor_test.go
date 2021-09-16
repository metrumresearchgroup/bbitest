package bbitest

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

var CovCorTestMods = []string{
	"acop",
	"example2_itsimp",
	"1001",
}

func TestCovCorHappyPath(tt *testing.T) {
	t := wrapt.WrapT(tt)

	for _, mod := range CovCorTestMods {
		commandAndArgs := []string{
			"nonmem",
			"covcor",
			filepath.Join(SUMMARY_TEST_DIR, mod, mod),
		}

		output, err := executeCommand(context.Background(), "bbi", commandAndArgs...)

		t.R.NoError(err)
		t.R.NotEmpty(output)

		gtd := GoldenFileTestingDetails{
			outputString:   output,
			goldenFilePath: filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden.covcor.json"),
		}

		if os.Getenv(":q") == "true" {
			UpdateGoldenFile(t, gtd)
		}

		RequireOutputMatchesGoldenFile(t, gtd)
	}
}

var CovCorErrorMods = []string{
	"12",
	"iovmm",
}

func TestCovCorErrors(tt *testing.T) {
	t := wrapt.WrapT(tt)

	for _, tm := range CovCorErrorMods {
		commandAndArgs := []string{
			"nonmem",
			"covcor",
			filepath.Join(SUMMARY_TEST_DIR, tm, tm),
		}

		// try without flag and get error
		output, err := executeCommandNoErrorCheck(context.Background(), "bbi", commandAndArgs...)
		t.R.Error(err)

		errorMatch, _ := regexp.MatchString(noFilePresentError, output)
		t.R.True(errorMatch)
	}
}
