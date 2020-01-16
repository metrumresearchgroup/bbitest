package babylontest

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBabylonCompletesLocalExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate
	Initialize()

	fs := afero.NewOsFs()

	tarIdentifier := ".tar.gz"

	//Prep for Running this beast
	//Get all .gz files in the EXECUTION dir
	modelSets, _ := afero.Glob(fs,filepath.Join(EXECUTION_DIR,"*" + tarIdentifier))

	//Test shouldn't take longer than 5 min in total
	//TODO use the context downstream in a runModel function
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Minute)
	defer cancel()

	//TODO Break this into a method that takes a function for execution
	for _, v := range modelSets{
		file := filepath.Base(v)
		modelSet := strings.TrimSuffix(file,tarIdentifier)

		log.Infof("Beginning local execution test for model set %s",modelSet)

		//create Target directory as this untar operation doesn't handle it for you
		fs.MkdirAll(filepath.Join(EXECUTION_DIR,modelSet),0755)

		reader, _ := os.Open(v)

		Untar(filepath.Join(EXECUTION_DIR,modelSet),reader)

		reader.Close()

		os.Chdir(filepath.Join(EXECUTION_DIR,modelSet))
		executeCommand(ctx, "bbi", "init","--dir",viper.GetString("nonmemroot"))

		models := findModelFiles(filepath.Join(EXECUTION_DIR,modelSet))

		nmVersion, err := findNonMemKey(filepath.Join(EXECUTION_DIR,modelSet,"babylon.yaml"))

		if err != nil {
			log.Fatal("Unable to locate nonmem version to run bbi!")
		}

		for _ , m := range models {
			output := executeCommand(ctx, "bbi", "nonmem","run","local", "--nmVersion",nmVersion,m)
			assert.Contains(t,output,"Beginning local work")
			assert.Contains(t,output,"Beginning cleanup")

			modelName := strings.Split(filepath.Base(m),".")[0]
			outputDir := filepath.Join(EXECUTION_DIR,modelSet,modelName)

			xmlControlStream, err := afero.Exists(fs,filepath.Join(outputDir,modelName + ".xml"))

			assert.Nil(t,err)
			assert.True(t,xmlControlStream)

			//TODO Nonmem output file and look for completion text
		}

	}
}

func findNonMemKey(pathToBabylonConfig string) (string, error) {
	file, _ := os.Open(pathToBabylonConfig)
	defer file.Close()

	contentBytes, _ := ioutil.ReadAll(file)
	contentLines := strings.Split(string(contentBytes),"\n")

	for k, l := range contentLines{
		if l == "nonmem:" {
			return strings.TrimSpace(strings.Split(contentLines[k+1],":")[0]),nil
		}
	}

	return "", errors.New("unable to locate a key for a valid nonmem config")
}