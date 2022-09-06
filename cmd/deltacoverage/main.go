package main

import (
	"fmt"
	"runtime"
)

func retrieveCallInfo() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func main() {
	fmt.Println(retrieveCallInfo())
}
