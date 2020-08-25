package babylontest

import (
	"context"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

type testConfig struct {
	bbiFlag string
	goldenExt string
}

var testConfigs = []testConfig{
	{
		"--json",
		".json",
	},
	{
		"",
		".txt",
	},
}

var SummaryHappyPathTestMods = []string{
	"acop",              // basic model
	"12",                // bootstrap model with no $COV step
	"example2_saemimp",  // two est methods SAEM => IMP
	"example2_itsimp",   // two est methods ITS => IMP (No Prior)
	"example2_bayes",    // Bayes (5 est methods, from NONMEM examples)
}

func TestSummaryHappyPath(t *testing.T) {
	for _, mod := range(SummaryHappyPathTestMods) {
		for _, tc := range(testConfigs) {

			commandAndArgs := []string{
				"nonmem",
				"summary",
				filepath.Join(SUMMARY_TEST_DIR, mod, mod+".lst"),
			}

			if (tc.bbiFlag != "") {
				commandAndArgs = append(commandAndArgs, tc.bbiFlag)
			}

			output, err := executeCommand(context.Background(),"bbi", commandAndArgs...)

			require.Nil(t,err)
			require.NotEmpty(t,output)

			gtd := GoldenFileTestingDetails{
				t:               t,
				outputString:    output,
				goldenFilePath:  filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden"+tc.goldenExt),
			}

			if *update_summary {
				UpdateGoldenFile(gtd)
			}

			RequireOutputMatchesGoldenFile(gtd)
		}
	}
}


type testModWithFlag struct {
	mod     string
	bbiFlag string
}

var SummaryFlagsTestMods = []testModWithFlag{
	{"66",      "--no-shk-file"},
	{"1001",    "--ext-file=1001.1.TXT"},
}

func TestSummaryFlags(t *testing.T) {
	for _, tm := range(SummaryFlagsTestMods) {
		for _, tc := range(testConfigs) {

			mod := tm.mod

			commandAndArgs := []string{
				"nonmem",
				"summary",
				filepath.Join(SUMMARY_TEST_DIR, mod, mod+".lst"),
			}

			commandAndArgs = append(commandAndArgs, tm.bbiFlag)
			if (tc.bbiFlag != "") {
				commandAndArgs = append(commandAndArgs, tc.bbiFlag)
			}

			output, err := executeCommand(context.Background(),"bbi", commandAndArgs...)

			require.Nil(t,err)
			require.NotEmpty(t,output)

			gtd := GoldenFileTestingDetails{
				t:               t,
				outputString:    output,
				goldenFilePath:  filepath.Join(SUMMARY_TEST_DIR, SUMMARY_GOLD_DIR, mod+".golden"+tc.goldenExt),
			}

			if *update_summary {
				UpdateGoldenFile(gtd)
			}

			RequireOutputMatchesGoldenFile(gtd)
		}
	}
}
