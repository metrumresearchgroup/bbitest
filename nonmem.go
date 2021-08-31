package bbitest

import (
	"crypto/md5"
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type NonMemTestingDetails struct {
	t *testing.T
	OutputDir string
	Model Model
	Output string
	Scenario *Scenario
}

func AssertNonMemCompleted(details NonMemTestingDetails){
	nmlines, err := fileLines(filepath.Join(details.OutputDir,details.Model.identifier + ".lst"))

	require.Nil(details.t,err)
	require.NotNil(details.t,nmlines)
	require.NotEmpty(details.t,nmlines)

	//Check for either finaloutput or Stop Time as Stop Time appears in older versions of nonmem

	require.True(details.t,strings.Contains(strings.Join(nmlines,"\n"),"finaloutput") || strings.Contains(strings.Join(nmlines,"\n"),"Stop Time:"))
}

func AssertNonMemCreatedOutputFiles( details NonMemTestingDetails){
	fs := afero.NewOsFs()
	expected := []string{
		".xml",
		".cpu",
		//".grd",
	}

	for _, v := range expected {
		ok, _ := afero.Exists(fs,filepath.Join(details.OutputDir,details.Model.identifier + v))
		require.True(details.t,ok,"Unable to locate expected file %s",v)
	}
}

func AssertBBIConfigJSONCreated( details NonMemTestingDetails){
	fs := afero.NewOsFs()

	ok, _ := afero.Exists(fs,filepath.Join(details.OutputDir,"bbi_config.json"))
	require.True(details.t,ok)
}

func AssertContainsBBIScript( details NonMemTestingDetails){

	fs := afero.NewOsFs()

	ok, _ := afero.Exists(fs,filepath.Join(details.OutputDir,details.Model.identifier + ".sh"))
	require.True(details.t,ok,"The required BBI execution script %s, is not present in the output dir", details.Model.identifier+".sh")
}


func AssertNonMemOutputContainsParafile( details NonMemTestingDetails){
	containsParafile := false

	lines, _ :=  fileLines(filepath.Join(details.OutputDir,details.Model.identifier + ".lst"))

	for _, v := range lines {
		if strings.Contains(v,"PARAFILE="){
			containsParafile = true
		}
	}

	require.True(details.t,containsParafile)
}

func AssertDefaultConfigLoaded (details NonMemTestingDetails){
	require.True(details.t,strings.Contains(details.Output,"Successfully loaded default configuration"))
}

func AssertSpecifiedConfigLoaded(details NonMemTestingDetails, specificFile string){
	message := fmt.Sprintf("Successfully loaded specified configuration from %s",specificFile)
	require.True(details.t, strings.Contains(details.Output,message))
}

func AssertContainsNMFEOptions(details NonMemTestingDetails, filepath string,  optionValue string) {
	content, _ := ioutil.ReadFile(filepath)
	contentString := string(content)
	require.True(details.t,strings.Contains(contentString,optionValue))
}

//Make sure that the BBIConfig json has a value for the data hash
func AssertDataSourceIsHashedAndCorrect(details NonMemTestingDetails) {
	file, _ := os.Open(filepath.Join(details.Scenario.Workpath,"scenario.json"))
	originalDetails := GetScenarioDetailsFromFile(file)

	bbiConfigJson, _ := os.Open(filepath.Join(details.OutputDir,"bbi_config.json"))
	savedHashes := GetBBIConfigJSONHashedValues(bbiConfigJson)

	require.NotEmpty(details.t,savedHashes.Data)

	//Get MD5 of current file
	datafile, _ := os.Open(filepath.Join(details.Scenario.Workpath,originalDetails.DataFile))
	defer datafile.Close()
	dataHash, _ := calculateMD5(datafile)

	//Make sure the calculated and saved values are the same
	require.Equal(details.t, savedHashes.Data, dataHash)
}

func AssertModelIsHashedAndCorrect(details NonMemTestingDetails){
	bbiConfigJson, _ := os.Open(filepath.Join(details.OutputDir,"bbi_config.json"))
	savedHashes := GetBBIConfigJSONHashedValues(bbiConfigJson)

	require.NotEmpty(details.t,savedHashes.Model)

	//Get MD5 of model ORIGINAL file
	//Getting a hash of the model's relative copy is pointless since we may or may not have modified it's $DATA location
	model, _ := os.Open(filepath.Join(details.Scenario.Workpath,details.Model.filename))
	defer model.Close()

	calculatedHash, _ := calculateMD5(model)

	require.Equal(details.t,savedHashes.Model,calculatedHash)
}


func calculateMD5(r io.Reader) (string, error){
	h := md5.New()

	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x",h.Sum(nil)),nil
}