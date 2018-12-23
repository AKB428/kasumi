package main

import (
	"fmt"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var glNum = 100


func main() {
	wg.Add(glNum)

	for i:=0; i < glNum; i++ {
		go sayHello(i)
	}

	wg.Wait()
}

func sayHello(i int){
	defer wg.Done()
	fmt.Print( strconv.Itoa(i) + " HELOOOO!\n")
}
