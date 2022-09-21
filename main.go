package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const serverAddr = ":8000"

func createServer(router http.Handler, addr string, timeout time.Duration) *http.Server {
	return &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}
}

func main() {
	register, err := NewRegisterFromTXTDictionary("./palabras.txt")
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/words", register.wordArrayInSpanishHandler)
	r.HandleFunc("/categories/{label}", register.categoriesForWordHandler)

	methodsOk := handlers.AllowedMethods([]string{"GET", "PUT"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	routerWithCors := handlers.CORS(allowedHeaders, methodsOk, originsOk)(r)
	server := createServer(routerWithCors, serverAddr, 10*time.Second)
	log.Fatal(server.ListenAndServe())

}
