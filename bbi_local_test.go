package babylontest

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBabylonCompletesLocalExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate
	scenarios := Initialize()
	fs := afero.NewOsFs()

	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range scenarios{
		modelSet := v.identifier

		log.Infof("Beginning local execution test for model set %s",modelSet)

		//create Target directory as this untar operation doesn't handle it for you
		fs.MkdirAll(v.Workpath,0755)

		reader, err := os.Open(filepath.Join(EXECUTION_DIR,v.archive))

		if err != nil{
			log.Errorf("An error occurred during the untar operation: %s", err)
		}

		Untar(v.Workpath,reader)

		reader.Close()

		os.Chdir(v.Workpath)
		executeCommand(ctx, "bbi", "init","--dir",viper.GetString("nonmemroot"))


		//TODO Import babylon configlib and serialize into Config struct. This will let us sanely iterate and just Pick one as opposed to file manipulation garbage
		nmVersion, err := findNonMemKey(filepath.Join(EXECUTION_DIR,modelSet,"babylon.yaml"))

		if err != nil {
			log.Fatal("Unable to locate nonmem version to run bbi!")
		}

		for _ , m := range v.models {
			executeCommand(ctx, "bbi", "nonmem","run","local", "--nmVersion",nmVersion,filepath.Join(v.Workpath,m.filename))

			testingDetails := NonMemTestingDetails{
				t:         t,
				OutputDir: filepath.Join(v.Workpath,m.identifier),
				Model:     m,
			}

			AssertNonMemCompleted(testingDetails)
			AssertNonMemCreatedOutputFiles(testingDetails)
		}

	}
}

