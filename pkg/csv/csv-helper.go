package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvData interface {
	AsCsvRecord() []string
}

func SaveIndex[T CsvData](filePath string, data []T) error {
	file, err := os.Create(filePath + ".new")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	// Write data
	for _, datum := range data {
		err := writer.Write(datum.AsCsvRecord())
		if err != nil {
			fmt.Println("Error writing CSV row:", err)
			return err
		}
	}

	writer.Flush()

	// Check for errors during Flush
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing CSV writer:", err)
		return err
	}

	os.Rename(filePath+".new", filePath)

	return nil
}

type ConverterFunc[T CsvData] func([]string) (T, error)

func LoadIndex[T CsvData](filePath string, converter ConverterFunc[T]) []T {

	rawCsv, error := loadCsv(filePath)

	if error != nil {
		fmt.Println("Could not load data:", error.Error())
		os.Exit(1)
	}

	var items []T

	for _, line := range rawCsv {
		data, error := converter(line)
		if error == nil {
			items = append(items, data)
		}
	}

	return items
}

func loadCsv(filePath string) ([][]string, error) {
	file, fileOpenError := os.Open(filePath)
	if fileOpenError != nil {
		return nil, fileOpenError
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	data, csvReadError := csvReader.ReadAll()
	if csvReadError != nil {
		return nil, csvReadError
	}

	return data, nil
}
