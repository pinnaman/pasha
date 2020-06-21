package main

import (
	"fmt"

	"github.com/pinnaman/pasha/internal/gen"
	"github.com/pinnaman/pasha/internal/stats"

	"sync"
)

func printSlice(ftype string, s []int) {
	//fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
	fmt.Printf(ftype+"=%v\n", s)
}

func main() {

	var wg sync.WaitGroup
	var result1 []int
	fmt.Println("")
	/* Create Random array for calcs.****/
	randInt1 := gen.RandomInt(6)
	randInt2 := gen.RandomInt(6)
	randFlt := gen.RandomFlt(6)
	fmt.Println("#************************************#")
	fmt.Println("Random Int array is->", randInt1)
	fmt.Println("Random Float array is->", randFlt)
	fmt.Println(len(randInt1))
	fmt.Println(len(randInt2))
	fmt.Println("")

	fmt.Println("/***** STATS*******/")

	/* Mean Function****/
	chan1 := make(chan float64, 1)
	// Called once...
	wg.Add(1)
	//nl := []float64{98, 93, 77, 82, 83}
	go stats.Mean(chan1, randFlt, &wg)
	wg.Wait()
	//fmt.Println("Mean =>", <-chan1)
	mean := <-chan1
	fmt.Println("Mean =>", mean)

	chan2 := make(chan float64, 1)
	wg.Add(1)
	fmt.Println("Channel 2 made and wg added..")
	// receive on chan 4
	go stats.StdDev(chan2, randFlt, mean, &wg)
	//fmt.Println("STD DEV =>", <-stdDev(nl, <-chan3))
	wg.Wait()
	//fmt.Println("STD DEV=>", <-chan2)
	stddev := <-chan2
	fmt.Println("STD DEV=>", stddev)

	// Fibonacci
	chan3 := make(chan int, 6)
	wg.Add(1)
	go stats.Fibonacci(cap(chan3), chan3, &wg)
	wg.Wait()
	for i := range chan3 {
		//fmt.Println(i)
		result1 = append(result1, i)
	}
	printSlice("Fibonacci", result1)

	//printSlice("Product=>", <-chan1)
	//fmt.Println("Mean =>", mean)

	chan4 := make(chan float64, 1)
	wg.Add(1)
	//fmt.Println("Channel 4 made and wg added..")
	fmt.Println("Array1 is->", randInt1, "Array2 is->", randInt2)
	go stats.CorrCoefficient(chan4, randInt1, randInt2, len(randInt1), &wg)
	wg.Wait()
	fmt.Println("Correlation Coeff=>", <-chan4)

	close(chan1)
	close(chan2)
	close(chan4)

	fmt.Println("/***** END STATS*******/")

}
