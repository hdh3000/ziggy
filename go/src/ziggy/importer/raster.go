package importer

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"time"
)

type rasterFile struct {
	systemName string // What is the system going to refer to this as?
	urlSource  string // Where did this come from?
	filePath   string // What is the path of the file on system to work with
}

func NewRaster(systemName, url, path string) SQLDataSource {
	return &rasterFile{systemName: systemName, urlSource: url, filePath: path}
}

func (r *rasterFile) ImpToSql() (*SQLImportMetaData, error) {
	prj, err := r.getPrj()
	if err != nil {
		return nil, err
	}

	rastCmd := exec.Command("raster2pgsql",
		"-I", "-M",
		"-s", fmt.Sprintf("%d", prj),
		"-t", "auto",
		r.filePath, toTableName(r.systemName),
	)

	if err := pipeToPsql(rastCmd); err != nil {
		return nil, err
	}

	return &SQLImportMetaData{
		TableName:   toTableName(r.systemName),
		Date:        time.Now(),
		Source:      r.urlSource,
		SourceType:  path.Ext(r.filePath),
		DetectedPrj: prj,
		ImportType:  "AUTO",
	}, nil

}

func (r *rasterFile) getPrj() (int, error) {
	gdalCmd := exec.Command("gdalinfo", "-json", r.filePath)

	relevantInfo := struct {
		CoordinateSystem struct {
			Wkt string `json:"wkt"`
		} `json:"coordinateSystem"`
	}{}

	rDesc, err := gdalCmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(rDesc, &relevantInfo); err != nil {
		return 0, err
	}

	return getPrj(relevantInfo.CoordinateSystem.Wkt)
}
