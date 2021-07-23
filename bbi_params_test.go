package bbitest

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

var testConfigsParams = []testConfig{
	{
		"--json",
		".json",
	},
	{
		"",
		".csv",
	},
}

func TestParamsSingleModel(t *testing.T) {
	mod := "example2_bayes"
	for _, tc := range(testConfigsParams) {

		commandAndArgs := []string{
			"nonmem",
			"params",
			filepath.Join(SUMMARY_TEST_DIR, mod),
		}

		if tc.bbiOption != "" {
			commandAndArgs = append(commandAndArgs, tc.bbiOption)
		}

		output, err := executeCommand(context.Background(),"bbi", commandAndArgs...)

		require.Nil(t,err)
		require.NotEmpty(t,output)

		gtd := GoldenFileTestingDetails{
			t:               t,
			outputString:    output,
			goldenFilePath:  filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden.params"+tc.goldenExt),
		}

		if os.Getenv("UPDATE_SUMMARY") == "true" {
			UpdateGoldenFile(gtd)
		}

		RequireOutputMatchesGoldenFile(gtd)
	}
}

func TestParamsDir(t *testing.T) {
	for _, tc := range(testConfigsParams) {

		commandAndArgs := []string{
			"nonmem",
			"params",
			"--dir",
			SUMMARY_TEST_DIR,
		}

		if tc.bbiOption != "" {
			commandAndArgs = append(commandAndArgs, tc.bbiOption)
		}

		output, err := executeCommand(context.Background(),"bbi", commandAndArgs...)

		require.Nil(t,err)
		require.NotEmpty(t,output)

		gtd := GoldenFileTestingDetails{
			t:               t,
			outputString:    output,
			goldenFilePath:  filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, "dir_bbi_summary.golden.params"+tc.goldenExt),
		}

		if os.Getenv("UPDATE_SUMMARY") == "true" {
			UpdateGoldenFile(gtd)
		}

		RequireOutputMatchesGoldenFile(gtd)
	}
}

