package utils

import (
	"encoding/csv"
	"os"

	"github.com/dimchansky/utfbom"
)

func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	sr, _ := utfbom.Skip(f)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(sr)
	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return csvData, nil
}
