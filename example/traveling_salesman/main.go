package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/sbiemont/galgogene/engine"
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

// coordinates of all cities
// coordinates[0] gives coordinates(x,y) of city #0
var coordinates = [][2]float64{
	{0.05, 0.38}, {0.06, 0.95},
	{0.10, 0.32}, {0.10, 0.39},
	{0.17, 0.52}, {0.18, 0.20},
	{0.32, 0.53}, {0.32, 0.04},
	{0.38, 0.85}, {0.39, 0.41},
	{0.44, 0.46}, {0.45, 0.23},
	{0.51, 1.00}, {0.60, 0.50},
	{0.70, 0.28}, {0.71, 0.48},
	{0.71, 0.81}, {0.80, 0.23},
	{0.83, 0.76}, {0.98, 0.85},
	{0.17, 0.21}, {0.59, 0.14},
	{0.13, 0.95}, {0.44, 0.53},
	{0.42, 0.68}, {0.77, 0.09},
	{0.91, 0.33}, {0.82, 0.27},
	{0.01, 0.87}, {0.19, 0.05},
}

// distances is a small cache for cities distances
var distances = make(map[string]float64)

func dist(cityA, cityB gene.B) float64 {
	key := fmt.Sprintf("%d-%d", cityA, cityB)
	if d, ok := distances[key]; ok {
		return d
	}

	coordA := coordinates[cityA]
	coordB := coordinates[cityB]
	if coordA[0] == coordB[0] && coordA[1] == coordB[1] {
		return 0
	}
	d := math.Sqrt(math.Pow(coordA[0]-coordB[0], 2) + math.Pow(coordA[1]-coordB[1], 2))
	distances[key] = d
	return d
}

type cities []gene.B

func newCities(chrm gene.Chromosome) cities {
	return cities(chrm.Raw)
}

func (cts cities) String() string {
	result := make([]string, len(cts))
	for i, c := range cts {
		result[i] = fmt.Sprintf("%d", c)
	}
	return "[" + strings.Join(result, ", ") + "]"
}

// Compute distances [A, B, C, D]
// dist = A->B + B->C + C->D + D->A
func (cts cities) Distance() float64 {
	var distance float64
	var cityA gene.B = cts[0]
	for _, cityB := range cts[1:] {
		distance += dist(cityA, cityB)
		cityA = cityB
	}
	return distance + dist(cityA, cts[0])
}

func (cts cities) Fitness() float64 {
	distance := cts.Distance()
	if distance == 0 {
		return 0
	}
	return 1.0 / distance
}

func main() {
	popSize := 600
	eng := engine.Engine{
		Initializer: gene.PermutationInitializer{},
		Selection: operator.MultiSelection{}.
			Use(0.01, operator.EliteSelection{}).
			Otherwise(operator.RouletteSelection{}),
		CrossOver: operator.MultiCrossOver{}.
			Use(0.1, operator.UniformOrderCrossOver{}).
			Use(1.0, operator.DavisOrderCrossOver{}),
		Mutation: operator.MultiMutation{}.
			Use(0.05, operator.InversionPermutation{}).
			Use(0.05, operator.SwapPermutation{}).
			Use(0.05, operator.ScramblePermutation{}),
		Survivor: operator.MultiSurvivor{}.
			Use(0.6, operator.EliteSurvivor{}).
			Otherwise(operator.RandomSurvivor{}),
		Termination: operator.MultiTermination{}.
			Use(&operator.GenerationTermination{K: 1000}).
			Use(&operator.ImprovementTermination{K: 10}).
			Use(&operator.DurationTermination{Duration: 5 * time.Second}),
		Fitness: func(chrm gene.Chromosome) float64 {
			return newCities(chrm).Fitness()
		},
		OnNewGeneration: func(pop gene.Population) {
			elite := newCities(pop.Elite().Code)
			fmt.Printf(
				"Generation #%-3d, dur: %.3fs, fit: %.4f, tot-fit: %.4f, uniq: %d/%d\n",
				pop.Stats.GenerationNb,
				pop.Stats.TotalDuration.Seconds(),
				elite.Distance(),
				pop.Stats.TotalFitness,
				len(pop.Unique()),
				popSize,
			)
		},
	}

	// Run the engine
	nbCities := len(coordinates)
	solution, err := eng.Run(popSize, popSize*2, nbCities)
	if err != nil {
		panic(err)
	}

	// Print solution (best individual & best gen)
	out := func(msg string, p gene.Population) {
		elite := newCities(p.Elite().Code)
		fmt.Printf(
			"\n%s, gen: #%d, dur: %s, fit: %f, tot-fit: %f,\nElite ID: %s\n%s\n",
			msg,
			p.Stats.GenerationNb,
			p.Stats.TotalDuration,
			elite.Distance(),
			p.Stats.TotalFitness,
			p.Elite().ID,
			elite.String(),
		)
	}
	out("Best individual", solution.PopWithBestIndividual)
	out("Best generation", solution.PopWithBestTotalFitness)
}
