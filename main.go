package main

func main() {
	defer initUI()
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	colors := []string{"", green, red, green, red, green, red, green, red, green, red}
	printGraph(arrGraph(arr, colors))
}
