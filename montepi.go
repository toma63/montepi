
package main

import ("fmt"
	"flag"
	"math"
	"math/big"
	"math/rand"
	"time"
)

type monteCount struct {
	inside int64
	outside int64
}


// individual worker
// compute pi via Monte Carlo
// count hits inside the circle embedded in a square
// accumulate results for runtime minutes
func MontePiWorker(cntCh chan monteCount, runtime int) {

	now := time.Now()

	// initialize random number generator
	r := rand.New(now.UnixNano())

	cnt := monteCount{0, 0}
	// loop taking points until time is up
	for time.Since(now) < (runtime * time.Minute) {
		xr := r.Float64() // 0 < xr < 1
		yr := r.float64()

		// edist from origin.  In the circle?
		r := math.Sqrt(xr * xr + yr * yr)
		if r < 1.0 {
			cnt.inside++
		} else {
			cnt.outside++
		}
	}
	cntCh <- cnt
}


// comnpute pi using the Monte Carlo method
//   cores specifies the number of parallel workers
//   runtime specifies the run time in minutes
func montePi(cores int, runtime int) big.Rat {

	cntCh := make(chan monteCount, cores)

	// launch the workers
	for i := 0 ; i < cores ; i++ {
		go montePiWorker(cntCh, runtime)
	}
	
	// drain the channel
	inside := new(int64)
	outside := new(int64)
	res := monteCount{0, 0}
	for i := 0 ; i < cores ; i++ {
		res <- cntCh
		inside += res.inside
		outside += res.outside
	}

	res = big.newRat(inside, outside)
	return res
}

// main function
func main() {

	// command line args
	runtime := flag.Int("runtime", 1, "run time in minutes")
	cores := flag.Int("cores", 2, "number of cores")
	flag.Parse()

	// compute pi
	pi := montePi(*runtime, *cores)
	
	fmt.Printf("Value of pin using %d cores for %d minutes is %f\n", *cores, *runtime, pi)
}
