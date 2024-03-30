package main

type Vector []float32

type NeighborVector struct {
	Key      string
	Vector   Vector
	Distance float32
}

func PadVector(v Vector, embedding_size int) Vector {
	if len(v) < embedding_size {
		for i := len(v); i < embedding_size; i++ {
			v = append(v, 0.0)
		}
	}
	return v
}
