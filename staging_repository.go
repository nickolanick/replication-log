package main

import (
)

// TODO: consider pointer
type StagingRepository []WriteConsistencyMessage

func (h StagingRepository) Len() int { return len(h) }
func (h StagingRepository) Less(i, j int) bool { return h[i].TotalOrder < h[j].TotalOrder }
func (h StagingRepository) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *StagingRepository) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(WriteConsistencyMessage))
}

func (h *StagingRepository) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
