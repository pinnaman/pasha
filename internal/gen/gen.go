package gen

import (
	"math/rand"
	"time"
)

func RandomInt(ln int) []int {
	var numbers []int
	var array [40]int
	rand.Seed(time.Now().UTC().UnixNano())
	for i := ln - 1; i > 0; i-- {
		//j := rand.Intn(i + 1)
		n := rand.Int() % len(array)
		//array[i], array[j] = array[j], array[i]
		numbers = append(numbers, n)
	}
	return numbers
}

// RandomFlt returns random Floats
func RandomFlt(ln int) []float64 {
	var numbers []float64
	//var array [40]float64
	rand.Seed(time.Now().UTC().UnixNano())
	for i := ln - 1; i > 0; i-- {
		//j := rand.Intn(i + 1)
		n := rand.Float64() * 90 // simulate number close to 90
		//array[i], array[j] = array[j], array[i]
		numbers = append(numbers, n)
	}
	return numbers
}

// Generator Integer channel
func Generator(max int) <-chan int {
	outChInt := make(chan int, 100)
	go func() {
		for i := 1; i <= max; i++ {
			outChInt <- i
		}
		close(outChInt)
	}()
	return outChInt
}
