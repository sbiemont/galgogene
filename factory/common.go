package factory

import (
	"time"

	"github.com/sbiemont/galgogene/operator"
)

// Selection

type commonSelection struct{}

func (f commonSelection) Roulette() operator.RouletteSelection {
	return operator.RouletteSelection{}
}

func (f commonSelection) Tournament(fighters int) operator.TournamentSelection {
	return operator.TournamentSelection{
		Fighters: fighters,
	}
}

func (f commonSelection) Elite() operator.EliteSelection {
	return operator.EliteSelection{}
}

func (f commonSelection) Multi() operator.MultiSelection {
	return operator.MultiSelection{}
}

// Survivor

type commonSurvivor struct{}

func (f commonSurvivor) Elite() operator.EliteSurvivor {
	return operator.EliteSurvivor{}
}

func (f commonSurvivor) Rank() operator.RankSurvivor {
	return operator.RankSurvivor{}
}

func (f commonSurvivor) Random() operator.RandomSurvivor {
	return operator.RandomSurvivor{}
}

func (f commonSurvivor) Multi() operator.MultiSurvivor {
	return operator.MultiSurvivor{}
}

// Termination

type commonTermination struct{}

func (f commonTermination) Generation(k int) *operator.GenerationTermination {
	return &operator.GenerationTermination{K: k}
}

func (f commonTermination) Improvement(k int) *operator.ImprovementTermination {
	return &operator.ImprovementTermination{K: k}
}

func (f commonTermination) Fitness(fitness float64) *operator.FitnessTermination {
	return &operator.FitnessTermination{Fitness: fitness}
}

func (f commonTermination) Duration(duration time.Duration) *operator.DurationTermination {
	return &operator.DurationTermination{Duration: duration}
}

func (f commonTermination) Multi() operator.MultiTermination {
	return operator.MultiTermination{}
}
