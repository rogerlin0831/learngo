package main

import (
	"net/http"

	routes "test/learngo/tryserver/router"
)

func main() {

	router := routes.NewRouter()
	http.ListenAndServe(":3000", router)

}
