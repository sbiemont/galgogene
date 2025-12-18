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

const (
	filenameCircle = "./example/traveling_salesman/circle.csv"
	filenameRandom = "./example/traveling_salesman/random.csv"
	filenameSquare = "./example/traveling_salesman/square.csv"
	maxDuration    = time.Minute
)

// coordinates of all cities
// coordinates[0] gives coordinates(x,y) of city #0
var coordinates [][2]float64

// distances is a small cache for cities distances
var distances map[string]float64

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

// Convert cities to 2D coordinates for printing
func (cts cities) Coordinates() [][2]float32 {
	result := make([][2]float32, len(cts))
	for i, city := range cts {
		coord := coordinates[city]
		result[i] = [2]float32{
			float32(coord[0]),
			float32(coord[1]),
		}
	}
	return result
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
	// err := writeCircle(30, filenameCircle)
	// if err != nil {
	// 	panic(err)
	// }

	// err = writeRandom(120, filenameRandom)
	// if err != nil {
	// 	panic(err)
	// }

	// Init game engine
	game, _ := NewGame(600, 600)
	go func() {
		err := Run(game)
		if err != nil {
			panic(err)
		}
	}()

	filenames := []string{filenameCircle, filenameRandom, filenameSquare}
	for _, filename := range filenames {
		run(game, DataCsv{
			Filename: filename,
		})
	}
}

func run(game *Game, dc DataCsv) {
	distances = make(map[string]float64)

	// Init coordinates
	var err error
	coordinates, err = dc.ReadCoordinates()
	if err != nil {
		panic(err)
	}

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
			Use(0.06, operator.InversionPermutation{}).
			Use(0.05, operator.SwapPermutation{}).
			Use(0.05, operator.ScramblePermutation{}),
		Survivor: operator.MultiSurvivor{}.
			Use(0.6, operator.EliteSurvivor{}).
			Otherwise(operator.RandomSurvivor{}),
		Termination: operator.MultiTermination{}.
			Use(&operator.GenerationTermination{K: 1500}).
			Use(&operator.ImprovementTermination{K: 2 * 100}).
			Use(&operator.DurationTermination{Duration: 2 * maxDuration}),
		Fitness: func(chrm gene.Chromosome) float64 {
			return newCities(chrm).Fitness()
		},
		OnNewGeneration: func(pop, withBestIndividual, _ gene.Population) {
			if pop.Stats.GenerationNb%10 == 0 {
				elite := newCities(withBestIndividual.Elite().Code)
				game.SetData(elite.Coordinates())
				game.Distance = elite.Distance()
			}
			game.GenerationNb = pop.Stats.GenerationNb

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

	fmt.Println("\nPress ENTER to exit")
	fmt.Scanln()
}
