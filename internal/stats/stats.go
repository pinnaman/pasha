package stats

import (
	"math"
	"sync"

	"github.com/pinnaman/gostats/gen"
)

var wg sync.WaitGroup
var mu sync.Mutex
var sum = 0

// LaunchStats need to implement for testing?????
func LaunchStats(n int) []float64 {
	//randFlt := gen.RandomFlt(n)
	//wg.Add(1)
	//	return <-sum(Mean(c,randFlt)
	return gen.RandomFlt(n)
}

func Prod(c chan int, num1 int, num2 int, wg *sync.WaitGroup) {
	defer wg.Done()
	//c <- someValue * 5
	c <- num1 * num2
}

func Fibonacci(n int, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	x, y := 0, 1
	for i := 0; i < n; i++ {

		c <- x
		x, y = y, x+y
	}
	close(c)
}

func Mean(c chan float64, data []float64, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0.0
	//fmt.Println(len(data))
	for _, d := range data {
		sum += float64(d)
	}
	c <- sum / float64(len(data))
}

func StdDev(c chan float64, data []float64, mean float64, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0.0
	//fmt.Println(len(data))
	for _, d := range data {
		sum += math.Pow(float64(d)-mean, 2)
	}
	//fmt.Println(sum)
	variance := sum / float64(len(data)-1)
	c <- math.Sqrt(variance)
}

func CorrCoefficient(c chan float64, X []int, Y []int, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum_X := 0
	sum_Y := 0
	sum_XY := 0
	squareSum_X := 0
	squareSum_Y := 0

	for i := 0; i < n; i++ {
		// sum of elements of array X.
		sum_X = sum_X + X[i]

		// sum of elements of array Y.
		sum_Y = sum_Y + Y[i]

		// sum of X[i] * Y[i].
		sum_XY = sum_XY + X[i]*Y[i]

		// sum of square of array elements.
		squareSum_X = squareSum_X + X[i]*X[i]
		squareSum_Y = squareSum_Y + Y[i]*Y[i]
	}

	// use formula for calculating correlation
	// coefficient.
	corr := float64((n*sum_XY - sum_X*sum_Y)) /
		(math.Sqrt(float64((n*squareSum_X - sum_X*sum_X) * (n*squareSum_Y - sum_Y*sum_Y))))

	//return corr
	c <- corr
}

// Add WVG, MWAVG, EMWA

// Linearly Weighted Moving Average For last 5 Days Weights=>(5/15,4/15..)
//((D5*5/15)+(D4*4/15)+(D3*3/15)+(D2*2/15)+(D1*1/15)) / (5+4+3+2+1)
