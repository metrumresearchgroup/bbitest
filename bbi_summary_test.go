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
	testConfig{
		"--json",
		".json",
	},
	testConfig{
		"",
		".txt",
	},
}

var SummaryHappyPathTestMods = []string{
	"acop",              // basic model
	"12",                // bootstrap model with no $COV step
	"example2_saemimp",  // two est methods SAEM => IMP
	"example2_itsimp",   // two est methods ITS => IMP (No Prior)
	//"example2_bayes",    // Bayes (5 est methods) CURRENTLY DIES ON no flag table
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



