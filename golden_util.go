package babylontest

import (
	"io/ioutil"
	"github.com/stretchr/testify/require"
	"testing"
)


type GoldenFileTestingDetails struct {
	t *testing.T
	outputString string
	goldenFilePath string
}

// check that string in outputString matches contents of file at goldenFilePath
func RequireOutputMatchesGoldenFile(details GoldenFileTestingDetails) {
	gb, err := ioutil.ReadFile(details.goldenFilePath)
	if err != nil {
		details.t.Fatalf("failed reading %s: %s", details.goldenFilePath, err)
	}
	gold := string(gb)

	require.Equal(details.t, gold, details.outputString, "output does not match .golden file "+details.goldenFilePath)
}

// Write string in outputString to file at goldenFilePath.
// If file already exists, it will be overwritten.
// User can then use git diff to see what has been updated.
func UpdateGoldenFile(details GoldenFileTestingDetails) {
	details.t.Logf("updating golden file %s", details.goldenFilePath)
	if err := ioutil.WriteFile(details.goldenFilePath, []byte(details.outputString), 0644); err != nil {
		details.t.Fatalf("failed to update %s: %s", details.goldenFilePath, err)
	}
}
