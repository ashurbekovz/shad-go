package main

import (
	"net/http"
	_ "net/http/pprof"
)

func main() {
	for i := range 10000 {
		go func() {
			select {}
		}()
	}

	http.ListenAndServe(":8080", nil)
}
