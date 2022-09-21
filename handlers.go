package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (register *SpanishRegister) wordArrayInSpanishHandler(w http.ResponseWriter, r *http.Request) {
	var words []string
	err := DecodeJSONBody(w, r, &words)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
			return
		}
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(register.wordsInSpanishFilter(words))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Println("New Request", r.Body)
	//allow cors
}

func (register *SpanishRegister) categoriesForWordHandler(w http.ResponseWriter, r *http.Request) {
	var categories []string
	vars := mux.Vars(r)
	word := vars["label"]
	if register.wordInSpanish(word) {
		categories = append(categories, "Español")
	}
	err := json.NewEncoder(w).Encode(categories)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
