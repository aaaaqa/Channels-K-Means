package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var host string = "192.168.1.39"

func hostNode(port string) error {
	hostname := host + ":" + port
	fmt.Println(hostname)
	ls, err := net.Listen("tcp", hostname)
	if err != nil {
		return err
	}
	defer ls.Close()

	fmt.Println("Nodo escuchando...")

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
	var data [][]float64
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
		data = append(data, cliente)
	}

	// EN ESTA SECCION SE DEBE PROCESAR LOS DATOS
	// PARA ACTUALIZAR LOS CENTROIDES
	fmt.Println("Datos enviados:")
	fmt.Println(data)

	// DESPUES SE DEBEN ENVIAR LOS DATOS DE VUELTA AL CLIENTE
	sendData(data, host+":8000")
}

func sendData(data [][]float64, address string) error {
	// LO QUE ESTA ACA ES PARA PRUEBAS, DEBERIA DEVOLVER LOS
	// DATOS PARA ACTUALIZAR LOS CENTROIDES

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	br := bufio.NewWriter(conn)
	// Numero de clientes
	fmt.Fprintf(br, "%d\n", len(data))

	// Enviar cada fila de datos
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
	// ingresar puerto
	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto: ")
	port, _ := br.ReadString('\n')
	port = strings.TrimSpace(port)

	err := hostNode(port)
	if err != nil {
		fmt.Println("Error hosteando el nodo:", err)
	}
}
