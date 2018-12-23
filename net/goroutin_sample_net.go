package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var glNum = 10


func main() {
	wg.Add(glNum)

	flag.Parse()
	url := flag.Arg(0)


	for i:=0; i < glNum; i++ {
		go httpAccess(i, url)
	}

	wg.Wait()
}

func httpAccess(i int, url string){
	defer wg.Done()
	response, error := http.Get(url)

	if error != nil {
		fmt.Println(error)
		return
	}
	fmt.Print(strconv.Itoa(i) + " ")
	fmt.Println(response.Status)

	// body, error
	ioutil.ReadAll(response.Body)

}