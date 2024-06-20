package main

import (
    "fmt"
    "os"
    "math/rand"
    "strconv"
    "math"
)

func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}

func main() {
    N := 1000000
    min := 1000.0
    max := 15000.0
    f, err := os.Create("dataset.txt")

    if err != nil {
        fmt.Println("")
    }

    for j:= 0; j < N; j++ {
        f.WriteString(strconv.FormatFloat(toFixed(min + rand.Float64() * (max - min), 2), 'f', -1, 64))
        if j == N - 1 {
            f.WriteString("\n")
            continue
        }
        f.WriteString(", ")
    }

    min = 18.0
    max = 80.0

    for j:= 0; j < N; j++ {
        f.WriteString(strconv.FormatFloat(math.Round(toFixed(min + rand.Float64() * (max - min), 2)), 'f', -1, 64))
        if j == N - 1 {
            f.WriteString("\n")
            continue
        }
        f.WriteString(", ")
    }

    f.Close()
}
