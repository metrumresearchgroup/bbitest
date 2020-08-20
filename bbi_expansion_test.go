package babylontest

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//Test that expansion works with 001-005 etc.
func TestBBIExpandsWithoutPrefix(t *testing.T){
//	lines :=`DEBU[0000] expanded models: [240/001.mod 240/002.mod 240/003.mod 240/004.mod 240/005.mod 240/006.mod 240/007.mod 240/008.mod 240/009.mod]
//INFO[0000] A total of 9 models have completed the initial preparation phase`

	Scenario := InitializeScenarios([]string {
		"bbi_expansion",
	})[0]

	Scenario.Prepare(context.Background())

	targets := `10[1:5].ctl`

	commandAndArgs := []string{
		"-d", //Needs to be in debug mode to generate the expected output
		"--threads",
		"2",
		"nonmem",
		"run",
		"local",
		"--nm_version",
		os.Getenv("NMVERSION"),
		filepath.Join(Scenario.Workpath,"model",targets),
	}

	output, err := executeCommand(context.Background(),"bbi",commandAndArgs...)

	assert.Nil(t,err)
	assert.NotEmpty(t,output)

	modelsLine, _ := findOutputLine(strings.Split(output,"\n"))
	modelsLine = strings.TrimSuffix(modelsLine,"\n")
	expandedModels := outputLineToModels(modelsLine)

	//Verify that we expanded to five models
	assert.Len(t,expandedModels,5)

	//Verify nonmem completed for all five
	for _, m := range expandedModels{
		file := filepath.Base(m)
		extension := filepath.Ext(file)
		identifier := strings.Replace(file,extension,"",1)
		outputDir := filepath.Join(Scenario.Workpath,"model",identifier)

		internalModel := Model{
			identifier:identifier,
			filename: file,
			extension:extension,
			path:outputDir,
		}

		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: internalModel.path,
			Model:     internalModel,
			Output:    output,
			Scenario: Scenario,
		}

		AssertNonMemCompleted(nmd)
		AssertNonMemCreatedOutputFiles(nmd)
	}
}

//Test that expansion works with 001-005 etc.
func TestBBIExpandsWithPrefix(t *testing.T){
	//	lines :=`DEBU[0000] expanded models: [240/001.mod 240/002.mod 240/003.mod 240/004.mod 240/005.mod 240/006.mod 240/007.mod 240/008.mod 240/009.mod]
	//INFO[0000] A total of 9 models have completed the initial preparation phase`

	Scenario := InitializeScenarios([]string {
		"bbi_expansion",
	})[0]

	Scenario.Prepare(context.Background())

	targets := `bbi_mainrun_10[1:3].ctl`

	commandAndArgs := []string{
		"-d", //Needs to be in debug mode to generate the expected output
		"--threads",
		"2",
		"nonmem",
		"run",
		"local",
		"--nm_version",
		os.Getenv("NMVERSION"),
		filepath.Join(Scenario.Workpath,"model",targets),
	}

	output, err := executeCommand(context.Background(),"bbi",commandAndArgs...)

	assert.Nil(t,err)
	assert.NotEmpty(t,output)

	modelsLine, _ := findOutputLine(strings.Split(output,"\n"))
	modelsLine = strings.TrimSuffix(modelsLine,"\n")
	expandedModels := outputLineToModels(modelsLine)

	//Verify that we expanded to three models
	assert.Len(t,expandedModels,3)

	//Verify nonmem completed for all five
	for _, m := range expandedModels{
		file := filepath.Base(m)
		extension := filepath.Ext(file)
		identifier := strings.Replace(file,extension,"",1)
		outputDir := filepath.Join(Scenario.Workpath,"model",identifier)

		internalModel := Model{
			identifier:identifier,
			filename: file,
			extension:extension,
			path:outputDir,
		}

		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: internalModel.path,
			Model:     internalModel,
			Output:    output,
		}

		AssertNonMemCompleted(nmd)
		AssertNonMemCreatedOutputFiles(nmd)
	}
}

//Test that expansion works with 001-005 etc.
func TestBBIExpandsWithPrefixToPartialMatch(t *testing.T){
	//	lines :=`DEBU[0000] expanded models: [240/001.mod 240/002.mod 240/003.mod 240/004.mod 240/005.mod 240/006.mod 240/007.mod 240/008.mod 240/009.mod]
	//INFO[0000] A total of 9 models have completed the initial preparation phase`

	Scenario := InitializeScenarios([]string {
		"bbi_expansion",
	})[0]

	Scenario.Prepare(context.Background())

	targets := `bbi_mainrun_10[2:3].ctl`

	commandAndArgs := []string{
		"nonmem",
		"run",
		"local",
		filepath.Join(Scenario.Workpath,"model",targets),
	}

	output, err := executeCommand(context.Background(),"bbi",commandAndArgs...)

	assert.Nil(t,err)
	assert.NotEmpty(t,output)

	modelsLine, _ := findOutputLine(strings.Split(output,"\n"))
	modelsLine = strings.TrimSuffix(modelsLine,"\n")
	expandedModels := outputLineToModels(modelsLine)

	//Verify that we expanded to three models
	assert.Len(t,expandedModels,2)

	//Verify nonmem completed for all five
	for _, m := range expandedModels{
		file := filepath.Base(m)
		extension := filepath.Ext(file)
		identifier := strings.Replace(file,extension,"",1)
		outputDir := filepath.Join(Scenario.Workpath,"model",identifier)

		internalModel := Model{
			identifier:identifier,
			filename: file,
			extension:extension,
			path:outputDir,
		}

		nmd := NonMemTestingDetails{
			t:         t,
			OutputDir: internalModel.path,
			Model:     internalModel,
			Output:    output,
		}

		AssertNonMemCompleted(nmd)
		AssertNonMemCreatedOutputFiles(nmd)
	}
}


func outputLineToModels(expansionLine string) []string{
	slfields := strings.Split(expansionLine,":")
	arrayComponent := slfields[len(slfields) - 1]
	arrayComponent = strings.TrimSpace(arrayComponent)
	primaryPieces := strings.Split(arrayComponent,"[")[1]
	secondaryPieces := strings.Split(primaryPieces,"]")[0]
	models := strings.Fields(secondaryPieces)
	return  models
}

func findOutputLine(outputLines []string) (string, error) {
	for _, v := range outputLines {
		if strings.Contains(v,"expanded models:"){
			return v, nil
		}
	}

	return "", errors.New("No matching line of text could be found in the provided output")
}
