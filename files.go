package babylontest

import (
	"errors"
	"github.com/metrumresearchgroup/babylon/configlib"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func findNonMemKey(pathToBabylonConfig string) (string, error) {
	fs := afero.NewOsFs()

	config := configlib.Config {}

	configFile, _ := fs.Open(pathToBabylonConfig)
	bytes, _ := afero.ReadAll(configFile)
	yaml.Unmarshal(bytes, &config)

	for k, _ := range config.Nonmem{
		return k, nil
	}

	return "", errors.New("unable to locate a key for a valid nonmem config")
}

func fileLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return []string{}, err
	}

	defer file.Close()

	contentBytes, err := ioutil.ReadAll(file)

	if err != nil {
		return []string{}, err
	}

	contentLines := strings.Split(string(contentBytes),"\n")

	return contentLines, nil
}
