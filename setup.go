package babylontest

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const EXECUTION_DIR string = "/tmp/working"
const BBI_VERSION string = "2.1.0-alpha.6"
const BINDIR string = "/usr/local/bin"
var BBI_RELEASE string = fmt.Sprintf("https://github.com/metrumresearchgroup/babylon/releases/download/v%s/bbi_%s_%s_amd64.tar.gz",BBI_VERSION,BBI_VERSION,runtime.GOOS)

type Scenario struct {
	identifier string
	SourcePath string
	Workpath string
	models []Model
	archive string //The name of the tar.gz file used
}

type Model struct {
	identifier string //acop or Executive_Mod
	filename string //acop.mod or Executive_Mod.mod
	extension string//.mod or .ctl
	path string //Path at which model resides.
}

func newScenario(path string) (Scenario, error) {

	scenario := Scenario{
		identifier: filepath.Base(path),
		models:     []Model{},
	}

	scenario.SourcePath = path
	scenario.Workpath = filepath.Join(EXECUTION_DIR,scenario.identifier)

	scenario.models = modelsFromOriginalScenarioPath(path)
	scenario.archive = scenario.identifier + ".tar.gz"

	if len(scenario.models) == 0{
		return scenario, errors.New("no model directories were located in the provided scenario")
	}

	return scenario, nil
}

func modelsFromOriginalScenarioPath(path string) []Model {

	models := []Model{}

	scenarioID := filepath.Base(path)
	newBaseDir := filepath.Join(EXECUTION_DIR,scenarioID)

	modelIdentifiers := []string{
		".ctl",
		".mod",
	}

	fs := afero.NewOsFs()

	for _, v := range modelIdentifiers {
		contents, _ := afero.Glob(fs,filepath.Join(path,"*" + v))
		for _, c := range contents {
			model := Model{
				filename:   filepath.Base(c),
			}

			modelPieces := strings.Split(filepath.Base(c),".")
			model.extension = filepath.Ext(filepath.Base(c))
			model.identifier = modelPieces[0]
			modelDir := filepath.Join(newBaseDir,model.identifier)
			model.path = modelDir

			models = append(models,model)

		}
	}

	return models
}

func Initialize()[]Scenario{
	downloadAndInstallBBI()

	viper.SetEnvPrefix("babylon")
	viper.AutomaticEnv()

	viper.SetDefault("nonmemroot","/opt/NONMEM")

	fs := afero.NewOsFs()
	if ok, _ := afero.DirExists(fs,EXECUTION_DIR); !ok {
		fs.MkdirAll(EXECUTION_DIR,0755)
	} else {
		fs.RemoveAll(EXECUTION_DIR)
		fs.MkdirAll(EXECUTION_DIR,0755)
	}

	scenarios := []Scenario{}



	dirs, _ := getScenarioDirs()
	whereami, _ := os.Getwd()

	//Let's navigate to each and try to tar it up
	//We'll use these later for execution layers by always starting with a clean slate from the tar content
	for _, v := range dirs {
		scenario, _ := newScenario(v)
		scenarios = append(scenarios,scenario)
		f, _ := os.Create(filepath.Join(whereami,"testdata", filepath.Base(v) + ".tar.gz"))
		err := Tar(filepath.Join(v),f)
		if err != nil {
			log.Error(err)
		}
		f.Close()
	}

	//Now let's find all the tar gz files and move them to the EXECUTIONDIR
	tars, _ := afero.Glob(afero.NewOsFs(),filepath.Join(whereami,"testdata","*.tar.gz"))

	for _, v := range tars {
		source, _ := os.Open(v)
		defer source.Close()

		dest, _ := os.Create(filepath.Join(EXECUTION_DIR,filepath.Base(v)))
		defer dest.Close()

		io.Copy(dest,source)
	}

	return scenarios
}

func getScenarioDirs() ([]string,error) {
	whereami, _ := os.Getwd()
	fs := afero.NewOsFs()
	directories := []string{}

	contents, err := afero.ReadDir(fs,filepath.Join(whereami,"testdata"))

	if err != nil {
		log.Error("Unable to parse directory contents of 'testdata'")
		return  directories,err
	}

	for _, v := range contents {
		if ok, _ := afero.IsDir(fs,filepath.Join(whereami,"testdata",v.Name())); ok{
			directories = append(directories,filepath.Join(whereami,"testdata",v.Name()))
		}
	}

	return directories, nil
}

func downloadAndInstallBBI(){


	fs := afero.NewOsFs()

	//Do nothing if the file already exists
	if ok, _ := afero.Exists(fs,"/usr/local/bin/bbi"); ok{
		return
	}

	downloadFile("/tmp/bbi.tar.gz",BBI_RELEASE)
	file, _ := os.Open("/tmp/bbi.tar.gz")
	Untar("/tmp", file)

	bbi, _ := os.Open("/tmp/bbi")
	defer bbi.Close()

	installed, _ := os.Create("/usr/local/bin/bbi")
	installed.Chmod(0755)
	defer installed.Close()

	io.Copy(installed,bbi)
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}

// Tar takes a source and variable writers and walks 'source' writing each file
// found to the tar writer; the purpose for accepting multiple writers is to allow
// for multiple outputs (for example a file, or md5 hash)
func Tar(src string, writers ...io.Writer) error {

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}


func findModelFiles(path string) []string {

	knownModelTypes := []string{
		".mod",
		".ctl",
	}

	foundModels := []string{}

	fs := afero.NewOsFs()

	for _, v := range knownModelTypes{
		located, _ := afero.Glob(fs,filepath.Join(path,"*" + v))
		for _ , l := range located{
			foundModels = append(foundModels,l)
		}
	}

	return foundModels

}