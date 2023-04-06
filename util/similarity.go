package util

import "math"

func dotProduct(a, b []float32) float32 {
	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func norm(a []float32) float32 {
	var sumSquares float32
	for _, x := range a {
		sumSquares += x * x
	}
	return float32(math.Sqrt(float64(sumSquares)))
}

func CosineSimilarity(a, b []float32) float32 {
	dotProduct := dotProduct(a, b)
	normA := norm(a)
	normB := norm(b)
	return dotProduct / (normA * normB)
}
