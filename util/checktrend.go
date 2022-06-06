package util

func IsTrendGoingUp(arr []float64) bool {
	for i := 0; i < len(arr)-2; i++ {
		if arr[i] > arr[i+1] {
			return false
		}
	}
	return true
}

func IsTrendGoingDown(arr []float64) bool {
	for i := 0; i < len(arr)-2; i++ {
		if arr[i] < arr[i+1] {
			return false
		}
	}
	return true
}
