package importer

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// This is a global projection with minimal distortion.
const projectProjection = 4326

type SQLImportMetaData struct {
	TableName   string `fd:"key"`
	Date        time.Time
	Source      string
	SourceType  string
	DetectedPrj int
	ImportType  string
}

type SQLDataSource interface {
	ImpToSql() (*SQLImportMetaData, error) // Calls various shell commands to import file
}

func toTableName(in string) string {
	return strings.Replace(strings.ToLower(in), "-", "_", -1)
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

func pipeToPsql(cmd *exec.Cmd) error {
	// very secure....
	// auth here is done by a cloudsqlproxy.
	psqlCmd := exec.Command(
		"psql",
		"host=127.0.0.1 sslmode=disable dbname=postgres user=postgres password=postgres",
	)

	// Pipe the commands together
	r, w := io.Pipe()
	defer w.Close()
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	psqlCmd.Stdin = r
	psqlCmd.Stdout = os.Stdout
	psqlCmd.Stderr = os.Stderr

	cmd.Start()
	psqlCmd.Start()
	if err := cmd.Wait(); err != nil {
		return err
	}
	w.Close()

	if err := psqlCmd.Wait(); err != nil {
		return err
	}

	return nil
}
