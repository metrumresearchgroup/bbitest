package babylontest

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBabylonCompletesLocalExecution(t *testing.T){
	//Get BB and make sure we have the test data moved over.
	//Clean Slate
	scenarios := Initialize()
	log.Info(scenarios)
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

		reader, _ := os.Open(filepath.Join(v.SourcePath,v.archive))

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
			output := executeCommand(ctx, "bbi", "nonmem","run","local", "--nmVersion",nmVersion,m.filename)
			assert.Contains(t,output,"Beginning local work")
			assert.Contains(t,output,"Beginning cleanup")

			modelName := m.filename
			outputDir := filepath.Join(EXECUTION_DIR,modelSet,modelName)

			xmlControlStream, err := afero.Exists(fs,filepath.Join(outputDir,modelName + ".xml"))

			assert.Nil(t,err)
			assert.True(t,xmlControlStream)

			nmlines, err := fileLines(filepath.Join(outputDir,modelName + ".lst"))

			assert.Nil(t,err)
			assert.NotNil(t,nmlines)
			assert.NotEmpty(t,nmlines)
			//Make sure that nonmem shows it finished and generated files
			assert.Contains(t,strings.Join(nmlines,"\n"),"finaloutput")
			//Make sure that nonmem records a stop time
			assert.Contains(t,strings.Join(nmlines,"\n"),"Stop Time:")


			expected := []string{
				".xml",
				".cpu",
				".grd",
			}

			for _, v := range expected {
					ok, _ := afero.Exists(fs,filepath.Join(outputDir,modelName + v))
					assert.True(t,ok)
			}
		}

	}
}

