package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"time"
)

func ExportShp(query, file string) (string, error) {
	dir := path.Dir(file)
	zipPath := fmt.Sprintf("%s.zip", dir)
	// Make directory
	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		return "", err
	}

	exportCmd := exec.Command("pgsql2shp",
		"-f", file,
		"-h", "127.0.0.1",
		"-p", "5432",
		"-P", "postgres",
		"-u", "postgres",
		"postgres",
		query,
	)

	if err := exportCmd.Run(); err != nil {
		return "", err
	}

	zipCmd := exec.Command("zip", "-jrm", zipPath, dir)

	if out, err := zipCmd.CombinedOutput(); err != nil {
		log.Println(string(out))
		return "", err
	}

	if err := os.RemoveAll(dir); err != nil {
		return "", err
	}

	return zipPath, nil

}

type SQLExportMetaData struct {
	Date  time.Time
	Query string
	Name  string
}

func WriteExportMetadata(storeLoc string, meta *SQLExportMetaData) error {
	rF, err := os.Open(storeLoc)
	if err != nil {
		return err
	}
	defer rF.Close()

	var data []SQLExportMetaData
	if err := json.NewDecoder(rF).Decode(&data); err != nil {
		return err
	}

	data = append(data, *meta)

	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
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
