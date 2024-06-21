package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	//"os"
	"strconv"
	"strings"
)

func readAsArray(s string) []float64 {
	X := make([]float64, 1000000)
	S := strings.Fields(s)
	for i, number := range S {
		X[i], _ = strconv.ParseFloat(strings.Replace(number, ",", "", -1), 64)
	}

	return X
}

var host string = "26.114.63.141"

func fetchDataset(url string) ([]float64, []float64, error) {
	response, err := http.Get(url)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	response.Body.Close()

	incomes := readAsArray(string(body)[:strings.IndexByte(string(body), '\n')])
	ages := readAsArray(string(body)[strings.IndexByte(string(body), '\n'):])
	fmt.Println(incomes[:5])
	fmt.Println(ages[:5])

	return incomes, ages, nil
}

func sendData(data [][]float64, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	for _, row := range data {
		for _, val := range row {
			fmt.Fprintf(writer, "%f ", val)
		}
		fmt.Fprintln(writer)
	}
	writer.Flush()
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewReader(conn)
	centroids, _ := scanner.ReadString('\n')
	temp_labels, _ := scanner.ReadString('\n')
	temp_labels = strings.Trim(temp_labels, "[]")
	labels := strings.Split(temp_labels, " ")

	fmt.Println("Data received:", centroids)
	fmt.Println("Data received:", labels[:5])
}

func main() {
	url := "https://raw.githubusercontent.com/aaaaqa/Channels-K-Means/main/dataset.txt"
	incomes, ages, err := fetchDataset(url)
	if err != nil {
		fmt.Println("Error fetching dataset:", err)
		return
	}

	data := make([][]float64, len(incomes))
	for i := range incomes {
		data[i] = []float64{incomes[i], ages[i]}
	}

	addresses := []string{host + ":8001"} // Add server addresses
	for _, address := range addresses {
		err = sendData(data, address)
		if err != nil {
			fmt.Println("Error sending data to", address, ":", err)
		}
	}

	ln, err := net.Listen("tcp", ":8002")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Client is listening...")

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	} else {
		handleConnection(conn)
	}
}
