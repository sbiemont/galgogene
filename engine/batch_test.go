package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBatches(t *testing.T) {
	Convey("batches inf", t, func() {
		froms := make([]int, 2)
		tos := make([]int, 2)
		err := runParallelBatch(15, 2, func(from, to, i int) error {
			froms[i] = from
			tos[i] = to
			return nil
		})
		So(err, ShouldBeNil)
		So(froms, ShouldResemble, []int{0, 8})
		So(tos, ShouldResemble, []int{8, 15})
	})

	Convey("batches fit", t, func() {
		froms := make([]int, 3)
		tos := make([]int, 3)
		err := runParallelBatch(15, 3, func(from, to, i int) error {
			froms[i] = from
			tos[i] = to
			return nil
		})
		So(err, ShouldBeNil)
		So(froms, ShouldResemble, []int{0, 5, 10})
		So(tos, ShouldResemble, []int{5, 10, 15})
	})

	Convey("batches sup", t, func() {
		froms := make([]int, 4)
		tos := make([]int, 4)
		err := runParallelBatch(15, 4, func(from, to, i int) error {
			froms[i] = from
			tos[i] = to
			return nil
		})
		So(err, ShouldBeNil)
		So(froms, ShouldResemble, []int{0, 4, 8, 12})
		So(tos, ShouldResemble, []int{4,8,12, 15})
	})
}
