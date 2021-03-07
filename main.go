package main

import (
	"fmt"
	"io"
	"net/http"
	// m "github.com/comerc/segezha4/mymodule"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!!!")
}

func main() {
	port := "8080" // os.Getenv("PORT")
	http.HandleFunc("/", hello)
	fmt.Println("Server listening!")
	// m.Yo()
	http.ListenAndServe(":"+port, nil)
}
