package main

import (
	"fmt"
	"time"
)

func RunBenchmark() {

	db := NewVecMemDB("test", 3072)

	// Start timer
	start := time.Now()

	for i := 0; i < 100_000; i++ {
		db.Insert(GetRandomVector(3072))
	}

	fmt.Println(db.Size())

	// perform 10 K-NN searches
	for i := 0; i < 10; i++ {
		db.KNearestNeighbor(GetRandomVector(3072), 5, 0)
	}

	// end timer
	end := time.Now()

	fmt.Printf("Time taken: %s\n", end.Sub(start).String())
}
