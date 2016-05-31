
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
func montePiWorker(cntCh chan monteCount, runtime int) {

	now := time.Now()

	// initialize random number generator
	r := rand.New(rand.NewSource(now.UnixNano()))

	cnt := monteCount{0, 0}
	// loop taking points until time is up
	for time.Since(now) < (time.Duration(runtime) * time.Minute) {
		xr := r.Float64() // 0 < xr < 1
		yr := r.Float64()

		// edist from origin.  In the circle?
		r := math.Sqrt(xr * xr + yr * yr)
		if r <= 1.0 {
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
func montePi(cores int, runtime int) *big.Rat {

	cntCh := make(chan monteCount, cores)

	// launch the workers
	for i := 0 ; i < cores ; i++ {
		go montePiWorker(cntCh, runtime)
	}
	
	// drain the channel
	var inside int64 = 0
	var outside int64 = 0
	res := monteCount{0, 0}
	for i := 0 ; i < cores ; i++ {
		res = <- cntCh
		inside += res.inside
		outside += res.outside
	}

	ratio := big.NewRat(outside, inside)
	one := big.NewRat(1, 1)
	four := big.NewRat(4, 1)
	pi := four.Quo(four, ratio.Add(ratio, one))
	return pi
}

// main function
func main() {

	// command line args
	runtime := flag.Int("runtime", 1, "run time in minutes")
	cores := flag.Int("cores", 2, "number of cores")
	digits := flag.Int("digits", 30, "number of digits to print")
	flag.Parse()

	// compute pi
	pi := montePi(*cores, *runtime)
	
	fmt.Printf("Value of pi using %d cores for %d minutes is %s\n", *cores, *runtime, pi.FloatString(*digits))
}
