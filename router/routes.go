package routes

import (
	"net/http"
	controller "test/learngo/tryserver/controller"

	"github.com/gorilla/mux"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("POST", "/api/addUser", controller.AddData, nil)
	register("POST", "/api/deleteUser", controller.DeleteData, nil)
	register("POST", "/api/UpdateUser", controller.UpdateData, nil)
	register("POST", "/api/login", controller.Login, nil)
	register("GET", "/api/GetPassword/{user}", controller.GetPassword, nil)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	for _, rotue := range routes {
		r.Methods(rotue.Method).Path(rotue.Pattern).Handler(rotue.Handler)
		if rotue.Middleware != nil {
			r.Use(rotue.Middleware)
		}
	}

	return r
}

func register(method string, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}
