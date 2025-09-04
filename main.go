package main

func main() {
	defer initUI()()
	rawArr := getSequenceArr(1, 40)
	shuffleArr(rawArr)
	arr := newArrObj(rawArr)
	colors := make([]string, arr.len())
	render(arrGraph(rawArr, colors))
	quickSort(arr, 0, arr.len()-1)
}
