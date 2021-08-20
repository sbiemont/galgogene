package example

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"genalgo.git/engine"
	"genalgo.git/gene"
	"genalgo.git/operator"
	. "github.com/smartystreets/goconvey/convey"
)

var coordinates = [][2]float64{
	{0.05, 0.38},
	{0.06, 0.95},
	{0.1, 0.32},
	{0.1, 0.39},
	{0.17, 0.52},
	{0.18, 0.2},
	{0.32, 0.53},
	{0.32, 0.04},
	{0.38, 0.85},
	{0.39, 0.41},
	{0.44, 0.46},
	{0.45, 0.23},
	{0.51, 1},
	{0.6, 0.5},
	{0.7, 0.28},
	{0.71, 0.48},
	{0.71, 0.81},
	{0.8, 0.23},
	{0.83, 0.76},
	{0.98, 0.85},
}

func dist(cityA, cityB uint8) float64 {
	coordA := coordinates[cityA]
	coordB := coordinates[cityB]
	if coordA[0] == coordB[0] && coordA[1] == coordB[1] {
		return 50000
	}
	return math.Sqrt(math.Pow(coordA[0]-coordB[0], 2) + math.Pow(coordA[1]-coordB[1], 2))
}

type cities []uint8

func toCities(bits gene.Bits) cities {
	return cities(bits.Raw)
}

func (cts cities) String() string {
	result := make([]string, len(cts))
	for i, c := range cts {
		result[i] = fmt.Sprintf("%d", c)
	}
	return "[" + strings.Join(result, ", ") + "]"
}

func (cts cities) Distance() float64 {
	var distance float64
	var cityA uint8 = cts[0]
	for _, cityB := range cts[1:] {
		distance += dist(cityA, cityB)
		cityA = cityB
	}
	return distance
}

func (cts cities) Fitness() float64 {
	distance := cts.Distance()
	if distance == 0 {
		return 0
	}
	return 1.0 / cts.Distance()
}

func TestTravelingSalesmanProblem(t *testing.T) {
	Convey("cities", t, func() {
		seed := time.Now().Unix()
		// var seed int64 = 1626634789
		rand.Seed(seed)

		popSize := 300
		eng := engine.Engine{
			Initializer: gene.PermuationInitializer{},
			Selection: operator.MultiSelection{}.
				Use(0.3, operator.TournamentSelection{Fighters: 2}).
				Otherwise(operator.RouletteSelection{}),
			CrossOver: operator.MultiCrossOver{}.
				Use(0.2, operator.UniformOrderCrossOver{}).
				Use(1, operator.DavisOrderCrossOver{}),
			Mutation: operator.MultiMutation{}.
				Use(0.15, operator.InversionPermutation{}),
			Survivor: operator.MultiSurvivor{}.
				Use(0.9, operator.SurvivorElite{}).
				Otherwise(operator.SurvivorChildren{}),
			Termination: operator.MultiTermination{}.
				Use(&operator.GenerationTermination{K: 1000}).
				Use(&operator.ImprovementTermination{K: 10}).
				Use(&operator.DurationTermination{Duration: 10 * time.Second}),
			OnNewGeneration: func(pop gene.Population) {
				elite := toCities(pop.Elite().Code)
				fmt.Printf(
					"Generation #%d, dur: %s dist: %f, tot: %f, cnt: %d/%d\n",
					pop.Stats.GenerationNb, pop.Stats.TotalDuration, elite.Distance(), pop.Stats.TotalFitness,
					pop.MapCount(), popSize,
				)
			},
		}

		last, _, err := eng.Run(popSize, 20, func(bits gene.Bits) float64 {
			cts := toCities(bits)
			return cts.Fitness()
		})

		if err != nil {
			panic(err)
		}

		elite := toCities(last.Elite().Code)
		fmt.Printf(
			"\nGeneration #%d, dur: %s dist: %f, tot: %f,\nseed: %d\n%s, dst: %f\n",
			last.Stats.GenerationNb,
			last.Stats.TotalDuration,
			elite.Distance(),
			last.Stats.TotalFitness,
			seed,
			elite.String(),
			elite.Distance(),
		)
	})
}
