package babylontest

import (
	"context"
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
	v.Prepare(ctx)

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"local",
				"--nmVersion",
				v.nmversion,
			}

			m.Execute(v,nonMemArguments...)

			testingDetails := NonMemTestingDetails{
				t:         t,
				OutputDir: filepath.Join(v.Workpath,m.identifier),
				Model:     m,
			}

			AssertNonMemCompleted(testingDetails)
			AssertNonMemCreatedOutputFiles(testingDetails)
			AssertContainsBBIScript(t,testingDetails)
		}

	}
}

