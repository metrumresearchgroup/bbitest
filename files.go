package babylontest

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func findNonMemKey(pathToBabylonConfig string) (string, error) {
	contentLines, err := fileLines(pathToBabylonConfig)

	if err != nil {
		log.Fatalf("There was an issue trying to read the contents of file %s. Error " +
			"details are %s",pathToBabylonConfig,err)
	}

	for k, l := range contentLines{
		if l == "nonmem:" {
			return strings.TrimSpace(strings.Split(contentLines[k+1],":")[0]),nil
		}
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
