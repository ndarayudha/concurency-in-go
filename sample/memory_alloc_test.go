package sample

import (
	"testing"
)

// Function with memory allocation
func allocateMemory() []int {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	return data
}

// Function without memory allocation (reuses a pre-allocated slice)
func reuseMemory(data []int) {
	for i := range data {
		data[i] = i
	}
}

var result []int

// Benchmark for the function that allocates memory
func BenchmarkWithAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result = allocateMemory()
	}
}

// Benchmark for the function that avoids allocations
func BenchmarkWithoutAllocation(b *testing.B) {
	data := make([]int, 100) // Pre-allocated memory
	b.ResetTimer()           // Reset the timer after setup
	for i := 0; i < b.N; i++ {
		reuseMemory(data)
	}
}
