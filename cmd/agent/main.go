package main

import "net/http"

func mainPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Привет!"))
}

func apiPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Это страница /api."))
}

func main() {
	http.HandleFunc(`/api`, apiPage)
	http.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
