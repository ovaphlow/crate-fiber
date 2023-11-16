package main

import (
	"ovaphlow/crate/internal/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Serve("localhost:8088")
}
