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

var host string = "192.168.1.39"

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

	addresses := []string{host + ":8001", host + ":8002"} // Add server addresses
	for _, address := range addresses {
		err = sendData(data, address)
		if err != nil {
			fmt.Println("Error sending data to", address, ":", err)
		}
	}
}
