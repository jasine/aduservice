package main

import (
	"fmt"
	"github.com/deepglint/aduservice/sdk"
	"net/http"
)

var (
	adu = sdk.NewAdu4vulcand()
)

func main() {
	stop := make(chan bool)
	adu.Start()
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", Auth)

	http.ListenAndServe(":1111", mux)
	<-stop
}
func Auth(w http.ResponseWriter, r *http.Request) {
	pass := adu.Check(r)
	if pass {
		fmt.Fprint(w, sdk.TRUE_BODY)
	} else {
		fmt.Fprint(w, sdk.FALSE_BODY)
	}
}
