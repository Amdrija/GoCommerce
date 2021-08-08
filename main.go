package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Amdrija/GoCommerce/router"
)

func do(w http.ResponseWriter, r *http.Request, routeParameters router.RouteParameters) {
	fmt.Fprint(w, r.URL.Path)
}

func main() {
	d := router.NewRouter()
	d.Get("/about/{id}", do)
	d.Get("/about/{id}/asdf", do)
	d.Get("/view/asdf", do)
	d.Post("/view", do)
	d.Get("/", do)

	http.HandleFunc("/", d.Dispatch)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
