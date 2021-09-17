package utils

import "math"

//Sigmoid sigmoid函数
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(float64(-x)))
}
