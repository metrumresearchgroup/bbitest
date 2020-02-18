package babylontest

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNMQUALExecutionSucceeds(t *testing.T){
	scenarios := Initialize()

	//Let's work with third Scenario
	scenario := scenarios[2]

	scenario.Prepare(context.Background())

	for _, m := range scenario.models {
		args := []string{
			"nonmem",
			"run",
			"--nm_version",
			os.Getenv("NMVERSION"),
			"--nmqual=true",
			"local",
		}

		output, err := m.Execute(scenario, args...)

		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: filepath.Join(scenario.Workpath,m.identifier),
			Model:     m,
			Output:    output,
		}

		assert.Nil(t,err)
		AssertNonMemCompleted(nmd)
		AssertNonMemCreatedOutputFiles(nmd)
		AssertScriptContainsAutologReference(nmd)
	}
}

func AssertScriptContainsAutologReference(details NonMemTestingDetails){
	scriptFile, _  := os.Open(filepath.Join(details.OutputDir,details.Model.identifier + ".sh"))
	bytes, _ := ioutil.ReadAll(scriptFile)
	scriptFile.Close()
	scriptFileContent := string(bytes)

	assert.Contains(details.t,scriptFileContent,"autolog.pl")
}