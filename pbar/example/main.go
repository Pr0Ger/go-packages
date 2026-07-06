package main

import (
	"fmt"
	"time"

	"go.pr0ger.dev/x/pbar"
)

func main() {
	bar := pbar.ProgressBar{}
	bar.Bars = []pbar.Bar{
		{Current: 0, Total: 100},
	}

	_ = bar.Start()
	for i := uint(0); i < 100; i++ {
		bar.SetBarValue(0, i)
		time.Sleep(100 * time.Millisecond)

		if i%10 == 0 {
			fmt.Println("stupid log message")

			time.Sleep(300 * time.Millisecond)
		}
	}
}
