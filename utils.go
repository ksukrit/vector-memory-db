package main

import (
	"math/rand"
)

// File contains custom sort functions for the KNN algorithm

type NeighborArr []NeighborVector

func (a NeighborArr) Len() int {
	return len(a)
}

func (a NeighborArr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a NeighborArr) Less(i, j int) bool {
	return a[i].Distance < a[j].Distance
}

func GetRandomVector(size int) []float32 {
	var vec []float32
	for i := 0; i < size; i++ {
		vec = append(vec, float32(rand.Float32()))
	}
	return vec
}
