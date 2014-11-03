package main

import (
	"fmt"
	"github.com/slimsag/swigcflags"
)

func main() {
	f := swigcflags.XOpenDisplay
	_ = f
	fmt.Println(f)
}
