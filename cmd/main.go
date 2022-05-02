package main

import (
	Cykle "Cykle_go"
	"encoding/csv"
	"fmt"
	"golang.org/x/exp/constraints"
	"log"
	"os"
	"time"
)

func main() {
	os.Remove("./wyniki.csv")

	var (
		n, n02, n06                    []int
		e02, e06, h02, h06, ha02, ha06 []float64
	)

	lastHa06 := float64(0)
	size := 5
	for size = size; lastHa06 < 5; size++ {
		n = append(n, size)

		g02 := Cykle.NewGraphWithDensity(size, 0.2)
		if !g02.CheckIfEulerian() {
			panic(fmt.Sprintf("Graph G02 is not Eulerian!!!"))
		}
		g06 := Cykle.NewGraphWithDensity(size, 0.6)
		if !g06.CheckIfEulerian() {
			panic(fmt.Sprintf("Graph G06 is not Eulerian!!!"))
		}

		nCycles02, timeFirst02, timeAll02 := g02.MeasureHamilton(true)
		h02 = append(h02, timeFirst02)
		ha02 = append(ha02, timeAll02)
		n02 = append(n02, nCycles02)

		nCycles06, timeFirst06, timeAll06 := g06.MeasureHamilton(true)
		h06 = append(h06, timeFirst06)
		ha06 = append(ha06, timeAll06)
		n06 = append(n06, nCycles06)

		lastHa06 = timeAll06

		log.Printf("Graph density: %f\n", g02.GetDensity())
		log.Printf("Graph density: %f\n", g06.GetDensity())
		e02 = append(e02, timeFunction(fmt.Sprintf("Eulerian circuit for n=%d, d=%.1f", size, 0.2), g02.GetEulerianCircuit))
		e06 = append(e06, timeFunction(fmt.Sprintf("Eulerian circuit for n=%d, d=%.1f", size, 0.6), g06.GetEulerianCircuit))
		log.Printf("First Hamilton circuit for n=%d, d=%.1f = %fs", size, 0.2, timeFirst02)
		log.Printf("All Hamilton circuit for n=%d, d=%.1f = %fs", size, 0.2, timeAll02)
		log.Printf("First Hamilton circuit for n=%d, d=%.1f = %fs", size, 0.6, timeFirst06)
		log.Printf("All Hamilton circuit for n=%d, d=%.1f = %fs", size, 0.6, timeAll06)
	}

	for i := 0; i < 10; i++ {
		size += 50

		n = append(n, size)

		log.Printf("Generating d=0.2 for n=%d", size)
		g02 := Cykle.NewGraphWithDensity(size, 0.2)
		if !g02.CheckIfEulerian() {
			panic(fmt.Sprintf("Graph G02 is not Eulerian!!!"))
		}
		log.Printf("Finished")

		log.Printf("Generating d=0.6 for n=%d", size)
		g06 := Cykle.NewGraphWithDensity(size, 0.6)
		if !g06.CheckIfEulerian() {
			panic(fmt.Sprintf("Graph G06 is not Eulerian!!!"))
		}
		log.Printf("Finished")

		_, timeFirst02, _ := g02.MeasureHamilton(false)
		h02 = append(h02, timeFirst02)

		_, timeFirst06, _ := g06.MeasureHamilton(false)
		h06 = append(h06, timeFirst06)

		log.Printf("Graph density: %f\n", g02.GetDensity())
		log.Printf("Graph density: %f\n", g06.GetDensity())
		e02 = append(e02, timeFunction(fmt.Sprintf("Eulerian circuit for n=%d, d=%.1f", size, 0.2), g02.GetEulerianCircuit))
		e06 = append(e06, timeFunction(fmt.Sprintf("Eulerian circuit for n=%d, d=%.1f", size, 0.6), g06.GetEulerianCircuit))
		log.Printf("First Hamilton circuit for n=%d, d=%.1f = %fs", size, 0.2, timeFirst02)
		log.Printf("First Hamilton circuit for n=%d, d=%.1f = %fs", size, 0.6, timeFirst06)
	}

	f, err := os.Create("./wyniki.csv")
	if err != nil {
		panic("Failed to open wyniki.csv")
	}
	defer f.Close()

	writer := csv.NewWriter(f)

	writer.Write(writeLineToCSV("number of elements", n))
	writer.Write(writeLineToCSV("Euler d=0.2", e02))
	writer.Write(writeLineToCSV("Euler d=0.6", e06))
	writer.Write(writeLineToCSV("Hamilton jeden d=0.2", h02))
	writer.Write(writeLineToCSV("Hamilton wszystkie d=0.2", ha02))
	writer.Write(writeLineToCSV("Hamilton liczba cykli d=0.2", n02))
	writer.Write(writeLineToCSV("Hamilton jeden d=0.6", h06))
	writer.Write(writeLineToCSV("Hamilton wszystkie d=0.2", ha06))
	writer.Write(writeLineToCSV("Hamilton liczba cykli d=0.6", n06))

	writer.Flush()
}

func writeLineToCSV[T constraints.Float | constraints.Integer](name string, line []T) []string {
	asStrings := make([]string, len(line)+1)

	asStrings[0] = name
	for i, val := range line {
		asStrings[i+1] = fmt.Sprintf("%v", val)
	}

	return asStrings
}

func timeFunction[T any](name string, f func() T) float64 {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	log.Printf("%s took %fs", name, elapsed.Seconds())
	return elapsed.Seconds()
}
