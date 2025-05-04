package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// scatterPlotWithRegression creates a scatter plot with a regression line for the given dataset.
func scatterPlotWithRegression(dataset DataSet, slope, intercept float64) {
	p := plot.New()

	p.Title.Text = "Anscombe Dataset " + dataset.name
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create scatter points
	pts := make(plotter.XYs, len(dataset.x))
	for i := range dataset.x {
		pts[i].X = dataset.x[i]
		pts[i].Y = dataset.y[i]
	}

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatalf("failed to create scatter: %v", err)
	}
	scatter.GlyphStyle.Radius = vg.Points(3)

	// Create regression line using min and max x
	minX, maxX := dataset.x[0], dataset.x[0]
	for _, val := range dataset.x {
		if val < minX {
			minX = val
		}
		if val > maxX {
			maxX = val
		}
	}
	line := plotter.XYs{
		{X: minX, Y: slope*minX + intercept},
		{X: maxX, Y: slope*maxX + intercept},
	}
	regLine, err := plotter.NewLine(line)
	if err != nil {
		log.Fatalf("failed to create line: %v", err)
	}
	regLine.LineStyle.Width = vg.Points(1)

	p.Add(scatter, regLine)
	p.Legend.Add("Regression Line", regLine)
	p.Legend.Add("Data Points", scatter)

	// Save to PNG
	filename := fmt.Sprintf("anscombe_%s.png", dataset.name)
	if err := p.Save(5*vg.Inch, 5*vg.Inch, filename); err != nil {
		log.Fatalf("failed to save plot: %v", err)
	}
	fmt.Printf("  Plot saved to %s\n", filename)
}

// DataSet holds the x, y data and the name of the dataset.
type DataSet struct {
	name string
	x    []float64
	y    []float64
}

func main() {
	// Define the Anscombe Quartet datasets.
	datasets := []DataSet{
		{
			name: "I",
			x:    []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7, 5},
			y:    []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82, 5.68},
		},
		{
			name: "II",
			x:    []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7, 5},
			y:    []float64{9.14, 8.14, 8.74, 8.77, 9.26, 8.1, 6.13, 3.1, 9.13, 7.26, 4.74},
		},
		{
			name: "III",
			x:    []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7, 5},
			y:    []float64{7.46, 6.77, 12.74, 7.11, 7.81, 8.84, 6.08, 5.39, 8.15, 6.42, 5.73},
		},
		{
			name: "IV",
			x:    []float64{8, 8, 8, 8, 8, 8, 8, 19, 8, 8, 8},
			y:    []float64{6.58, 5.76, 7.71, 8.84, 8.47, 7.04, 5.25, 12.5, 5.56, 7.91, 6.89},
		},
	}

	for _, ds := range datasets {
		fmt.Printf("\nAnalyzing Dataset %s\n", ds.name)

		start := time.Now()
		var memStart, memEnd runtime.MemStats
		runtime.ReadMemStats(&memStart)

		slope, intercept := linearRegression(ds.x, ds.y)
		r2 := rSquared(ds.x, ds.y, slope, intercept)
		rse := residualStandardError(ds.x, ds.y, slope, intercept)
		fstat := fStatistic(r2, len(ds.x))

		runtime.ReadMemStats(&memEnd)
		duration := time.Since(start)
		memUsed := memEnd.Alloc - memStart.Alloc

		fmt.Printf("  Slope               : %.4f\n", slope)
		fmt.Printf("  Intercept           : %.4f\n", intercept)
		fmt.Printf("  R-squared           : %.4f\n", r2)
		fmt.Printf("  Residual Std. Error : %.4f\n", rse)
		fmt.Printf("  F-statistic         : %.4f\n", fstat)
		fmt.Printf("  Time (ms)           : %.4f\n", float64(duration.Microseconds())/1000.0)
		fmt.Printf("  Memory (bytes)      : %d\n", memUsed)

		scatterPlotWithRegression(ds, slope, intercept)
	}
}

// linearRegression calculates the slope and intercept for simple linear regression.
func linearRegression(x, y []float64) (slope, intercept float64) {
	n := float64(len(x))
	var sumX, sumY, sumXY, sumX2 float64

	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	xMean := sumX / n
	yMean := sumY / n

	numerator := sumXY - n*xMean*yMean
	denominator := sumX2 - n*xMean*xMean

	slope = numerator / denominator
	intercept = yMean - slope*xMean

	return
}

// rSquared calculates the coefficient of determination.
func rSquared(x, y []float64, slope, intercept float64) float64 {
	var ssTot, ssRes float64
	yMean := mean(y)

	for i := range x {
		predicted := slope*x[i] + intercept
		ssTot += (y[i] - yMean) * (y[i] - yMean)
		ssRes += (y[i] - predicted) * (y[i] - predicted)
	}

	return 1 - (ssRes / ssTot)
}

// residualStandardError calculates the RSE.
func residualStandardError(x, y []float64, slope, intercept float64) float64 {
	var ssRes float64
	n := float64(len(x))

	for i := range x {
		predicted := slope*x[i] + intercept
		ssRes += math.Pow(y[i]-predicted, 2)
	}

	return math.Sqrt(ssRes / (n - 2))
}

// fStatistic calculates the F-statistic.
func fStatistic(r2 float64, n int) float64 {
	if r2 == 1.0 {
		return math.Inf(1)
	}
	return (r2 / (1 - r2)) * float64(n-2)
}

// mean computes the mean of a slice.
func mean(vals []float64) float64 {
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}
