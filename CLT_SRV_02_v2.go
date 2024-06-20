package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"kmeans"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	var data [][]float64
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		row := make([]float64, len(fields))
		for i, field := range fields {
			row[i], _ = strconv.ParseFloat(field, 64)
		}
		data = append(data, row)
	}

	if len(data) > 0 {
		nClusters := 3
		maxIters := 100
		kmeans := kmeans.NewKMeans(data, nClusters, maxIters)
		kmeans.Fit()

		labels := kmeans.Labels()
		centroids := kmeans.Centroids()

		// Send the results back to the client
		for _, label := range labels {
			fmt.Fprintf(conn, "%d\n", label)
		}

		fmt.Fprintf(conn, "Centroids:\n")
		for _, centroid := range centroids {
			for _, val := range centroid {
				fmt.Fprintf(conn, "%f ", val)
			}
			fmt.Fprint(conn, "\n")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8001...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}