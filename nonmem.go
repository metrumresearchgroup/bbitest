package babylontest

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

type NonMemTestingDetails struct {
	t *testing.T
	OutputDir string
	Model Model
	Output string
}

func AssertNonMemCompleted(details NonMemTestingDetails){
	nmlines, err := fileLines(filepath.Join(details.OutputDir,details.Model.identifier + ".lst"))

	assert.Nil(details.t,err)
	assert.NotNil(details.t,nmlines)
	assert.NotEmpty(details.t,nmlines)
	//Make sure that nonmem shows it finished and generated files
	assert.Contains(details.t,strings.Join(nmlines,"\n"),"finaloutput")
	//Make sure that nonmem records a stop time
	assert.Contains(details.t,strings.Join(nmlines,"\n"),"Stop Time:")
}

func AssertNonMemCreatedOutputFiles( details NonMemTestingDetails){
	fs := afero.NewOsFs()
	expected := []string{
		".xml",
		".cpu",
		".grd",
	}

	for _, v := range expected {
		ok, _ := afero.Exists(fs,filepath.Join(details.OutputDir,details.Model.identifier + v))
		assert.True(details.t,ok,"Unable to locate expected file %s",v)
	}
}

func AssertBBIConfigJSONCreated( details NonMemTestingDetails){
	fs := afero.NewOsFs()

	ok, _ := afero.Exists(fs,filepath.Join(details.OutputDir,"bbi_config.json"))
	assert.True(details.t,ok)
}

func AssertContainsBBIScript( details NonMemTestingDetails){

	fs := afero.NewOsFs()

	ok, _ := afero.Exists(fs,filepath.Join(details.OutputDir,details.Model.identifier + ".sh"))
	assert.True(details.t,ok,"The required BBI execution script %s, is not present in the output dir", details.Model.identifier+".sh")
}


func AssertNonMemOutputContainsParafile( details NonMemTestingDetails){
	containsParafile := false

	lines, _ :=  fileLines(filepath.Join(details.OutputDir,details.Model.identifier + ".lst"))

	for _, v := range lines {
		if strings.Contains(v,"PARAFILE="){
			containsParafile = true
		}
	}

	assert.True(details.t,containsParafile)
}

func AssertDefaultConfigLoaded (details NonMemTestingDetails){
	assert.True(details.t,strings.Contains(details.Output,"Successfully loaded default configuration"))
}

func AssertSpecifiedConfigLoaded(details NonMemTestingDetails, specificFile string){
	message := fmt.Sprintf("Successfully loaded specified configuration from %s",specificFile)
	assert.True(details.t, strings.Contains(details.Output,message))
}

func AssertContainsNMFEOptions(details NonMemTestingDetails, filepath string,  optionValue string) {
	content, _ := ioutil.ReadFile(filepath)
	contentString := string(content)
	assert.True(details.t,strings.Contains(contentString,optionValue))
}