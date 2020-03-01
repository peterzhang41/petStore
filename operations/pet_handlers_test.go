package operations

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/peterzhang41/petStore/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetPetByIdSuccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/pet/1", nil)
	expected := `{"id":1,"name":"test1","photoUrls":null,"status":"available"}`
	testTemplate(t, req, "/v2/pet/{petId}", GetPetById, 200, expected)
}

func TestGetPetByIdInvalidInput(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/pet/TWO", nil)
	testTemplate(t, req, "/v2/pet/{petId}", GetPetById, 400, "Invalid ID supplied")
}

func TestGetPetByIdPetNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/pet/99", nil)
	testTemplate(t, req, "/v2/pet/{petId}", GetPetById, 404, "Pet not found")
}

func TestDeletePetSuccess(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/v2/pet/1", nil)
	testTemplate(t, req, "/v2/pet/{petId}", DeletePet, 200, "null")
}

func TestDeletePetInvalidInput(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/v2/pet/TWO", nil)
	testTemplate(t, req, "/v2/pet/{petId}", DeletePet, 400, "Invalid input supplied")

}

func TestDeletePetPetNotFound(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/v2/pet/99", nil)
	testTemplate(t, req, "/v2/pet/{petId}", DeletePet, 404, "Pet not found")
}

func TestAddPetSuccess(t *testing.T) {
	body := `{"id":0,"name":"test1","photoUrls":["test.log"],"status":"available"}`
	req := httptest.NewRequest("POST", "/v2/pet", bytes.NewBufferString(body))
	expect := `{"id":4,"name":"test1","photoUrls":["test.log"],"status":"available"}`
	testTemplate(t, req, "/v2/pet", AddPet, 200, expect)
}

func TestAddPetInvalidInput(t *testing.T) {
	body := `{"id":1,"name":"test1","photoUrls":null,"status":"available"}`
	req := httptest.NewRequest("POST", "/v2/pet", bytes.NewBufferString(body))
	testTemplate(t, req, "/v2/pet", AddPet, 405, "Invalid input")

	body = `{"id":1,"name":"","photoUrls":["photoPath"],"status":"available"}`
	req = httptest.NewRequest("POST", "/v2/pet", bytes.NewBufferString(body))
	testTemplate(t, req, "/v2/pet", AddPet, 405, "Invalid input")
}

func TestUpdatePetSuccess(t *testing.T) {
	body := `{"id":1,"name":"test2","photoUrls":["photoOurPath"],"status":"available"}`
	req := httptest.NewRequest("PUT", "/v2/pet", bytes.NewBufferString(body))
	testTemplate(t, req, "/v2/pet", UpdatePet, 200, body)
}

func TestUpdatePetInvalidInput(t *testing.T) {
	body := `"invalidInput":"testing"`
	req := httptest.NewRequest("PUT", "/v2/pet", bytes.NewBufferString(body))
	testTemplate(t, req, "/v2/pet", UpdatePet, 400, "Invalid ID supplied")
}

func TestUpdatePetPetNotFound(t *testing.T) {
	body := `{"id":999,"name":"test2","photoUrls":["photoOurPath"],"status":"available"}`
	req := httptest.NewRequest("PUT", "/v2/pet", bytes.NewBufferString(body))
	testTemplate(t, req, "/v2/pet", UpdatePet, 404, "Pet not found")
}

func TestUpdatePetValidationException(t *testing.T) {
	body := `{"id":1,"name":"","photoUrls":["photoOurPath"],"status":"available"}`
	req := httptest.NewRequest("PUT", "/v2/pet", bytes.NewBufferString(body))
	testTemplate(t, req, "/v2/pet", UpdatePet, 405, "Validation exception")

	body = `{"id":1,"name":"test1","photoUrls":[],"status":"available"}`
	req = httptest.NewRequest("PUT", "/v2/pet", bytes.NewBufferString(body))
	testTemplate(t, req, "/v2/pet", UpdatePet, 405, "Validation exception")
}

func TestUpdatePetWithFormSuccess(t *testing.T) {
	form := url.Values{}
	form.Add("name", "testName")
	form.Add("tags", "testTags")
	req, err := http.NewRequest("POST", "/v2/pet/1", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	testTemplate(t, req, "/v2/pet/{petId}", UpdatePetWithForm, 200, "null")
}

func TestUpdatePetWithFormInvalidInput(t *testing.T) {
	form := url.Values{}
	form.Add("name", "testName")
	form.Add("tags", "testTags")

	//id = TWO
	req, err := http.NewRequest("POST", "/v2/pet/TWO", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	testTemplate(t, req, "/v2/pet/{petId}", UpdatePetWithForm, 405, "Invalid input")

	//body = nil
	req, err = http.NewRequest("POST", "/v2/pet/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	testTemplate(t, req, "/v2/pet/{petId}", UpdatePetWithForm, 405, "Invalid input")

	//id not found
	req, err = http.NewRequest("POST", "/v2/pet/9999", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	testTemplate(t, req, "/v2/pet/{petId}", UpdatePetWithForm, 405, "Invalid input")
}

func TestUpdatePetUploadFileSuccess(t *testing.T) {
	path := "../test/upload_test_file.log"
	body, writer := readFile(t, path, "file")
	req := httptest.NewRequest("POST", "/v2/pet/1/uploadImage", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	expect := `{"code":200,"message":"additionalMetadata: \nFile uploaded to ./upload_test_file.log,928 bytes"}`
	testTemplate(t, req, "/v2/pet/{petId}/uploadImage", UploadFile, 200, expect)
}

func TestUpdatePetUploadFilePetNotFound(t *testing.T) {
	path := "./upload_test_file.log"
	body, writer := readFile(t, path, "file")
	req := httptest.NewRequest("POST", "/v2/pet/99/uploadImage", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	expect := "Pet not found"
	testTemplate(t, req, "/v2/pet/{petId}/uploadImage", UploadFile, 404, expect)
}

func TestUpdatePetUploadFilePetInvalidInput(t *testing.T) {
	path := "./upload_test_file.log"
	body, writer := readFile(t, path, "invalidInput")
	req := httptest.NewRequest("POST", "/v2/pet/99/uploadImage", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	expect := "Invalid input supplied"
	testTemplate(t, req, "/v2/pet/{petId}/uploadImage", UploadFile, 400, expect)
}

func TestFindPetsByStatusSuccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/pet/findByStatus?status=available", nil)
	expect := `[{"id":1,"name":"test1","photoUrls":null,"status":"available"}]`
	testTemplate(t, req, "/v2/pet/findByStatus", FindPetsByStatus, 200, expect)
}

func TestFindPetsByStatusInvalidStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/pet/findByStatus?status=notAvailable", nil)
	expect := "Invalid status value"
	testTemplate(t, req, "/v2/pet/findByStatus", FindPetsByStatus, 400, expect)
}

func testTemplate(t *testing.T, req *http.Request, handlePath string, handler func(http.ResponseWriter,
	*http.Request), statusCode int, expectedReturn string) {
	models.ClearData()
	models.MockData()

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc(handlePath, handler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != statusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, statusCode)
	}

	//Check the response
	if strings.TrimRight(rr.Body.String(), "\n") != expectedReturn {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedReturn)
	}

}

func readFile(t *testing.T, path string, fieldName string) (body *bytes.Buffer, writer *multipart.Writer) {
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	return body, writer
}
