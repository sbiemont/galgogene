package engine

import (
	"golang.org/x/sync/errgroup"
)

func runParallelBatch(totalSize, nbBatches int, fct func(from, to, i int) error) error {
	grp := errgroup.Group{}
	batchSize := int((float64(totalSize) / float64(nbBatches)) + 0.5)

	for from, i := 0, 0; from < totalSize; from, i = from+batchSize, i+1 {
		to := min(from+batchSize, totalSize)
		grp.Go(func() error {
			return fct(from, to, i)
		})
	}

	return grp.Wait()
}
