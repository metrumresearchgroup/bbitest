package bbitest

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//Verifies that if we have a CTL file we don't experience issues with path location of the data file
func TestHasValidDataPathForCTL(t *testing.T){
	scenarios := InitializeScenarios([]string{
		"ctl_test",
	})

	//Take the 3rd scenario for the CTL file
	scenario := scenarios[0]

	scenario.Prepare(context.Background())

	//Directories et all should be prepared.
	for _, m := range scenario.models {

		t.Run(fmt.Sprintf("validPathCTL_%s",m.filename),func(t *testing.T){
			args := []string{
				"nonmem",
				"run",
				"local",
				"--nm_version",
				os.Getenv("NMVERSION"),
			}

			output, err := m.Execute(scenario,args...)

			ntd := NonMemTestingDetails{
				t:         t,
				OutputDir:  filepath.Join(scenario.Workpath,m.identifier),
				Model:     m,
				Output:    output,
			}



			require.Nil(t,err)
			AssertNonMemCompleted(ntd)
			AssertNonMemCreatedOutputFiles(ntd)
		})
	}
}


//Verifies that if we have a CTL file we don't experience issues with path location of the data file
func TestHasInvalidDataPath(t *testing.T){
	scenarios := InitializeScenarios([]string{
		"ctl_test",
	})

	//Take the 3rd scenario for the CTL file
	scenario := scenarios[0]

	scenario.Prepare(context.Background())

	//Directories et all should be prepared.
	for _, m := range scenario.models {
		//We need to manipulate the file to contain an invalid file reference
		file, _ := os.Open(filepath.Join(scenario.Workpath,m.filename))
		b, _ := ioutil.ReadAll(file)
		file.Close() //Explicitly close so we can write it again
		lines := strings.Split(string(b),"\n")

		for k, line := range lines {
			if strings.Contains(line,"$DATA") {
				lines[k] = "$DATA      ../FData.csv IGNORE=@"
			}
		}

		adjusted := strings.Join(lines,"\n")
		ab := []byte(adjusted)

		err := ioutil.WriteFile(filepath.Join(scenario.Workpath,m.filename),ab,0755)

		if err != nil {
			t.Log("Had a problem writing the file")
		}

		t.Run(fmt.Sprintf("invalidPathCTL_%s",m.filename),func(t *testing.T){
			args := []string{
				"nonmem",
				"run",
				"local",
				"--nm_version",
				os.Getenv("NMVERSION"),
			}

			_, err := m.Execute(scenario,args...)

			//ntd := NonMemTestingDetails{
			//	t:         t,
			//	OutputDir:  filepath.Join(scenario.Workpath,m.identifier),
			//	Model:     m,
			//	Output:    output,
			//}


			require.NotNil(t,err)
			require.Error(t,err)

		})
	}
}


//Verifies that if we have a CTL file we don't experience issues with path location of the data file
func TestHasValidComplexPathCTLAndMod(t *testing.T){
	scenarios := InitializeScenarios([]string{
		"metrum_std",
	})

	//Take the 3rd scenario for the CTL file
	scenario := scenarios[0]

	scenario.Prepare(context.Background())

	//Because we have a relatively complex model Structure we're going to set the model manually
	scenario.models = modelsFromOriginalScenarioPath(filepath.Join(scenario.Workpath,"model","nonmem","test_suite_1"))
	intermediary := scenario
	intermediary.Workpath = filepath.Join(scenario.Workpath,"model","nonmem","test_suite_1")

	//Directories et all should be prepared.
	for _, m := range scenario.models {

		t.Run(fmt.Sprintf("validComplexPathFor_%s",m.filename),func(t *testing.T){
			args := []string{
				"nonmem",
				"run",
				"local",
				"--nm_version",
				os.Getenv("NMVERSION"),
			}



			output, err := m.Execute(intermediary,args...)

			println(output)

			ntd := NonMemTestingDetails{
				t:         t,
				OutputDir:  filepath.Join(intermediary.Workpath,m.identifier),
				Model:     m,
				Output:    output,
			}


			require.Nil(t,err)
			AssertNonMemCompleted(ntd)
			AssertNonMemCreatedOutputFiles(ntd)
		})
	}
}

