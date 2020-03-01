package operations

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/peterzhang41/petStore/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

func UpdatePet(w http.ResponseWriter, r *http.Request) {
	newPet, err := readBodyAndUnmarshalPet(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid ID supplied\n", 400)
		return
	}

	pet, appErr := models.UpdatePet(newPet)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	respondWithJSON(w, 200, pet)
}

func AddPet(w http.ResponseWriter, r *http.Request) {
	newPet, err := readBodyAndUnmarshalPet(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input\n", 405)
		return
	}

	pet, appErr := models.AddPet(newPet)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	respondWithJSON(w, 200, pet)
}

func GetPetById(w http.ResponseWriter, r *http.Request) {
	id, err := readPetIdFromRequest(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid ID supplied", 400)
		return
	}

	pet, appErr := models.GetPetById(id)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	respondWithJSON(w, 200, pet)
}

func UpdatePetWithForm(w http.ResponseWriter, r *http.Request) {
	id, err := readPetIdFromRequest(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input\n", 405)
		return
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input\n", 405)
		return
	}
	name := r.FormValue("name")
	status := r.FormValue("status")

	appErr := models.UpdatePetWithForm(id, name, status)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	respondWithJSON(w, 200, nil)
}

func DeletePet(w http.ResponseWriter, r *http.Request) {
	//apiKey := r.Header.Get("api_key")
	//TODO: validation required
	//fmt.Println("api_key: ", apiKey)

	id, err := readPetIdFromRequest(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input supplied", 400)
		return
	}
	appErr := models.DeletePet(id)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	respondWithJSON(w, 200, nil)

}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	id, err := readPetIdFromRequest(r)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input supplied\n", 400)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input supplied\n", 400)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid input supplied\n", 400)
		return
	}
	metadata := r.FormValue("additionalMetadata")

	apiResponse, appErr := models.UploadFile(id, file, header, metadata)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	respondWithJSON(w, 200, apiResponse)

}

//Deprecated
func FindPetsByTags(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query()["tags"]

	result, appErr := models.FindPetsByTags(tags)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	respondWithJSON(w, 200, result)
}

func FindPetsByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query()["status"]
	result, appErr := models.FindPetsByStatus(status)
	if appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	respondWithJSON(w, 200, result)
}

func readBodyAndUnmarshalPet(r *http.Request) (*models.Pet, error) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var newPet models.Pet
	err = json.Unmarshal(reqBody, &newPet)
	if err != nil {
		return nil, err
	}

	return &newPet, nil
}

func readPetIdFromRequest(r *http.Request) (int64, error) {
	petId := mux.Vars(r)["petId"]
	id, err := strconv.ParseInt(petId, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
