package main

import (
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

// TestLinearRegressionBasic checks regression against a perfect line y = 2x + 1
func TestLinearRegressionBasic(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{3, 5, 7, 9, 11}

	slope, intercept := linearRegression(x, y)

	if math.Abs(slope-2.0) > 1e-6 {
		t.Errorf("Expected slope 2.0, got %.6f", slope)
	}
	if math.Abs(intercept-1.0) > 1e-6 {
		t.Errorf("Expected intercept 1.0, got %.6f", intercept)
	}
}

// TestRSquaredPerfectLine should return R² = 1
func TestRSquaredPerfectLine(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{3, 5, 7, 9, 11}
	slope, intercept := linearRegression(x, y)
	r2 := rSquared(x, y, slope, intercept)

	if math.Abs(r2-1.0) > 1e-6 {
		t.Errorf("Expected R² 1.0, got %.6f", r2)
	}
}

// TestRSEZeroResidual checks that RSE = 0 when residuals are zero
func TestRSEZeroResidual(t *testing.T) {
	x := []float64{1, 2, 3}
	y := []float64{2, 4, 6}
	slope, intercept := linearRegression(x, y)
	rse := residualStandardError(x, y, slope, intercept)

	if math.Abs(rse) > 1e-10 {
		t.Errorf("Expected RSE 0, got %.10f", rse)
	}
}

// BenchmarkLinearRegression runs Monte Carlo-style performance benchmarks
func BenchmarkLinearRegression(b *testing.B) {
	const n = 10000
	x := make([]float64, n)
	y := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = float64(i)
		y[i] = 3*x[i] + 7 + rand.NormFloat64()*0.5 // y = 3x + 7 + noise
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		linearRegression(x, y)
	}
}

// TestMemoryAndTiming measures basic memory allocation and timing
func TestMemoryAndTiming(t *testing.T) {
	const n = 100000
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i)
		y[i] = 5*x[i] - 2 + rand.NormFloat64()
	}

	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)

	start := time.Now()
	slope, intercept := linearRegression(x, y)
	duration := time.Since(start)

	runtime.ReadMemStats(&mEnd)
	memUsed := mEnd.Alloc - mStart.Alloc

	t.Logf("Slope: %.2f, Intercept: %.2f", slope, intercept)
	t.Logf("Time taken: %s", duration)
	t.Logf("Memory used: %d bytes", memUsed)
}
