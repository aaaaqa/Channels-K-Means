package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var host string = "192.168.1.39"

func hostClient(port int) {
	// Abrir el host para escuchar los resultados de los nodos
	ls, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error hosteando cliente:", err)
		return
	}
	defer ls.Close()

	fmt.Println("Servidor escuchando...")

	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("Error en conexion:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)

	// ESTA PARTE DEBERIA CAMBIARSE PARA RECIBIR
	// LOS DATOS PARA ACTUALIZAR LOS CENTROIDES

	// Leer numero de clientes
	numRowsStr, err := br.ReadString('\n')
	if err != nil {
		fmt.Println("Error en entrada de filas:", err)
		return
	}
	numRows, err := strconv.Atoi(numRowsStr[:len(numRowsStr)-1])
	if err != nil {
		fmt.Println("Error convirtiendo numero de filas:", err)
		return
	}

	// Leer cada fila de datos
	for i := 0; i < numRows; i++ {
		rowStr, err := br.ReadString('\n')
		if err != nil {
			fmt.Println("Error leyendo datos:", err)
			return
		}

		var cliente []float64
		words := strings.Fields(rowStr)
		for _, word := range words {
			// 64 for float64
			val, err := strconv.ParseFloat(word, 64)
			if err != nil {
				fmt.Println("Error convirtiendo valor:", err)
				return
			}
			cliente = append(cliente, val)
		}
		// DE LO QUE ESCUCHAMOS DEBEMOS ACTUALIZAR CENTROIDES ACA
		// CUIDADO CON MANTENER EL INDICE DE LOS CLIENTES ENVIADOS
		fmt.Println(cliente)
	}
}

func sendData(data [][]float64, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Enviar el numero de clientes
	fmt.Println("Enviando datos")
	fmt.Println(data)
	br := bufio.NewWriter(conn)
	fmt.Fprintf(br, "%d\n", len(data))

	// Enviar los datos por cliente
	for _, cliente := range data {
		for _, val := range cliente {
			fmt.Fprintf(br, "%f ", val)
		}
		fmt.Fprintln(br)
	}

	br.Flush()
	return nil
}

func main() {
	// array de prueba para probar, probare con 3 por ahora
	data := [][]float64{
		{20, 1000},
		{30, 4000},
		{40, 7000},
		{50, 8000},
		{60, 10000},
		{70, 15000}}

	// ejemplo de data particionada
	// preferiblemente se podrian usar arreglos dinamicos
	// pero no hay necesidad para la prueba
	var dataP [][][]float64
	for i := 0; i < 6; i += 2 {
		dataP = append(dataP, data[i:i+2])
	}

	ports := []string{host + ":8001", host + ":8002", host + ":8003"}
	for i, node := range ports {
		err := sendData(dataP[i], node)
		if err != nil {
			fmt.Println("Error enviando datos al nodo:", err)
		}
	}
	hostClient(8000)
}
