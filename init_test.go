package bbitest

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/babylon/configlib"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInitialization(t *testing.T){
	scenarios := InitializeScenarios([]string{
		"240",
		"acop",
		"ctl_test",
		"metrum_std",
	})


	for _, s := range scenarios {
		s.Prepare(context.Background())

		t.Run(fmt.Sprintf("init_%s",s.identifier),func(t *testing.T){
			_, err := executeCommand(context.Background(), "bbi", "init","--dir",os.Getenv("NONMEMROOT"))

			assert.Nil(t,err)

			assert.FileExists(t,filepath.Join(s.Workpath,"bbi.yaml"))

			//Verify that we have nonmem contents!
			c := configlib.Config{}

			configHandle, _ := os.Open(filepath.Join(s.Workpath,"bbi.yaml"))
			bytes, _  := ioutil.ReadAll(configHandle)
			yaml.Unmarshal(bytes,&c)

			assert.Greater(t,len(c.Nonmem),0)
			configHandle.Close()
		})
	}
}
