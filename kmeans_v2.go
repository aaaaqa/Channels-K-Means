package kmeans

import (	
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
	}

	for i := 0; i < kMeans.nClusters; i++ {
		idx := rand.Intn(len(kMeans.data))
		copy(kMeans.centroids[i], kMeans.data[idx])
	}

	for iter := 0; iter < kMeans.maxIters; iter++ {
		kMeans.labels = make([]int, len(kMeans.data))
		for i, point := range kMeans.data {
			kMeans.labels[i] = kMeans.closestCentroid(point)
		}

		newCentroids := make([][]float64, kMeans.nClusters)
		counts := make([]int, kMeans.nClusters)
		for i := range newCentroids {
			newCentroids[i] = make([]float64, len(kMeans.data[0]))
		}

		for i, point := range kMeans.data {
			label := kMeans.labels[i]
			for j := range point {
				newCentroids[label][j] += point[j]
			}
			counts[label]++
		}

		for i := range newCentroids {
			if counts[i] == 0 {
				continue
			}
			for j := range newCentroids[i] {
				newCentroids[i][j] /= float64(counts[i])
			}
		}

		kMeans.centroids = newCentroids
	}
}

func (kMeans *KMeans) closestCentroid(point []float64) int {
	minDist := math.MaxFloat64
	label := 0
	for i, centroid := range kMeans.centroids {
		dist := kMeans.euclideanDistance(point, centroid)
		if dist < minDist {
			minDist = dist
			label = i
		}
	}
	return label
}

func (kMeans *KMeans) euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

func (kMeans *KMeans) Labels() []int {
	return kMeans.labels
}

func (kMeans *KMeans) Centroids() [][]float64 {
	return kMeans.centroids
}