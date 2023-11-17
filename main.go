package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	Serve("localhost:8088")
}
