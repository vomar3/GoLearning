package main

func main() {

}

func SumSlices(firstSlice, secondSlice []int) []int {
	length := min(len(firstSlice), len(secondSlice))
	answer := make([]int, 0, length)

	for i := 0; i < length; i++ {
		answer[i] = firstSlice[i] + secondSlice[i]
	}

	return answer
}
