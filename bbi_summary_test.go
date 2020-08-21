package babylontest

import (
	//"bufio"
	//"bytes"
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	"io/ioutil"

	//"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

const TEST_DIR = "testdata/bbi_summary"
const GOLD_DIR = "aa_golden_files"

var update = flag.Bool("update", false, "update .golden files")

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

func TestSummaryHappyPath(t *testing.T) {

	mod := `acop`

	for _, tc := range(testConfigs) {

		commandAndArgs := []string{
			"nonmem",
			"summary",
			filepath.Join(TEST_DIR, mod, mod+".lst"),
		}

		if (tc.bbiFlag != "") {
			commandAndArgs = append(commandAndArgs, tc.bbiFlag)
		}

		output, err := executeCommand(context.Background(),"bbi", commandAndArgs...)

		assert.Nil(t,err)
		assert.NotEmpty(t,output)

		gf := mod+".golden"+tc.goldenExt
		gp := filepath.Join(TEST_DIR, GOLD_DIR, gf)
		if *update {
			t.Log("updating golden file "+gf)
			if err := ioutil.WriteFile(gp, []byte(output), 0644); err != nil {
				t.Fatalf("failed to update golden file: %s", err)
			}
		}
		g, err := ioutil.ReadFile(gp)
		if err != nil {
			t.Fatalf("failed reading .golden: %s", err)
		}

		//t.Log(output[1:100])
		//t.Log("===== output ^ ======= gold >")
		//t.Log(string(g)[1:100])
		if string(g) != output {
			t.Log("Diff between golden file and output:")
			t.Log(diff.Diff(string(g), output))
			t.Errorf("output does not match .golden file")
		}
	}
}



