package babylontest

import (
	"context"
	log "github.com/sirupsen/logrus"
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

	const qsub string = "/usr/local/bin/qsub"
	purgeBinary(qsub)
	fakeBinary(qsub)

	scenarios := Initialize()



	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range scenarios{
		log.Infof("Beginning SGE execution test for model set %s",v.identifier)
		v.Prepare(ctx)

		bbiBinary, _ := exec.LookPath("bbi")

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"sge",
				"--nmVersion",
				os.Getenv("NMVERSION"),
				"--babylonBinary",
				bbiBinary,
			}

			m.Execute(v,nonMemArguments...)



			//Now let's run the script that was generated
			executeCommand(ctx,filepath.Join(v.Workpath,m.identifier,"grid.sh"))

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

	purgeBinary(qsub)
}


func TestBabylonCompletesParallelSGEExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate

	const qsub string = "/usr/local/bin/qsub"
	purgeBinary(qsub)
	fakeBinary(qsub)

	scenarios := Initialize()



	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range scenarios{
		log.Infof("Beginning SGE parallel execution test for model set %s",v.identifier)
		v.Prepare(ctx)

		bbiBinary, _ := exec.LookPath("bbi")

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"sge",
				"--nmVersion",
				os.Getenv("NMVERSION"),
				"--babylonBinary",
				bbiBinary,
				"--parallel=true",
				"--mpiExecPath",
				os.Getenv("MPIEXEC_PATH"),
			}

			m.Execute(v,nonMemArguments...)



			//Now let's run the script that was generated
			executeCommand(ctx,filepath.Join(v.Workpath,m.identifier,"grid.sh"))

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

	purgeBinary(qsub)
}


func fakeBinary(name string) {
	contents := `#!/bin/bash
	echo $0 $@
	exit 0`

	err := ioutil.WriteFile(name, []byte(contents), 0755)
	if err != nil {
		log.Error("Unable to create the file", err)
	}
}

func purgeBinary(name string) {
	err := os.Remove(name)

	if err != nil {
		log.Error("Unable to create the file", err)
	}
}

