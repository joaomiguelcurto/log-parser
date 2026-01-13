package utils

import (
	"math"
	"strconv"
)

func RoundFloat(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

func FormatNumberSimple(number float64) string {
	if number >= 1000000 {
		return strconv.FormatFloat(number/1000000, 'f', 2, 64) + "M"
	} else if number >= 1000 {
		return strconv.FormatFloat(number/1000, 'f', 2, 64) + "K"
	}

	return strconv.FormatFloat(number, 'f', 0, 64)
}
