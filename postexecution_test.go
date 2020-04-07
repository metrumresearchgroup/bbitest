package babylontest

import (
	"bufio"
	"context"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//Will need to set a custom env for execution
const postExecutionScriptString string = `#!/bin/bash

env > /tmp/working/${SCENARIO}/${BABYLON_MODEL_FILENAME}.out
`

func TestPostExecutionSucceeds(t *testing.T){

	//Test if feature flag enabled
	if FeatureEnabled("POST_EXECUTION"){
		Scenarios := InitializeScenarios([]string{
			"240",
			"acop",
			"ctl_test",
			"metrum_std",
		})

		ioutil.WriteFile(filepath.Join("/tmp","post.sh"),[]byte(postExecutionScriptString),0755)

		for _, v := range Scenarios {
			v.Prepare(context.Background())

			arguments := []string{
				"-d",
				"nonmem",
				"--nm_version",
				os.Getenv("NMVERSION"),
				"run",
				"local",
				"--post_work_executable",
				"/tmp/post.sh",
				"--additional_post_work_envs=\"SCENARIO=" + v.identifier + "\"",
			}

			//Do the actual execution
			for _, m := range v.models {
				t.Run(v.identifier + "_post_execution",func(t *testing.T){
					output, err := m.Execute(v,arguments...)
					assert.Nil(t,err)

					nmd := NonMemTestingDetails{
						t:         t,
						OutputDir: filepath.Join(v.Workpath,m.identifier),
						Model:     m,
						Output:    output,
					}

					AssertNonMemCompleted(nmd)
					AssertNonMemCreatedOutputFiles(nmd)

					exists, err := afero.Exists(afero.NewOsFs(),filepath.Join("/tmp","working",v.identifier,m.identifier + ".out") )

					assert.Nil(t,err)
					assert.True(t,exists)

					//Does the file contain the expected Details:
					//SCENARIO (Additional provided value)
					file, _ := os.Open(filepath.Join("/tmp","working",v.identifier, m.identifier + ".out"))

					scanner := bufio.NewScanner(file)
					scanner.Split(bufio.ScanLines)

					var lines []string

					for scanner.Scan() {
						lines = append(lines,scanner.Text())
					}


					assert.True(t, doesOutputFileContainKeyWithValue(lines,"BABYLON_MODEL",m.filename))
					assert.True(t, doesOutputFileContainKeyWithValue(lines, "BABYLON_MODEL_FILENAME", m.identifier))
					assert.True(t, doesOutputFileContainKeyWithValue(lines, "BABYLON_MODEL_EXT", strings.Replace(m.extension,".","",1)))
					assert.True(t, doesOutputFileContainKeyWithValue(lines, "BABYLON_SUCCESSFUL", "true"))
					assert.True(t, doesOutputFileContainKeyWithValue(lines, "BABYLON_ERROR", ""))

				})

			}
		}
	}
}


func doesOutputFileContainKeyWithValue(lines []string, key string, value string) bool {

	for _, v := range lines {
		if strings.Contains(v,key+"=") {
			components := strings.Split(v,"=")
			return components[0] == key && components[1] == value
		}
	}

	return false
}




