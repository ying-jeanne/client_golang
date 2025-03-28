// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

var output float64

func TestAtomicUpdateFloat(t *testing.T) {
	var val float64 = 0.0
	bits := (*uint64)(unsafe.Pointer(&val))
	var wg sync.WaitGroup
	numGoroutines := 100000
	increment := 1.0

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomicUpdateFloat(bits, func(f float64) float64 {
				return f + increment
			})
		}()
	}

	wg.Wait()
	expected := float64(numGoroutines) * increment
	if val != expected {
		t.Errorf("Expected %f, got %f", expected, val)
	}
}

// Benchmark for atomicUpdateFloat with single goroutine (no contention).
func BenchmarkAtomicUpdateFloat_SingleGoroutine(b *testing.B) {
	var val float64 = 0.0
	bits := (*uint64)(unsafe.Pointer(&val))

	for i := 0; i < b.N; i++ {
		atomicUpdateFloat(bits, func(f float64) float64 {
			return f + 1.0
		})
	}

	output = val
}

// Benchmark for old implementation with single goroutine (no contention) -> to check overhead of backoff
func BenchmarkAtomicNoBackoff_SingleGoroutine(b *testing.B) {
	var val float64 = 0.0
	bits := (*uint64)(unsafe.Pointer(&val))

	for i := 0; i < b.N; i++ {
		for {
			loadedBits := atomic.LoadUint64(bits)
			newBits := math.Float64bits(math.Float64frombits(loadedBits) + 1.0)
			if atomic.CompareAndSwapUint64(bits, loadedBits, newBits) {
				break
			}
		}
	}

	output = val
}

// Benchmark varying the number of goroutines.
func benchmarkAtomicUpdateFloatConcurrency(b *testing.B, numGoroutines int) {
	var val float64 = 0.0
	bits := (*uint64)(unsafe.Pointer(&val))
	b.SetParallelism(numGoroutines)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomicUpdateFloat(bits, func(f float64) float64 {
				return f + 1.0
			})
		}
	})

	output = val
}

func benchmarkAtomicNoBackoffFloatConcurrency(b *testing.B, numGoroutines int) {
	var val float64 = 0.0
	bits := (*uint64)(unsafe.Pointer(&val))
	b.SetParallelism(numGoroutines)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for {
				loadedBits := atomic.LoadUint64(bits)
				newBits := math.Float64bits(math.Float64frombits(loadedBits) + 1.0)
				if atomic.CompareAndSwapUint64(bits, loadedBits, newBits) {
					break
				}
			}
		}
	})

	output = val
}

func BenchmarkAtomicUpdateFloat_1Goroutine(b *testing.B) {
	benchmarkAtomicUpdateFloatConcurrency(b, 1)
}

func BenchmarkAtomicNoBackoff_1Goroutine(b *testing.B) {
	benchmarkAtomicNoBackoffFloatConcurrency(b, 1)
}

func BenchmarkAtomicUpdateFloat_2Goroutines(b *testing.B) {
	benchmarkAtomicUpdateFloatConcurrency(b, 2)
}

func BenchmarkAtomicNoBackoff_2Goroutines(b *testing.B) {
	benchmarkAtomicNoBackoffFloatConcurrency(b, 2)
}

func BenchmarkAtomicUpdateFloat_4Goroutines(b *testing.B) {
	benchmarkAtomicUpdateFloatConcurrency(b, 4)
}

func BenchmarkAtomicNoBackoff_4Goroutines(b *testing.B) {
	benchmarkAtomicNoBackoffFloatConcurrency(b, 4)
}

func BenchmarkAtomicUpdateFloat_8Goroutines(b *testing.B) {
	benchmarkAtomicUpdateFloatConcurrency(b, 8)
}

func BenchmarkAtomicNoBackoff_8Goroutines(b *testing.B) {
	benchmarkAtomicNoBackoffFloatConcurrency(b, 8)
}

func BenchmarkAtomicUpdateFloat_16Goroutines(b *testing.B) {
	benchmarkAtomicUpdateFloatConcurrency(b, 16)
}

func BenchmarkAtomicNoBackoff_16Goroutines(b *testing.B) {
	benchmarkAtomicNoBackoffFloatConcurrency(b, 16)
}

func BenchmarkAtomicUpdateFloat_32Goroutines(b *testing.B) {
	benchmarkAtomicUpdateFloatConcurrency(b, 32)
}

func BenchmarkAtomicNoBackoff_32Goroutines(b *testing.B) {
	benchmarkAtomicNoBackoffFloatConcurrency(b, 32)
}
