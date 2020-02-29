package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/peterzhang41/petStore/models"
	"log"
	"net/http"

	op "github.com/peterzhang41/petStore/operations"
)

func main() {

	models.MockData()
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

	fmt.Println("server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
