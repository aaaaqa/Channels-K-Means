package kmeans

import (
	"fmt"
	"math"
	"math/rand"
)

type KMeans struct {
	data       [][]float64
	centroids  [][]float64
	labels     []int
	nClusters  int
	maxIters   int
}

func NewKMeans(data [][]float64, nClusters, maxIters int) *KMeans {
	return &KMeans{
		data:      data,
		nClusters: nClusters,
		maxIters:  maxIters,
	}
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

func (kMeans *KMeans) assignLabels() {
	kMeans.labels = make([]int, len(kMeans.data))
	N := len(kMeans.data) / 4
	results := make(chan [2]int, len(kMeans.data))

	for i := 0; i < 4; i++ {
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

func euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(sum)
}

func (kMeans *KMeans) Centroids() [][]float64 {
	return kMeans.centroids
}

func (kMeans *KMeans) Labels() []int {
	return kMeans.labels
}
