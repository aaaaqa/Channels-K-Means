package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	//"os"
	"strconv"
	"strings"
)

var centroids []string
var labels []string

func readAsArray(s string) []float64 {
	// X := make([]float64, 1000000)
	var X []float64
	S := strings.Split(s, ", ")
	for _, s := range S {
		temp, _ := strconv.ParseFloat(s, 64)
		X = append(X, temp)
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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Read data from centroids as an array
	scanner := bufio.NewReader(conn)
	tempCentroids, _ := scanner.ReadString('\n')
	tempCentroids = strings.TrimSpace(tempCentroids)
	tempCentroids = tempCentroids[:len(tempCentroids)-1] // Erase last comma
	tempCentroidsArray := strings.Split(tempCentroids, ",")

	// Read data from labels as an array
	tempLabels, _ := scanner.ReadString('\n')
	tempLabels = strings.Trim(tempLabels, "[]")
	tempLabelsArray := strings.Split(tempLabels, " ")
	centroids = tempCentroidsArray
	labels = tempLabelsArray
	fmt.Println("Data recibida:", centroids)
	fmt.Println("Data recibida:", labels[:5])
}

func trainKmeans(res http.ResponseWriter, req *http.Request) {
	// Build response
	log.Println("KMeans empezo entrenamiento...")
	res.Header().Set("Content-Type", "application/json")

	// Fetching the data
	url := "https://raw.githubusercontent.com/aaaaqa/Channels-K-Means/main/dataset.txt"
	incomes, ages, err := fetchDataset(url)
	if err != nil {
		fmt.Println("Error buscando dataset:", err)
		return
	}

	data := make([][]float64, len(incomes))
	for i := range incomes {
		data[i] = []float64{incomes[i], ages[i]}
	}

	// Open connection with the worker nodes
	addresses := []string{host + ":8001"} // Add server addresses
	for _, address := range addresses {
		err = sendData(data, address)
		if err != nil {
			io.WriteString(res, "Error enviando datos al nodo.")
			log.Println("Error enviado datos a", address, ":8001", err)
			return
		}
	}

	// Await a response
	ln, err := net.Listen("tcp", ":8002")
	if err != nil {
		io.WriteString(res, "Error empezando servidor.")
		log.Println("Error empezando servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Println("El cliente esta escuchando...")

	conn, err := ln.Accept()
	if err != nil {
		io.WriteString(res, "Error aceptando conexion.")
		log.Println("Error aceptando conexion:", err)
		return
	} else {
		handleConnection(conn)
	}

	// Success message
	io.WriteString(res, "Kmeans entrenado exitosamente!")
	log.Println("Kmeans entrenado exitosamente!")
}

func getLabels(res http.ResponseWriter, req *http.Request) {
	// Build response
	log.Println("Buscando labels...")
	res.Header().Set("Content-Type", "application/json")
	// Check for errors
	if labels == nil {
		io.WriteString(res, "No hay labels. Entrena el modelo primero.")
		log.Println("No hay labels. Entrena el modelo primero.")
		return
	}
	// Respond as json
	jsonBytes, _ := json.Marshal(labels)
	io.WriteString(res, string(jsonBytes))
	log.Println("Labels encontrados!")
}

func getCentroids(res http.ResponseWriter, req *http.Request) {
	// Build response
	log.Println("Buscando centroides...")
	res.Header().Set("Content-Type", "application/json")
	// Check for errors
	if labels == nil {
		io.WriteString(res, "No hay centroides. Entrena el modelo primero.")
		log.Println("No hay centroides. Entrena el modelo primero.")
		return
	}
	// Respond as json
	jsonBytes, _ := json.Marshal(centroids)
	io.WriteString(res, string(jsonBytes))
	log.Println("Centroides encontrados!")
}

func handleRequests() {
	// Abrir los endpoints
	http.HandleFunc("/train", trainKmeans)
	http.HandleFunc("/getLabels", getLabels)
	http.HandleFunc("/getCentroids", getCentroids)

	// Abrir api
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
