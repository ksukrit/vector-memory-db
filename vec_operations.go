package main

import (
	"math"
	"sort"
)

// This is a separate file to deal with vector related operations
// like dot product, cosine similarity, sum etc.

func (db *InMemVectorDB) EuclideanDist(v1 Vector, v2 Vector) float32 {
	var sum float32
	for i := 0; i < db.embedding_size; i++ {
		sum += (v1[i] - v2[i]) * (v1[i] - v2[i])
	}
	return sum
}

func (db *InMemVectorDB) KNearestNeighbor(query Vector, k int, op_code int) []NeighborVector {
	var neighbors []NeighborVector

	for key, vector := range db.data {
		switch op_code {
		case 0:
			dist := db.EuclideanDist(query, vector)
			neighbors = append(neighbors, NeighborVector{key, vector, dist})
		case 1:
			sim := db.CosineSimilarity(query, vector)
			neighbors = append(neighbors, NeighborVector{key, vector, sim})
		}
	}

	switch op_code {
	case 0:
		sort.Sort(NeighborArr(neighbors))
	case 1:
		sort.Sort(sort.Reverse(NeighborArr(neighbors)))
	}

	if len(neighbors) > k {
		neighbors = neighbors[:k]
	}

	return neighbors
}

func (db *InMemVectorDB) CosineSimilarity(v1 Vector, v2 Vector) float32 {
	var sum float32
	var mag_1 float32
	var mag_2 float32
	for i := 0; i < db.embedding_size; i++ {
		sum += v1[i] * v2[i]
		mag_1 += v1[i] * v1[i]
		mag_2 += v2[i] * v2[i]
	}
	mag_1 = float32(math.Sqrt(float64(mag_1)))
	mag_2 = float32(math.Sqrt(float64(mag_2)))
	return sum / (mag_1 * mag_2)
}

func (db *InMemVectorDB) AddVector(key1, key2 string) Vector {
	vec1 := db.Get(key1)
	vec2 := db.Get(key2)
	if len(vec1) != len(vec2) {
		return Vector{}
	}
	vec_res := make([]float32, db.embedding_size)
	for i := 0; i < len(vec1); i++ {
		vec_res[i] = vec1[i] + vec2[i]
	}
	return vec_res
}

func (db *InMemVectorDB) SubVector(key1, key2 string) Vector {
	vec1 := db.Get(key1)
	vec2 := db.Get(key2)
	if len(vec1) != len(vec2) {
		return Vector{}
	}
	vec_res := make([]float32, db.embedding_size)
	for i := 0; i < len(vec1); i++ {
		vec_res[i] = vec1[i] - vec2[i]
	}
	return vec_res
}

func (db *InMemVectorDB) Scale_Vector(key string, scalar float32) Vector {
	vec1 := db.Get(key)
	for i := 0; i < len(vec1); i++ {
		vec1[i] = vec1[i] * scalar
	}
	return vec1
}
