package main

import (
	"fmt"

	"example.com/price-calculator/filemanager"
	"example.com/price-calculator/prices"
)

func main() {
	taxRates := []float64{0, 0.07, 0.1, 0.15}
	doneCh := make([]chan bool, len(taxRates))
	errorCh := make([]chan error, len(taxRates))

	for index, taxRate := range taxRates {
		doneCh[index] = make(chan bool)
		errorCh[index] = make(chan error)
		fm := filemanager.New("prices.txt", fmt.Sprintf("result_%.0f.json", taxRate*100))
		//cmdm := cmdmanager.New()
		priceJob := prices.NewTaxIncludedPriceJob(fm, taxRate)
		go priceJob.Process(doneCh[index], errorCh[index])

		// if err != nil {
		// 	fmt.Println("could not process job")
		// 	fmt.Println(err)
		// }
	}

	for index := range taxRates {
		select {
		case err := <-errorCh[index]:
			if err != nil {
				fmt.Println(err)
			}
		case <-doneCh[index]:
			fmt.Println("Done")
		}
	}
}
