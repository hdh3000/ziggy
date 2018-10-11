package exporter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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
