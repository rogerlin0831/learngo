package main

import (
	"net/http"

	controller "test/learngo/tryserver/controller"
	db "test/learngo/tryserver/db"
	routes "test/learngo/tryserver/router"
)

func main() {

	controller.InitData()
	router := routes.NewRouter()
	http.ListenAndServe(":3000", router)

	defer db.Close()
}
