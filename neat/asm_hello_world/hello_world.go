package main

func add(a, b int32) int32

func main() {
	if add(1, 4) != 5 {
		panic("Incorrect result!")
	}
}
