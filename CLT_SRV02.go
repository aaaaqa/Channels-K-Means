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
	fmt.Println("Data received:", data[:5])

	k := 3
	maxIter := 100
	kmeansInstance := kmeans.NewKMeans(data, k, maxIter)
	kmeansInstance.Fit()

	fmt.Println("Centroids:", kmeansInstance.Centroids())
	fmt.Println("Labels:", kmeansInstance.Labels())
}

func main() {
	ln, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server is listening...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
