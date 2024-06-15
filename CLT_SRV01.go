package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	//"os"
	"strconv"
	"strings"
)

var host string = "192.168.1.47"

func fetchDataset(url string) ([]float64, []float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	incomeLine := scanner.Text()
	scanner.Scan()
	ageLine := scanner.Text()

	incomes := strings.Split(incomeLine, ", ")
	ages := strings.Split(ageLine, ", ")

	var incomeData []float64
	var ageData []float64
	for _, incomeStr := range incomes {
		income, err := strconv.ParseFloat(incomeStr, 64)
		if err != nil {
			return nil, nil, err
		}
		incomeData = append(incomeData, income)
	}
	for _, ageStr := range ages {
		age, err := strconv.ParseFloat(ageStr, 64)
		if err != nil {
			return nil, nil, err
		}
		ageData = append(ageData, age)
	}
	return incomeData, ageData, nil
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