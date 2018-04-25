package main

import (
	"strings"
	"fmt"
)

func main(){

	str := "(-)He(0)llo, (+)nice(0)play(-)g(0)round"

	deletionSign := "(-)"
	additionSign := "(+)"
	normalSign := "(0)"

	var norm,backColor string
	var idx int

	for {
		idxAdd := strings.Index(str, additionSign)
		fmt.Println("idxAdd-",idxAdd)
		if idxAdd == -1{idxAdd = 1000}
		idxDel := strings.Index(str, deletionSign)
		if idxDel == -1{idxDel = 1000}
		fmt.Println("idxDel-",idxDel)
		idxNil := strings.Index(str, normalSign)
		if idxNil == -1{idxNil = 1000}
		fmt.Println("idxNil-",idxNil)

		//no sighn left
		if idxAdd+idxDel+idxNil == 3000 {
			fmt.Println("nothing");
			return
		}
		//finding the closest sighn to get index and type

		if (idxAdd < idxDel && idxAdd < idxNil) {
			idx = idxAdd;
			backColor = "green"
		} else if (idxDel < idxNil && idxDel < idxAdd) {
			idx = idxDel;
			backColor = "red"
		} else {
			idx = idxNil;
			backColor = "white"
		}

		fmt.Println("after if-",idx)

		norm = str[0:idx]
		fmt.Println(norm)
		fmt.Println("Change background color-", backColor)
		idx = idx + 3
		str = str[idx:]
		fmt.Println("Remaining string-", str)


	}












}

