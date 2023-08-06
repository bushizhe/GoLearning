package sorts

func BubbleSort(arr []int) []int {
	swapped := true
	for swapped {
		swapped = false // 如果arr有序，只需循环一次即可
		for i := 0; i < len(arr)-1; i++ {
			if arr[i+1] < arr[i] {
				arr[i+1], arr[i] = arr[i], arr[i+1]
				swapped = true
			}
		}
	}
	return arr
}
