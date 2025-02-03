package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
)

const (
	center   = 150.0
	radius   = 100.0
	minCoord = 10.0
	maxCoord = 290.0
)

func toStrCoordinates(x, y float64) []string {
	return []string{
		fmt.Sprintf("%.2f", x),
		fmt.Sprintf("%.2f", y),
	}
}

func writeCircle(nbPoints int, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	wr := csv.NewWriter(f)

	var records [][]string
	dt := math.Pi / float64(nbPoints)
	for theta := 0.0; theta < 2*math.Pi; theta += dt {
		x := radius*math.Cos(theta) + center
		y := radius*math.Sin(theta) + center
		records = append(records, toStrCoordinates(x, y))
	}
	return wr.WriteAll(records)
}

func writeRandom(nbPoints int, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	wr := csv.NewWriter(f)

	var records [][]string
	for range nbPoints {
		x := minCoord + rand.Float64()*(maxCoord-minCoord)
		y := minCoord + rand.Float64()*(maxCoord-minCoord)
		records = append(records, toStrCoordinates(x, y))
	}
	return wr.WriteAll(records)
}
