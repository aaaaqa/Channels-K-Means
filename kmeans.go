package kmeans

import (
	"math"
	"math/rand"
	"sync"
)

type Point struct {
	Age    float64
	Income float64
}

type Centroid struct {
	Age    float64
	Income float64
}

var numWorkers  = 4

func InitializeCentroids(points []Point, k int) []Centroid {
	centroids := make([]Centroid, k)
	for i := range centroids {
		centroids[i] = Centroid{
			Age:    points[rand.Intn(len(points))].Age,
			Income: points[rand.Intn(len(points))].Income,
		}
	}
	return centroids
}

func AssignPointsToClusters(points []Point, centroids []Centroid, assignments chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done()
	clusterAssignments := make([]int, len(points))
	for i, point := range points {
		minDist := math.MaxFloat64
		minIndex := 0
		for j, centroid := range centroids {
			dist := Distance(point, centroid)
			if dist < minDist {
				minDist = dist
				minIndex = j
			}
		}
		clusterAssignments[i] = minIndex
	}
	assignments <- clusterAssignments
}

func Distance(p1, p2 []float64) float64 {
	if len(p1) != len(p2) {
		return math.MaxFloat64 // Return a large number if dimensions do not match
	}
	var sum float64
	for i := range p1 {
		diff := p1[i] - p2[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func RecalculateCentroids(points []Point, assignments []int, k int) []Centroid {
	sumAges := make([]float64, k)
	sumIncomes := make([]float64, k)
	counts := make([]int, k)

	for i, point := range points {
		cluster := assignments[i]
		sumAges[cluster] += point.Age
		sumIncomes[cluster] += point.Income
		counts[cluster]++
	}

	newCentroids := make([]Centroid, k)
	for i := range newCentroids {
		if counts[i] == 0 {
			newCentroids[i] = Centroid{
				Age:    points[rand.Intn(len(points))].Age,
				Income: points[rand.Intn(len(points))].Income,
			}
		} else {
			newCentroids[i] = Centroid{
				Age:    sumAges[i] / float64(counts[i]),
				Income: sumIncomes[i] / float64(counts[i]),
			}
		}
	}
	return newCentroids
}

func KMeans(points []Point, k, maxIter int) []Centroid {
	centroids := InitializeCentroids(points, k)
	assignments := make(chan []int, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < maxIter; i++ {
		wg.Add(numWorkers)
		for w := 0; w < numWorkers; w++ {
			go AssignPointsToClusters(points[w*len(points)/numWorkers:(w+1)*len(points)/numWorkers], centroids, assignments, &wg)
		}
		wg.Wait()
		close(assignments)

		allAssignments := make([]int, len(points))
		for assignment := range assignments {
			for j, a := range assignment {
				allAssignments[j] = a
			}
		}

		newCentroids := RecalculateCentroids(points, allAssignments, k)
		if CentroidsEqual(centroids, newCentroids) {
			break
		}
		centroids = newCentroids
	}
	return centroids
}

func CentroidsEqual(a, b []Centroid) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}