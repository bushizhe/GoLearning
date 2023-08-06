package sorts

func SelectSort(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		min := i // 默认选择第一个元素最小
		for j := i + 1; j < len(arr); j++ {
			if arr[j] < arr[min] {
				arr[j], arr[min] = arr[j], arr[min]
			}
		}
	}
	return arr
}
