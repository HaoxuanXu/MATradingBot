package util

func ResizeFloatArray(arr []float64, arrayLength int) []float64 {
	if len(arr) <= arrayLength {
		return arr
	}
	cutoff := len(arr) - arrayLength
	return arr[cutoff:]
}
