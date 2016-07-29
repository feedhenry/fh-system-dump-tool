package main

import (
	"fmt"
	"io/ioutil"
)

func helpTask() int {
	help, err := ioutil.ReadFile("./help.txt")
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Print(string(help))
	return 0
}
