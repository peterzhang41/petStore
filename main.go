package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	op "github.com/peterzhang41/petStore/operations"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	s := router.PathPrefix("/v2").Subrouter()
	s.HandleFunc("/pet/findByTags", op.FindPetsByTags).Queries("tags", "{tags}")
	s.HandleFunc("/pet/findByStatus", op.FindPetsByStatus).Queries("status", "{status}")
	s.HandleFunc("/pet", op.AddPet).Methods("POST")
	s.HandleFunc("/pet", op.UpdatePet).Methods("PUT")
	s.HandleFunc("/pet/{petId}", op.GetPetById).Methods("GET")
	s.HandleFunc("/pet/{petId}", op.UpdatePetWithForm).Methods("POST")
	s.HandleFunc("/pet/{petId}", op.DeletePet).Methods("DELETE")
	s.HandleFunc("/pet/{petId}/uploadImage", op.UploadFile).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}