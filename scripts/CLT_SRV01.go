package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	//"os"
	"html/template"
	"strconv"
	"strings"
)

type Variables struct {
	Label_count map[int]int
	Centroid    []string
	Income      []float64
	Age         []float64
	Label       []string
}

var centroids []string
var labels []string
var incomes []float64
var ages []float64

func readAsArray(s string) []float64 {
	// X := make([]float64, 1000000)
	var X []float64
	S := strings.Split(strings.TrimSpace(s), ", ")
	for _, s := range S {
		temp, _ := strconv.ParseFloat(s, 64)
		X = append(X, temp)
	}

	return X
}

var remote []string
var currentWorker = 0

func fetchDataset(url string) ([]float64, []float64, error) {
	response, err := http.Get(url)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	response.Body.Close()

	incomes = readAsArray(string(body)[:strings.IndexByte(string(body), '\n')])
	ages = readAsArray(string(body)[strings.IndexByte(string(body), '\n'):])
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

	var fileName = "train.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Error parseando el archivo: ", err)
		return
	}
	err = t.ExecuteTemplate(res, fileName, nil)
	if err != nil {
		fmt.Println("Error ejecutando el archivo: ", err)
		return
	}

	log.Println("KMeans empezo entrenamiento...")
	//res.Header().Set("Content-Type", "application/json")

	// Fetching the data
	url := "https://raw.githubusercontent.com/aaaaqa/Channels-K-Means/main/dataset1.txt"
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
	addresses := []string{remote[currentWorker] + ":8001"} // Add server addresses
	currentWorker = (currentWorker + 1) % 2
	for _, address := range addresses {
		err = sendData(data, address)
		if err != nil {
			io.WriteString(res, "Error enviando datos al nodo.")
			log.Println("Error enviado los datos...", err)
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
	//io.WriteString(res, "Kmeans entrenado exitosamente!")
	log.Println("Kmeans entrenado exitosamente!")
}

func getLabels(res http.ResponseWriter, req *http.Request) {
	log.Println("Buscando labels...")

	// Check for errors
	if labels == nil {
		http.Error(res, "No hay labels. Entrena el modelo primero.", http.StatusNotFound)
		log.Println("No hay labels. Entrena el modelo primero.")
		return
	}

	// Count labels
	labelCount := make(map[int]int)
	for _, label := range labels {
		i, err := strconv.Atoi(strings.TrimSpace(strings.Replace(label, "]", "", -1)))
		if err != nil {
			log.Printf("Error al convertir label '%s' a entero: %v", label, err)
			continue
		}
		labelCount[i] = labelCount[i] + 1
	}

	// Build response
	var variable Variables
	variable = Variables{
		Label_count: labelCount,
		Income:      incomes,
		Age:         ages,
		Label:       labels,
	}

	log.Println("Labels encontrados!")

	const fileName = "labels.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		log.Printf("Error parseando el archivo: %v", err)
		http.Error(res, "Error parseando el archivo", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(res, fileName, variable)
	if err != nil {
		log.Printf("Error ejecutando el archivo: %v", err)
		http.Error(res, "Error ejecutando el archivo", http.StatusInternalServerError)
		return
	}
}

func getCentroids(res http.ResponseWriter, req *http.Request) {
	log.Println("Buscando centroides...")

	// Check for errors
	if labels == nil {
		http.Error(res, "No hay centroides. Entrena el modelo primero.", http.StatusNotFound)
		log.Println("No hay centroides. Entrena el modelo primero.")
		return
	}

	// Build response
	var variable Variables
	variable = Variables{
		Centroid: centroids,
		Income:   incomes,
		Age:      ages,
		Label:    labels,
	}

	const fileName = "centroids.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		log.Printf("Error parseando el archivo: %v", err)
		http.Error(res, "Error parseando el archivo", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(res, fileName, variable)
	if err != nil {
		log.Printf("Error ejecutando el archivo: %v", err)
		http.Error(res, "Error ejecutando el archivo", http.StatusInternalServerError)
		return
	}

	log.Println("Centroides encontrados y enviados!")
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	var fileName = "main.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Error parseando el archivo: ", err)
		return
	}
	err = t.ExecuteTemplate(w, fileName, nil)
	if err != nil {
		fmt.Println("Error ejecutando el archivo: ", err)
		return
	}
}

func handleRequests() {
	// Abrir los endpoints
	http.HandleFunc("/train", trainKmeans)
	http.HandleFunc("/getLabels", getLabels)
	http.HandleFunc("/getCentroids", getCentroids)
	http.HandleFunc("/", mainPage)
	// Abrir api
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	remote = append(remote, "172.212.83.254")
	remote = append(remote, "48.217.83.50")
	handleRequests()
}
