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
	var points []kmeans.Point
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		if len(values) == 2 {
			income, _ := strconv.ParseFloat(values[0], 64)
			age, _ := strconv.ParseFloat(values[1], 64)
			points = append(points, kmeans.Point{Age: age, Income: income})
		}
	}
	centroids := kmeans.KMeans(points, 3, 100)
	fmt.Println("Calculated centroids:", centroids)
}

func main() {
	port := "8001" // Change to appropriate port
	ls, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error hosting server:", err)
		return
	}
	defer ls.Close()

	fmt.Println("Server listening on port", port)

	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("Error in connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}