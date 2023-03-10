package chromedp

import (
	"math"
	"strings"
)

func QueryMode(sel string) string {
	if strings.HasPrefix(sel, "//") {
		return "xpath"
	} else {
		return "queryAll"
	}
}

func calcSpeeds(distance float64, times int) (total float64, points []float64) {
	/*x平方的曲线，计算每步的步长*/
	totalX := float64(0)
	steps := make([]float64, 0)
	for i := 1; i < times; i++ {
		v := math.Pow(float64(i), 2)
		if i == 1 {
			steps = append([]float64{v}, steps...)
		} else {
			v = v - steps[i-2]
			steps = append([]float64{v}, steps...)
		}
		totalX += v
	}
	/*根据每步的步长，归一化distance对应的步长*/
	points = make([]float64, 0)
	total = 0.0
	for i := 0; i < times-1; i++ {
		v := distance * (steps[i] / totalX)
		points = append(points, v)
		total += v
	}
	return total, points
}

func BuildTracks(distance float64) []float64 {
	times := 33
	_, points := calcSpeeds(distance, times)
	//先加20，然后在退回来！
	//backDistance := total_ - distance
	//_, backPoints := calcSpeeds(backDistance, 10)
	//for _, p := range backPoints {
	//	points = append(points, -p)
	//}
	return points
}
