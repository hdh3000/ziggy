package importer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type shpFile struct {
	systemName  string // What is the system going to refer to this as?
	urlSource   string // Where did this come from?
	shpFilePath string // What is the shpFile path to work with
}

func NewShpFile(systemName, url, file string) SQLDataSource {
	return &shpFile{urlSource: url, shpFilePath: file, systemName: systemName}
}

func (sf *shpFile) ImpToSql() (*SQLImportMetaData, error) {
	filePrj, err := sf.getPrj()
	if err != nil {
		return nil, err
	}

	// shp2pgsql is an opengeo utility that produces a set of sql commands for importing files
	shpCmd := exec.Command("shp2pgsql",
		// Create indices
		"-I",
		// Transform shp prj to project prj on import
		"-s", fmt.Sprintf("%d:%d", filePrj, projectProjection),
		sf.shpFilePath,
		toTableName(sf.systemName),
	)

	if err = pipeToPsql(shpCmd); err != nil {
		return nil, err
	}

	return &SQLImportMetaData{
		TableName:   toTableName(sf.systemName),
		Date:        time.Now(),
		Source:      sf.urlSource,
		SourceType:  ".shp",
		DetectedPrj: filePrj,
		ImportType:  "AUTO",
	}, nil
}

func (sf *shpFile) getPrj() (int, error) {
	prjPath := strings.Replace(sf.shpFilePath, ".shp", ".prj", -1)
	f, err := os.Open(prjPath)
	if err != nil {
		return 0, fmt.Errorf("couldn't open a .prj in %s's dir...\n%s", sf.shpFilePath, err)
	}
	defer f.Close()

	prjString, _ := ioutil.ReadAll(f)
	return getPrj(string(prjString))

}
