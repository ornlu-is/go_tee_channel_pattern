package main

import "fmt"

func numberStream() <-chan float64 {
	ch := make(chan float64)
	numberStrings := []float64{1., 2., 3., 4., 5., 6., 7., 8., 9., 10.}

	go func() {
		for _, numberString := range numberStrings {
			ch <- numberString
		}

		close(ch)
		return
	}()

	return ch
}

func teeChannel(c <-chan float64) (<-chan float64, <-chan float64) {
	tee1 := make(chan float64)
	tee2 := make(chan float64)

	go func() {
		defer func() {
			close(tee1)
			close(tee2)
		}()

		for val := range c {
			for i := 0; i < 2; i++ {
				var tee1, tee2 = tee1, tee2
				select {
				case tee1 <- val:
					tee1 = nil
				case tee2 <- val:
					tee2 = nil
				}
			}
		}

		return
	}()

	return tee1, tee2
}

func main() {
	done := make(chan struct{})
	defer close(done)

	dataStream := numberStream()

	teedStream1, teedStream2 := teeChannel(dataStream)

	for val1 := range teedStream1 {
		fmt.Printf("tee1: %f\n", val1)
		fmt.Printf("tee2: %f\n", <-teedStream2)
	}
}
