package babylontest

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io"
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
		//log.Infof("Beginning local execution test for model set %s",v.identifier)
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

			_, err := m.Execute(v,nonMemArguments...)

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

func TestNMFEOptionsEndInScript(t *testing.T){
	scenarios := Initialize()

	whereami, _ := os.Getwd()

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
				os.Getenv("NMVERSION"),
				"--background=true",
				"--prcompile=true",
			}

			_, err := m.Execute(v,nonMemArguments...)

			if err != nil {
				t.Error(err)
			}



			//Now let's run the script that was generated
			os.Chdir(filepath.Join(v.Workpath,m.identifier))
			_, err = executeCommand(ctx,filepath.Join(v.Workpath,m.identifier,m.identifier + ".sh"))
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
			AssertContainsNMFEOptions(testingDetails,filepath.Join(testingDetails.OutputDir,m.identifier+".sh"),"-background")
			AssertContainsNMFEOptions(testingDetails,filepath.Join(testingDetails.OutputDir,m.identifier+".sh"),"-prcompile")
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
		//log.Infof("Beginning localized parallel execution test for model set %s",v.identifier)
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

			_, err := m.Execute(v,nonMemArguments...)

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

func TestDefaultConfigLoaded(t *testing.T){
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()
	scenarios := Initialize()
	//Only work on the first one.
	scenario := scenarios[0]

	nonMemArguments := []string{
		"-d",
		"nonmem",
		"run",
		"local",
		"--nmVersion",
		os.Getenv("NMVERSION"),
	}

	scenario.Prepare(ctx)

	for _, v := range scenario.models {
		out, _ := v.Execute(scenario,nonMemArguments...)
		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: "",
			Model:     v,
			Output:    out,
		}

		AssertDefaultConfigLoaded(nmd)
	}
}

func TestSpecifiedConfigByAbsPathLoaded(t *testing.T){
	fs := afero.NewOsFs()

	if ok, _  := afero.DirExists(fs, "/tmp/meow"); ok {
		fs.RemoveAll("/tmp/meow")
	}



	fs.MkdirAll("/tmp/meow",0755)
	//Copy the babylon file here
	source, _ := fs.Open("babylon.yaml")
	defer source.Close()
	dest, _ := fs.Create("/tmp/meow/babylon.yaml")
	defer dest.Close()

	io.Copy(dest,source)


	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()
	scenarios := Initialize()
	//Only work on the first one.
	scenario := scenarios[0]

	//Copy config to /tmp/meow/babylon.yaml


	nonMemArguments := []string{
		"-d",
		"--config",
		"/tmp/meow/babylon.yaml",
		"nonmem",
		"run",
		"local",
		"--nmVersion",
		os.Getenv("NMVERSION"),
	}

	scenario.Prepare(ctx)

	for _, v := range scenario.models {
		out, _ := v.Execute(scenario,nonMemArguments...)
		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: "",
			Model:     v,
			Output:    out,
		}

		AssertSpecifiedConfigLoaded(nmd,"/tmp/meow/babylon.yaml")
	}
}

func TestSpecifiedConfigByRelativePathLoaded(t *testing.T){
	fs := afero.NewOsFs()

	if ok, _  := afero.DirExists(fs, "tmp/meow"); ok {
		fs.RemoveAll("tmp/meow")
	}



	fs.MkdirAll("tmp/meow",0755)
	//Copy the babylon file here
	source, _ := fs.Open("babylon.yaml")
	defer source.Close()
	dest, _ := fs.Create("tmp/meow/babylon.yaml")
	defer dest.Close()

	io.Copy(dest,source)


	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()
	scenarios := Initialize()
	//Only work on the first one.
	scenario := scenarios[0]

	//Copy config to /tmp/meow/babylon.yaml


	nonMemArguments := []string{
		"-d",
		"--config",
		"tmp/meow/babylon.yaml",
		"nonmem",
		"run",
		"local",
		"--nmVersion",
		os.Getenv("NMVERSION"),
	}

	scenario.Prepare(ctx)

	for _, v := range scenario.models {
		out, _ := v.Execute(scenario,nonMemArguments...)
		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: "",
			Model:     v,
			Output:    out,
		}

		AssertSpecifiedConfigLoaded(nmd,"tmp/meow/babylon.yaml")
	}
}

