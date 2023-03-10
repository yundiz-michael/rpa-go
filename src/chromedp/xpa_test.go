package chromedp

import (
	"fmt"
	"testing"
)

func TestTracks(t *testing.T) {
	points := BuildTracks(float64(80))
	total := 0.0
	for i := 0; i < len(points); i++ {
		fmt.Printf("index= %d value=%f \n", i, points[i])
		total += points[i]
	}
	fmt.Printf("%f", total)
}
