package engine

import (
	"reflect"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

// Solution after running the engine
// * best population (with elite individual)
// * best population (with max total fitness)
// * termination operator triggered
type Solution struct {
	PopWithBestIndividual   gene.Population      // Population with best computed individual
	PopWithBestTotalFitness gene.Population      // Population with best total fitness computed
	Termination             operator.Termination // Termination that triggered the end of computation
}

func (sol Solution) TerminationType(term operator.Termination) bool {
	return reflect.TypeOf(sol.Termination) == reflect.TypeOf(term)
}
