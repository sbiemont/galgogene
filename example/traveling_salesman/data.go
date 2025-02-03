package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type DataCsv struct {
	Filename string
}

// ReadCoordinates from the given csv file
func (dc DataCsv) ReadCoordinates() ([][2]float64, error) {
	// Open file
	f, err := os.Open(dc.Filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Convert [][]string to [][]float64
	var xMin, xMax, yMin, yMax float64
	result := make([][2]float64, len(data))
	for i, row := range data {
		result[i] = [2]float64{}
		for j, col := range row {
			value, err := strconv.ParseFloat(col, 64)
			if err != nil {
				return nil, err
			}
			result[i][j] = value
		}
		xMin = min(xMin, result[i][0])
		xMax = max(xMax, result[i][0])
		yMin = min(yMin, result[i][1])
		yMax = max(yMax, result[i][1])
	}

	// Make it bigger
	for i := range data {
		// result[i][0] -= xMin
		// result[i][1] -= yMin
		result[i][0] *= 2
		result[i][1] *= 2
	}
	return result, nil
}
