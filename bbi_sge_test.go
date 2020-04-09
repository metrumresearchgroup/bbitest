package babylontest

import (
	"context"
	"github.com/metrumresearchgroup/gogridengine"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
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

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"sge",
				"--nm_version",
				os.Getenv("NMVERSION"),
			}

			_, err := m.Execute(v,nonMemArguments...)

			if err != nil {
				t.Error(err)
			}

			WaitForSGEToTerminate(v)

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

		for _ , m := range v.models {

			nonMemArguments := []string{
				"-d",
				"nonmem",
				"run",
				"sge",
				"--nm_version",
				os.Getenv("NMVERSION"),
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






func fakeBinary(name string) {
	contents := `#!/bin/bash
	echo $0 $@
	exit 0`

	ioutil.WriteFile(name, []byte(contents), 0755)
}

func purgeBinary(name string) {
	os.Remove(name)
}


func WaitForSGEToTerminate(scenario *Scenario) {
	for CountOfPendingJobs() > 0 {
		log.Infof("Located %d pending jobs. Waiting for 30 seconds to check again", CountOfPendingJobs())
		time.Sleep(30 * time.Second)
	}

	log.Info("Looks like all queued and running jobs have terminated")
}

func CountOfPendingJobs() int {
	jobs, _ := gogridengine.GetJobsWithFilter(func(j gogridengine.Job) bool {
		return j.State == "qw" || j.State == "r"
	})

	return len(jobs)
}