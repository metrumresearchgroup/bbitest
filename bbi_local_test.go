package babylontest

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBabylonCompletesLocalExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate
	scenarios := Initialize()


	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range scenarios{
		log.Infof("Beginning local execution test for model set %s",v.identifier)
		v.Prepare(ctx)

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"local",
				"--nmVersion",
				os.Getenv("NMVERSION"),
			}

			err := m.Execute(v,nonMemArguments...)

			if err != nil {
				t.Error(err)
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


func TestBabylonParallelExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate
	scenarios := Initialize()

	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range scenarios{
		log.Infof("Beginning localized parallel execution test for model set %s",v.identifier)
		v.Prepare(ctx)



		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"local",
				"--nmVersion",
				os.Getenv("NMVERSION"),
				"--parallel=true",
				"--mpiExecPath",
				os.Getenv("MPIEXEC_PATH"),
			}

			err := m.Execute(v,nonMemArguments...)

			if err != nil {
				t.Error(err)
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

