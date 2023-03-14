package main

import (
	"fmt"
	"net/http"
	"sample-app/repository"
)

func handler(w http.ResponseWriter, r *http.Request) {
	resp, _ := repository.Get()
	fmt.Fprintf(w, resp)
}

func main() {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(":8080", nil)
}
