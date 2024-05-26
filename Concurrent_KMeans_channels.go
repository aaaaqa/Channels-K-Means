package main

import (
	"fmt"
	"math"
	"math/rand"	
	"time"
)

type KMeans struct {
	nClusters int
	maxIters  int
	centroids [][]float64
	labels    []int
	data      [][]float64
}

func NewKMeans(nClusters, maxIters int, X [][]float64) *KMeans {
	return &KMeans{
		nClusters: nClusters,
		maxIters:  maxIters,
		data:      X,
		labels:    make([]int, len(X)),
	}
}

func euclideanDistance(a, b []float64) float64 {
	var sum float64
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func (kMeans *KMeans) assignLabels() {
	WorkGroups := 10
	N := len(kMeans.data) / WorkGroups
	results := make(chan [2]int, len(kMeans.data))

	for i := 0; i < WorkGroups; i++ {
		go func(start, end int, data [][]float64, results chan<- [2]int) {
			for k := start; k < end; k++ {
				minDist := math.MaxFloat64
				minIdx := 0
				for j, c := range kMeans.centroids {
					dist := euclideanDistance(data[k], c)
					if dist < minDist {
						minDist = dist
						minIdx = j
					}
				}
				results <- [2]int{k, minIdx}
			}
		}(i*N, (i+1)*N, kMeans.data, results)
	}

	for i := 0; i < len(kMeans.data); i++ {
		result := <-results
		kMeans.labels[result[0]] = result[1]
	}
	close(results)
}

func (kMeans *KMeans) updateCentroids() bool {
	newCentroids := make([][]float64, kMeans.nClusters)
	for i := range newCentroids {
		newCentroids[i] = make([]float64, len(kMeans.data[0]))
	}
	counts := make([]int, kMeans.nClusters)

	for i, label := range kMeans.labels {
		for j, val := range kMeans.data[i] {
			newCentroids[label][j] += val
		}
		counts[label]++
	}

	for i := range newCentroids {
		for j := range newCentroids[i] {
			if counts[i] > 0 {
				newCentroids[i][j] /= float64(counts[i])
			}
		}
	}

	if kMeans.checkConvergence(newCentroids) {
		return true
	}

	kMeans.centroids = newCentroids
	return false
}

func (kMeans *KMeans) checkConvergence(newCentroids [][]float64) bool {
	for i, c := range kMeans.centroids {
		for j, v := range c {
			if math.Abs(v-newCentroids[i][j]) > 1e-2 {
				return false
			}
		}
	}
	return true
}

func (kMeans *KMeans) Fit() {
	kMeans.centroids = make([][]float64, kMeans.nClusters)
	for i := range kMeans.centroids {
		kMeans.centroids[i] = make([]float64, len(kMeans.data[0]))
		for j := range kMeans.centroids[i] {
			kMeans.centroids[i][j] = kMeans.data[rand.Intn(len(kMeans.data))][j]
		}
	}

	for iter := 0; iter < kMeans.maxIters; iter++ {
		kMeans.assignLabels()
		if kMeans.updateCentroids() {
			fmt.Println("Converged.")
			break
		}
	}
}

func createArrayValues(min, max float64) [][]float64 {
	X := make([][]float64, 1000000)
	for i := range X {
		X[i] = make([]float64, 2)
		for j := range X[i] {
			X[i][j] = min + rand.Float64()*(max-min)
		}
	}
	return X
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var sum time.Duration
	totalDuration := make([]time.Duration, 1000)

	for i := 0; i < 1000; i++ {
		start := time.Now()
		X := createArrayValues(0.0, 1000.0)
		kmeans := NewKMeans(10, 100, X)
		kmeans.Fit()
		fmt.Println("Final Centroids:", kmeans.centroids)
		fmt.Println("Execution Time: ", time.Since(start))
		totalDuration[i] = time.Since(start)
	}

	totalDuration = totalDuration[50 : len(totalDuration)-51]
	for _, duration := range totalDuration {
		sum += duration
	}

	fmt.Println("\nAverage time: ", float64(sum)/float64(len(totalDuration)))
}