package bbitest

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/metrumresearchgroup/babylon/configlib"
	"github.com/metrumresearchgroup/wrapt"
	"gopkg.in/yaml.v2"
)

func TestInitialization(tt *testing.T) {
	t := wrapt.WrapT(tt)

	scenarios := InitializeScenarios([]string{
		"240",
		"acop",
		"ctl_test",
		"metrum_std",
	})

	for _, s := range scenarios {
		s.Prepare(context.Background())

		t.Run(fmt.Sprintf("init_%s", s.identifier), func(t *wrapt.T) {
			_, err := executeCommand(context.Background(), "bbi", "init", "--dir", os.Getenv("NONMEMROOT"))

			t.A.Nil(err)

			t.A.FileExists(filepath.Join(s.Workpath, "bbi.yaml"))

			// Verify that we have nonmem contents!
			c := configlib.Config{}

			configHandle, _ := os.Open(filepath.Join(s.Workpath, "bbi.yaml"))
			bytes, _ := ioutil.ReadAll(configHandle)
			yaml.Unmarshal(bytes, &c)

			t.A.Greater(len(c.Nonmem), 0)
			configHandle.Close()
		})
	}
}
