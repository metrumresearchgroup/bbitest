package babylontest

import (
	"context"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNMQUALExecutionSucceeds(t *testing.T){


	if !FeatureEnabled("NMQUAL"){
		t.Skip("Testing for NMQUAL not enabled")
	}

	scenarios := InitializeScenarios([]string{
		"ctl_test",
	})

	//Let's work with third Scenario
	scenario := scenarios[0]

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
			Scenario: scenario,
		}

		require.Nil(t,err)
		AssertNonMemCompleted(nmd)
		AssertNonMemCreatedOutputFiles(nmd)
		AssertScriptContainsAutologReference(nmd)
		AssertDataSourceIsHashedAndCorrect(nmd)
		AssertModelIsHashedAndCorrect(nmd)
	}
}

//This test targets a model with a .mod extension to make sure that
//After cloning and re-creating as a .ctl, that the application
//knows to look for what was originally there; the .mod file
func TestHashingForNMQualWorksWithOriginalModFile(t *testing.T){

	if !FeatureEnabled("NMQUAL"){
		t.Skip("Testing for NMQUAL not enabled")
	}


	scenarios := InitializeScenarios([]string{
		"240",
	})

	//Let's work with third Scenario
	scenario := scenarios[0]

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
			Scenario: scenario,
		}

		require.Nil(t,err)
		AssertNonMemCompleted(nmd)
		AssertNonMemCreatedOutputFiles(nmd)
		AssertScriptContainsAutologReference(nmd)
		AssertDataSourceIsHashedAndCorrect(nmd)
		AssertModelIsHashedAndCorrect(nmd)
	}
}

func AssertScriptContainsAutologReference(details NonMemTestingDetails){
	scriptFile, _  := os.Open(filepath.Join(details.OutputDir,details.Model.identifier + ".sh"))
	bytes, _ := ioutil.ReadAll(scriptFile)
	scriptFile.Close()
	scriptFileContent := string(bytes)

	require.Contains(details.t,scriptFileContent,"autolog.pl")
}


