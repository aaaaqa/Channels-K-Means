package main

import (
	"bufio"
	"fmt"
	"kmeans"
	"net"
	"strconv"
	"strings"
)

var host string = "26.114.63.141"

func sendData(centroids [][]float64, labels []int, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Println("Centroids:", centroids)

	writer := bufio.NewWriter(conn)
	for _, centroid := range centroids {
		fmt.Fprintf(writer, "%f,", centroid)
	}
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, labels)

	writer.Flush()
	return nil
}

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

	k := 5
	maxIter := 100
	kmeansInstance := kmeans.NewKMeans(data, k, maxIter)
	kmeansInstance.Fit()

	err := sendData(kmeansInstance.Centroids(), kmeansInstance.Labels(), host+":8002")
	if err != nil {
		fmt.Println("Error sending data to", host+":8002", ":", err)
	}
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
