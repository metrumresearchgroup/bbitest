package bbitest

import (
	"context"
	"encoding/json"
	"github.com/metrumresearchgroup/babylon/cmd"
	"github.com/metrumresearchgroup/babylon/configlib"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBBIConfigJSONCreated(t *testing.T){
	scenarios := InitializeScenarios([]string{
		"240",
		"acop",
		"ctl_test",
		"metrum_std",
	})


	for _, v := range scenarios{
		v.Prepare(context.Background())

		for _, m := range v.models {
			args := []string{
				"nonmem",
				"run",
				"local",
				"--nm_version",
				os.Getenv("NMVERSION"),
			}

			output, err := m.Execute(v,args...)

			assert.Nil(t,err)
			assert.NotNil(t,output)

			nmd := NonMemTestingDetails{
				t:         t,
				OutputDir: filepath.Join(v.Workpath,m.identifier),
				Model:     m,
				Output:    output,
			}

			AssertNonMemCompleted(nmd)
			AssertNonMemCreatedOutputFiles(nmd)
			AssertBBIConfigJSONCreated(nmd)
			AssertBBIConfigContainsSpecifiedNMVersion(nmd,os.Getenv("NMVERSION"))
		}
	}
}

func AssertBBIConfigContainsSpecifiedNMVersion(details NonMemTestingDetails, nmVersion string) {
	configFile, _ := os.Open(filepath.Join(details.OutputDir, "bbi_config.json"))
	cbytes, _ := ioutil.ReadAll(configFile)
	configFile.Close() //Go ahead and close the file handle

	nm := cmd.NonMemModel{}

	json.Unmarshal(cbytes, &nm)

	assert.NotNil(details.t,nm)
	assert.NotEqual(details.t,nm,cmd.NonMemModel{})
	assert.Equal(details.t,nm.Configuration.NMVersion,nmVersion)
}


func TestConfigValuesAreCorrectInWrittenFile(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate
	//Pick a few critical configuration components such as
	/*
		--clean_level 3
		--copy_level 1
		--debug
		--parallel=true <- make sure it's present
		--mpi_exec_path
	*/

	Scenario := InitializeScenarios([]string{
		"240",
	})[0]

	Scenario.Prepare(context.Background())

	commandAndArgs := []string{
		"--debug=true", //Needs to be in debug mode to generate the expected output
		"nonmem",
		"run",
		"--clean_lvl",
		"3",
		"--copy_lvl",
		"1",
		"--parallel=true",
		"--mpi_exec_path",
		os.Getenv("MPIEXEC_PATH"),
		"local",
		"--nm_version",
		os.Getenv("NMVERSION"),
	}

	for _, m := range Scenario.models {
		output, err := m.Execute(Scenario,commandAndArgs...)

		assert.Nil(t,err)
		assert.NotEmpty(t,output)

		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: filepath.Join(Scenario.Workpath,m.identifier),
			Model:     m,
			Output:    output,
		}

		AssertNonMemCompleted(nmd)
		AssertNonMemCreatedOutputFiles(nmd)
		AssertNonMemOutputContainsParafile(nmd)

		//Now read the Config Lib
		configFile := filepath.Join(Scenario.Workpath,m.identifier,"bbi.yaml")
		file, _ := os.Open(configFile)
		Config := configlib.Config{}
		bytes, _ := ioutil.ReadAll(file)
		err = yaml.Unmarshal(bytes,&Config)

		assert.Nil(t,err)

		assert.Equal(t,3,Config.CleanLvl)
		assert.Equal(t,1,Config.CopyLvl)
		assert.Equal(t, true,Config.Parallel,)
		assert.Equal(t,os.Getenv("NMVERSION"), Config.NMVersion)

		assert.Equal(t,os.Getenv("MPIEXEC_PATH"),Config.MPIExecPath )
		assert.Equal(t,false,Config.Overwrite)
	}

}