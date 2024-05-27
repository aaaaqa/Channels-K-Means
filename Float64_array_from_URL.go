package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

func readAsArray(s string) []float64 {
	X := make([]float64, 1000000)
	S := strings.Fields(s)
	for i, number := range S {
		X[i], _ = strconv.ParseFloat(strings.Replace(number, ",", "", -1), 64)
	}

	return X
}

func main() {

	response, err := http.Get("https://raw.githubusercontent.com/aaaaqa/Channels-K-Means/main/dataset.txt")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	response.Body.Close()

	X := readAsArray(string(body)[:strings.IndexByte(string(body), '\n')])
	Y := readAsArray(string(body)[strings.IndexByte(string(body), '\n'):])

	fmt.Println(X[:5])
	fmt.Println(Y[:5])

}
