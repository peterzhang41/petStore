package operations

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/peterzhang41/petStore/models"
)

// For Testing Only
var pets = make(map[int64]*models.Pet)
var lastID int64
var petsLock = &sync.Mutex{}

// TODO: Database construction and CURD implementation. The init.sql file has been written in the db folder

func newPetID() int64 {
	return atomic.AddInt64(&lastID, 1)
}


func UpdatePet(w http.ResponseWriter, r *http.Request) {
	var newPet models.Pet
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid ID supplied")
		return
	}

	err = json.Unmarshal(reqBody, &newPet)
	if err != nil {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Validation exception")
		return
	}

	if newPet.Id == 0 {
		newPet.Id = lastID
	}else {
		_, exists := pets[newPet.Id]
		if !exists {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Pet not found")
			return
		}
	}

	if newPet.Name == ""|| len(newPet.PhotoUrls) == 0{
		w.WriteHeader(405)
		fmt.Fprintf(w, "Validation exception")
		return
	}

	petsLock.Lock()
	defer petsLock.Unlock()

	pets[newPet.Id] = &newPet

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pets[newPet.Id])

}

func AddPet(w http.ResponseWriter, r *http.Request) {
	var newPet models.Pet
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Invalid input\n")
		return
	}

	json.Unmarshal(reqBody, &newPet)

	if newPet.Name == ""|| len(newPet.PhotoUrls) == 0{
		w.WriteHeader(405)
		fmt.Fprintf(w, "Invalid input\n")
		return
	}

	newID := newPetID()
	newPet.Id= newID
	pets[newID] = &newPet

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(newPet)
}

func GetPetById(w http.ResponseWriter, r *http.Request) {
	petId := mux.Vars(r)["petId"]
	id, err := strconv.ParseInt(petId, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid ID supplied\n")
		return
	}

	_, exists := pets[id]
	if !exists {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Pet not found\n")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pets[id])

}

func UpdatePetWithForm (w http.ResponseWriter, r *http.Request) {
	petId := mux.Vars(r)["petId"]

	id, err := strconv.ParseInt(petId, 10, 64)
	if err != nil {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Invalid input\n")
		return
	}
	_, exists := pets[id]
	if !exists {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Invalid input\n")
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Invalid input\n")
		return
	}

	pets[id].Name = r.FormValue("name")
	pets[id].Status = r.FormValue("status")
	w.WriteHeader(http.StatusOK)
}

func DeletePet(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("api_key")
	//TODO: validation required in prod
	fmt.Println("api_key: ", apiKey)

	petId := mux.Vars(r)["petId"]

	id, err := strconv.ParseInt(petId, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid input supplied\n")
		return
	}

	petsLock.Lock()
	defer petsLock.Unlock()
	_, exists := pets[id]
	if !exists {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Pet not found\n")
		return
	}
	delete(pets,id)
	w.WriteHeader(http.StatusOK)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	petId := mux.Vars(r)["petId"]
	id, err := strconv.ParseInt(petId, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid input supplied\n")
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	path := "./"+ header.Filename
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)


	var apiResponse models.ApiResponse
	apiResponse.Code = 200
	apiResponse.Type = ""
	apiResponse.Message = "additionalMetadata: " + r.FormValue("additionalMetadata") +"\n"+ "File uploaded to " + path + "," + strconv.FormatInt(header.Size,10) + " bytes"

	_, exists := pets[id]
	if !exists {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Pet not found\n")
		return
	}
	pets[id].PhotoUrls = append(pets[id].PhotoUrls,path)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse)

}

//Deprecated
func FindPetsByTags (w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query()["tags"]
	tags = removeDuplicates(tags)
	fmt.Println(tags)

	tagsFilter := make(map[string]bool)
	for _, v := range tags {
		tagsFilter[v] = true
	}
	var result []*models.Pet
	for _, pet := range pets {
		for _, tag := range pet.Tags{
			if tagsFilter[tag.Name] ==true {
				result = append(result, pet)
				continue
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}


func FindPetsByStatus (w http.ResponseWriter, r *http.Request) {
	statusCode := map[string]bool{"available": true, "pending": true, "sold":true}
	status := r.URL.Query()["status"]

	//TODO: O(2n) could be refactored to O(1n)
	status = removeDuplicates(status)
	for _, v := range status {
		if statusCode[v] != true {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid status value\n")
			return
		}
	}
	fmt.Println(status)

	var result []*models.Pet
	for _, v := range pets {
		if statusCode[v.Status] == true {
			result = append(result, v)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}