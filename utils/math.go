package utils

import (
	"math"
)

//Sigmoid sigmoid函数
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(float64(-x)))
}

//CombinationString 排列组合
func CombinationString(input []string, n int, cache []string, k int, results *[][]string) {
	//var cache []string = make([]string, count)
	if k == len(cache) {
		//fmt.Printf("[-]%s %s %s\n", cache[0], cache[1], cache[2])
		var result []string
		result = append(result, cache...)
		*results = append(*results, result)
		return
	}
	if len(input)-n < len(cache)-k {
		return
	}

	cache[k] = input[n]
	CombinationString(input, n+1, cache, k+1, results)
	CombinationString(input, n+1, cache, k, results)
}

func PermuteString(input []string) [][]string {
	var ans [][]string
	var dfs func(l []string, temp []string)
	dfs = func(l []string, temp []string) {
		if len(l) == 0 {
			ans = append(ans, temp)
		}
		for i := 0; i < len(l); i++ {
			n := append([]string{}, l...)
			dfs(append(n[:i], n[i+1:]...), append(temp, l[i]))
		}
	}
	dfs(input, []string{})
	return ans
}
