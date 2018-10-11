package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type shpFile struct {
	systemName  string // What is the system going to refer to this as?
	urlSource   string // Where did this come from?
	shpFilePath string // What is the shpFile path to work with
}

func NewShpFile(systemName, url, file string) *shpFile {
	return &shpFile{urlSource: url, shpFilePath: file, systemName: systemName}
}

func (sf *shpFile) Path() string {
	return sf.shpFilePath
}

func (sf *shpFile) TableName() string {
	return strings.Replace(strings.ToLower(sf.systemName), "-", "_", -1)
}

func (sf *shpFile) Source() string {
	return sf.urlSource
}

func (sf *shpFile) ImpToSql() (*SQLImportMetaData, error) {
	filePrj, err := sf.GetPrj()
	if err != nil {
		return nil, err
	}

	// shp2pgsql is an opengeo utility that produces a set of sql commands for importing files
	shpCmd := exec.Command("shp2pgsql",
		// Create indices
		"-I",
		// Transform shp prj to project prj on import
		"-s", fmt.Sprintf("%d:%d", filePrj, projectProjection),
		sf.Path(),
		sf.TableName(),
	)

	// very secure....
	// auth here is done by a cloudsqlproxy.
	psqlCmd := exec.Command(
		"psql",
		"host=127.0.0.1 sslmode=disable dbname=postgres user=postgres password=postgres",
	)

	// Pipe the commands together
	r, w := io.Pipe()
	defer w.Close()
	shpCmd.Stdout = w
	shpCmd.Stderr = os.Stderr
	psqlCmd.Stdin = r
	psqlCmd.Stdout = os.Stdout // just so that it logs output as it tends to be long running.
	psqlCmd.Stderr = os.Stderr // just so that it logs output as it tends to be long running.

	shpCmd.Start()
	psqlCmd.Start()
	if err := shpCmd.Wait(); err != nil {
		return nil, err
	}
	w.Close()

	if err := psqlCmd.Wait(); err != nil {
		return nil, err
	}

	return &SQLImportMetaData{
		TableName:   sf.TableName(),
		Date:        time.Now(),
		Source:      sf.Source(),
		SourceType:  ".shp",
		DetectedPrj: filePrj,
		ImportType:  "AUTO",
	}, nil
}

func (sf *shpFile) GetPrj() (int, error) {
	prjPath := strings.Replace(sf.shpFilePath, ".shp", ".prj", -1)
	f, err := os.Open(prjPath)
	if err != nil {
		return 0, fmt.Errorf("couldn't open a .prj in %s's dir...\n%s", sf.shpFilePath, err)
	}
	defer f.Close()

	prjString, _ := ioutil.ReadAll(f)
	return getPrj(string(prjString))

}

func getPrj(prjDef string) (int, error) {
	// This could look for an exact match, but I don't care that much.
	type prjResp struct {
		Exact     bool   `json:"exact"`
		HTMLTerms string `json:"html_terms"`
		Codes     []struct {
			Name string `json:"name"`
			Code string `json:"code"`
			URL  string `json:"url"`
		} `json:"codes"`
		HTMLShowResults bool `json:"html_showResults"`
	}

	params := &url.Values{}
	params.Add("mode", "wkt")
	params.Add("terms", prjDef)

	target := url.URL{
		Scheme:   "http",
		Host:     "prj2epsg.org",
		Path:     "/search.json",
		RawQuery: params.Encode(),
	}

	resp, err := http.Get(target.String())
	if err != nil {
		return 0, err
	}

	var data prjResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if len(data.Codes) < 1 {
		return 0, errors.New("unable to find projection")
	}

	return strconv.Atoi(data.Codes[0].Code)
}
