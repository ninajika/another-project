package main

import "fmt"

func bubbleSort(arr []int) {
	n := len(arr)
	swapped := true
	for i := 0; i < n-1 && swapped; i++ {
		swapped = false
		for j := 0; j < n-1-i; j++ {
			if arr[j] > arr[j+1] {
				swapped = true
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func main() {
	a := []int{6, 16, 4, 10, 17}
	bubbleSort(a)
	fmt.Println(a)
}
