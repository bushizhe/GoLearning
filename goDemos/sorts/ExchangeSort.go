package sorts

func ExchangeSort(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[j] < arr[i] {
				arr[j+1], arr[i] = arr[i], arr[j+1]
			}
		}
	}
	return arr
}
