package babylontest

import (
	"context"
	"encoding/json"
	"github.com/metrumresearchgroup/babylon/cmd"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBBIConfigJSONCreated(t *testing.T){
	scenarios := Initialize()

	for _, v := range scenarios{
		v.Prepare(context.Background())

		for _, m := range v.models {
			args := []string{
				"nonmem",
				"run",
				"local",
				"--nm_version",
				os.Getenv("NMVERSION"),
			}

			output, err := m.Execute(v,args...)

			assert.Nil(t,err)
			assert.NotNil(t,output)

			nmd := NonMemTestingDetails{
				t:         t,
				OutputDir: filepath.Join(v.Workpath,m.identifier),
				Model:     m,
				Output:    output,
			}

			AssertNonMemCompleted(nmd)
			AssertNonMemCreatedOutputFiles(nmd)
			AssertBBIConfigJSONCreated(nmd)
			AssertBBIConfigContainsSpecifiedNMVersion(nmd,os.Getenv("NMVERSION"))
		}
	}
}

func AssertBBIConfigContainsSpecifiedNMVersion(details NonMemTestingDetails, nmVersion string) {
	configFile, _ := os.Open(filepath.Join(details.OutputDir, "bbi_config.json"))
	cbytes, _ := ioutil.ReadAll(configFile)
	configFile.Close() //Go ahead and close the file handle

	nm := cmd.NonMemModel{}

	json.Unmarshal(cbytes, &nm)

	assert.NotNil(details.t,nm)
	assert.NotEqual(details.t,nm,cmd.NonMemModel{})
	assert.Equal(details.t,nm.Configuration.NMVersion,nmVersion)
}