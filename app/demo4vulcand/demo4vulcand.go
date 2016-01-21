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
	adu.Start()
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", Auth)

	if err := http.ListenAndServe(":1111", mux); err != nil {
		fmt.Println(err.Error())
	}
}
func Auth(w http.ResponseWriter, r *http.Request) {
	pass := adu.Check(r)
	if pass {
		fmt.Fprint(w, sdk.TRUE_BODY)
	} else {
		fmt.Fprint(w, sdk.FALSE_BODY)
	}
}
