package models

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"sync/atomic"
)

type Pet struct {
	Id int64 `json:"id,omitempty"`

	Category *Category `json:"category,omitempty"`

	Name string `json:"name"`

	PhotoUrls []string `json:"photoUrls"`

	Tags []Tag `json:"tags,omitempty"`

	// pet status in the uploadedFiles
	Status string `json:"status,omitempty"`
}

type appError struct {
	Error   error
	Code    int
	Message string
}

// For Testing Only
var pets = make(map[int64]*Pet)
var lastID int64

func newPetID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

// Mock Data
func MockData() {
	pets[1] = &Pet{Id: 1, Name: "test1", Status: "available"}
	pets[2] = &Pet{Id: 2, Name: "test2", Status: "pending"}
	pets[3] = &Pet{Id: 3, Name: "test3", Status: "sold"}
	lastID = 3
}

// TODO: Database construction and CURD implementation. The init.sql file has been written in the db folder

func UpdatePet(newPet *Pet) (*Pet, *appError) {
	if newPet.Id == 0 {
		newPet.Id = lastID
	} else {
		_, exists := pets[newPet.Id]
		if !exists {
			return nil, &appError{nil, 404, "Pet not found"}
		}
	}

	if newPet.Name == "" || len(newPet.PhotoUrls) == 0 {
		return nil, &appError{nil, 405, "Validation exception"}
	}

	pets[newPet.Id] = newPet

	return pets[newPet.Id], nil

}

func AddPet(newPet *Pet) (*Pet, *appError) {
	if newPet.Name == "" || len(newPet.PhotoUrls) == 0 {
		return nil, &appError{nil, 405, "Invalid input\n"}
	}

	newID := newPetID()
	newPet.Id = newID
	pets[newID] = newPet

	return pets[newID], nil
}

func GetPetById(id int64) (*Pet, *appError) {
	_, exists := pets[id]
	if !exists {
		return nil, &appError{nil, 404, "Pet not found\n"}
	}
	return pets[id], nil
}

func UpdatePetWithForm(id int64, name string, status string) *appError {
	_, exists := pets[id]
	if !exists {
		return &appError{nil, 405, "Invalid input\n"}
	}

	pets[id].Name = name
	pets[id].Status = status

	return nil
}

func DeletePet(id int64) *appError {
	_, exists := pets[id]
	if !exists {
		return &appError{nil, 404, "Pet not found\n\n"}
	}
	delete(pets, id)
	return nil
}

func UploadFile(id int64, file multipart.File, header *multipart.FileHeader, metadata string) (*ApiResponse, *appError) {
	_, exists := pets[id]
	if !exists {
		return nil, &appError{nil, 400, "Invalid input supplied\n"}
	}

	defer file.Close()
	path := "./" + header.Filename
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, &appError{err, 404, "Pet not found\n\n"}
	}
	defer f.Close()
	io.Copy(f, file)

	var apiResponse ApiResponse
	apiResponse.Code = 200
	apiResponse.Type = ""
	apiResponse.Message = "additionalMetadata: " + metadata + "\n" + "File uploaded to " + path + "," + strconv.FormatInt(header.Size, 10) + " bytes"

	pets[id].PhotoUrls = append(pets[id].PhotoUrls, path)

	return &apiResponse, nil

}

//Deprecated
func FindPetsByTags(tags []string) (result []*Pet, error *appError) {
	tags = removeDuplicates(tags)
	fmt.Println(tags)

	tagsFilter := make(map[string]bool)
	for _, v := range tags {
		tagsFilter[v] = true
	}

	for _, pet := range pets {
		for _, tag := range pet.Tags {
			if tagsFilter[tag.Name] == true {
				result = append(result, pet)
				continue
			}
		}
	}

	if len(result) == 0 {
		return nil, &appError{nil, 400, "Invalid tag value\n"}
	}

	return result, nil

}

func FindPetsByStatus(status []string) (result []*Pet, error *appError) {
	statusCode := map[string]bool{"available": true, "pending": true, "sold": true}

	//TODO: O(2n) could be refactored to O(1n)
	status = removeDuplicates(status)
	for _, v := range status {
		if statusCode[v] != true {
			return nil, &appError{nil, 400, "Invalid status value\n"}
		}
	}
	fmt.Println(status)

	for _, v := range pets {
		if statusCode[v.Status] == true {
			result = append(result, v)
		}
	}

	return result, nil

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
