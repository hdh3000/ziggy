package importer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

func WriteSqlMetaData(storeLoc string, meta *SQLImportMetaData) error {
	rF, err := os.Open(storeLoc)
	if err != nil {
		return err
	}
	defer rF.Close()

	var data []SQLImportMetaData
	if err := json.NewDecoder(rF).Decode(&data); err != nil {
		return err
	}

	data = append(data, *meta)

	sort.Slice(data, func(i, j int) bool {
		return data[i].TableName < data[j].TableName
	})

	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	rF.Close()

	// Write to a backup...
	bkFilePath := fmt.Sprintf("%s.%s", storeLoc, "bk")
	bkW, err := os.Create(bkFilePath)
	if _, err := bkW.Write(b); err != nil {
		return err
	}

	wF, err := os.Create(storeLoc)
	if err != nil {
		return err
	}
	defer wF.Close()

	if _, err := wF.Write(b); err != nil {
		return err
	}

	os.Remove(bkFilePath)

	return nil
}
