package main

// type arrObj interface {
// 	get(ind int) int
// 	set(ind, val int)
// 	swap(i, j int)
// 	len() int
// }

func bubbleSort(arr arrObj) {
	n := arr.len()
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr.get(j) > arr.get(j+1) {
				arr.swap(j, j+1)
			}
		}
	}
}
