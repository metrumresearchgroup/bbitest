package babylontest

import (
	"context"
	"github.com/ghodss/yaml"
	"github.com/metrumresearchgroup/babylon/configlib"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestBabylonCompletesSGEExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate

	if ! FeatureEnabled("SGE"){
		t.Skip("Skipping SGE as it's not enabled")
	}

	scenarios := InitializeScenarios([]string{
		"240",
		"acop",
		"ctl_test",
	})

	whereami, _ := os.Getwd()

	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range scenarios{
		//log.Infof("Beginning SGE execution test for model set %s",v.identifier)
		v.Prepare(ctx)

		bbiBinary, _ := exec.LookPath("bbi")

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"sge",
				"--nm_version",
				os.Getenv("NMVERSION"),
				"--babylon_binary",
				bbiBinary,
			}

			_, err := m.Execute(v,nonMemArguments...)

			if err != nil {
				t.Error(err)
			}


			os.Chdir(filepath.Join(v.Workpath,m.identifier))
			//Now let's run the script that was generated
			_, err = executeCommand(ctx,filepath.Join(v.Workpath,m.identifier,"grid.sh"))
			os.Chdir(whereami)

			if err != nil {
				log.Error(err)
			}

			testingDetails := NonMemTestingDetails{
				t:         t,
				OutputDir: filepath.Join(v.Workpath,m.identifier),
				Model:     m,
			}

			AssertNonMemCompleted(testingDetails)
			AssertNonMemCreatedOutputFiles(testingDetails)
			AssertContainsBBIScript(testingDetails)
		}
	}
}


func TestBabylonCompletesParallelSGEExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate

	if ! FeatureEnabled("SGE"){
		t.Skip("Skipping SG Parallel execution as it's not enabled")
	}

	scenarios := InitializeScenarios([]string{
		"240",
		"acop",
		"ctl_test",
	})

	whereami, _ := os.Getwd()



	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range scenarios[0:3]{
		//log.Infof("Beginning SGE parallel execution test for model set %s",v.identifier)
		v.Prepare(ctx)

		bbiBinary, _ := exec.LookPath("bbi")

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"sge",
				"--nm_version",
				os.Getenv("NMVERSION"),
				"--babylon_binary",
				bbiBinary,
				"--parallel=true",
				"--mpi_exec_path",
				os.Getenv("MPIEXEC_PATH"),
			}

			_, err := m.Execute(v,nonMemArguments...)

			if err != nil {
				t.Error(err)
			}



			//Now let's run the script that was generated
			os.Chdir(filepath.Join(v.Workpath,m.identifier))
			_, err = executeCommand(ctx,filepath.Join(v.Workpath,m.identifier,"grid.sh"))
			os.Chdir(whereami)

			if err != nil {
				log.Error(err)
			}

			testingDetails := NonMemTestingDetails{
				t:         t,
				OutputDir: filepath.Join(v.Workpath,m.identifier),
				Model:     m,
			}

			AssertNonMemCompleted(testingDetails)
			AssertNonMemCreatedOutputFiles(testingDetails)
			AssertContainsBBIScript(testingDetails)
			AssertNonMemOutputContainsParafile(testingDetails)
		}
	}
}

func TestConfigValuesAreCorrectInWrittenFile(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate

	const qsub string = "/usr/local/bin/qsub"
	purgeBinary(qsub)
	fakeBinary(qsub)

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
		configFile := filepath.Join(Scenario.Workpath,m.identifier,"babylon.yaml")
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

	purgeBinary(qsub)
}




func fakeBinary(name string) {
	contents := `#!/bin/bash
	echo $0 $@
	exit 0`

	ioutil.WriteFile(name, []byte(contents), 0755)
}

func purgeBinary(name string) {
	os.Remove(name)
}
