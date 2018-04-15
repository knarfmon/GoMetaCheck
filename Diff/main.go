package main

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = "Frank is great"
	text2 = "Frank is Great"
)

func main() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

	fmt.Println(diffs)
fmt.Println("---------------------")
	fmt.Println(dmp.DiffPrettyText(diffs))
}
