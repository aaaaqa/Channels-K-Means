package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	file, err := os.Open("dataset.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, strings.TrimSpace(line))
	}
	if len(lines) < 2 {
		fmt.Println("Not enough data in file")
		return
	}

	x1Str := strings.Split(lines[0], ", ")
	x2Str := strings.Split(lines[1], ", ")

	var x1, x2 []float64
	for _, val := range x1Str {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println("Error parsing float:", err)
			return
		}
		x1 = append(x1, f)
	}
	for _, val := range x2Str {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println("Error parsing float:", err)
			return
		}
		x2 = append(x2, f)
	}

	data := make([][]float64, len(x1))
	for i := range x1 {
		data[i] = []float64{x1[i], x2[i]}
	}

	K := 10
	distortions := make([]float64, K)
	inertias := make([]float64, K)

	for k := 1; k <= K; k++ {
		centroids, clusters := kmeans(data, k)

		distortion := 0.0
		for i := 0; i < len(data); i++ {
			minDist := math.MaxFloat64
			for _, centroid := range centroids {
				dist := euclideanDist(data[i], centroid)
				if dist < minDist {
					minDist = dist
				}
			}
			distortion += minDist
		}
		distortion /= float64(len(data))
		distortions[k-1] = distortion

		inertia := 0.0
		for i, cluster := range clusters {
			for _, point := range cluster {
				inertia += euclideanDist(point, centroids[i])
			}
		}
		inertias[k-1] = inertia
	}

	fmt.Println("Distortions:", distortions)
	fmt.Println("Inertias:", inertias)
}

func kmeans(data [][]float64, k int) ([][]float64, [][][]float64) {
	rand.Seed(time.Now().UnixNano())

	centroids := make([][]float64, k)
	for i := 0; i < k; i++ {
		centroids[i] = data[rand.Intn(len(data))]
	}

	clusters := make([][][]float64, k)
	for {
		for i := range clusters {
			clusters[i] = nil
		}

		for _, point := range data {
			minDist := math.MaxFloat64
			clusterIdx := 0
			for i, centroid := range centroids {
				dist := euclideanDist(point, centroid)
				if dist < minDist {
					minDist = dist
					clusterIdx = i
				}
			}
			clusters[clusterIdx] = append(clusters[clusterIdx], point)
		}

		newCentroids := make([][]float64, k)
		for i, cluster := range clusters {
			newCentroids[i] = mean(cluster)
		}

		converged := true
		for i := range centroids {
			if !equal(centroids[i], newCentroids[i]) {
				converged = false
				break
			}
		}
		if converged {
			break
		}

		centroids = newCentroids
	}

	return centroids, clusters
}

func euclideanDist(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func mean(data [][]float64) []float64 {
	if len(data) == 0 {
		return nil
	}
	mean := make([]float64, len(data[0]))
	for _, point := range data {
		for i, val := range point {
			mean[i] += val
		}
	}
	for i := range mean {
		mean[i] /= float64(len(data))
	}
	return mean
}

func equal(a, b []float64) bool {
	const epsilon = 1e-9
	for i := range a {
		if math.Abs(a[i]-b[i]) > epsilon {
			return false
		}
	}
	return true
}
