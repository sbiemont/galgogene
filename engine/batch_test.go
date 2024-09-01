package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBatches(t *testing.T) {
	Convey("batches", t, func() {
		froms := make([]int, 3)
		tos := make([]int, 3)
		runParallelBatch(15, 3, func(from, to, i int) error {
			froms[i] = from
			tos[i] = to
			return nil
		})
		So(froms, ShouldResemble, []int{0, 5, 10})
		So(tos, ShouldResemble, []int{5, 10, 15})
	})
}
