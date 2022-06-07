package util

func ValueInArrayLarger(arr []float64, value float64) bool {
	for _, val := range arr {
		if val > value {
			return true
		}
	}
	return false
}

func ValueInArraySmaller(arr []float64, value float64) bool {
	for _, val := range arr {
		if val < value {
			return true
		}
	}
	return false
}
