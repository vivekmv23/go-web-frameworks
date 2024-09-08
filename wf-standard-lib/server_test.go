package wfstandardlib

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vivekmv23/go-web-frameworks/database"
)

var (
	error_generic   error              = fmt.Errorf("generic error")
	error_not_found *database.NotFound = &database.NotFound{Id: "some-id"}
	error_outdated  *database.Outdated = &database.Outdated{}
	error_conflict  *database.Conflict = &database.Conflict{}
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("../testdata/" + name)
	if err != nil {
		t.Errorf("Could not read %v", name)
	}

	return content
}

func TestServer_GetAll(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)

	r := httptest.NewRequest(http.MethodGet, "/items", nil)
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)
	assert.NotEmpty(t, res.Body)

	d = database.NewMockedDatabase(error_generic)
	ih = NewItemsHandler(d)
	w = httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 500, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestServer_GetById(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)

	r := httptest.NewRequest(http.MethodGet, "/items/fe9dd883-7b95-4d7a-80d9-0c80423a8e16", nil)
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)
	assert.NotEmpty(t, res.Body)

	d = database.NewMockedDatabase(error_generic)
	ih = NewItemsHandler(d)
	w = httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 500, res.StatusCode)
	assert.NotEmpty(t, res.Body)

	d = database.NewMockedDatabase(error_not_found)
	ih = NewItemsHandler(d)
	w = httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 404, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestServer_GetById_Unauthorized(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)
	r := httptest.NewRequest(http.MethodGet, "/items/fe9dd883-7b95-4d7a-80d9-0c80423a8e16", nil)
	w := httptest.NewRecorder()
	r.Header.Add("unauthorized", "true")

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 401, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestServer_SaveItem(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)
	itemToSave := readTestData(t, "item-payload.json")
	itemToSaveReader := bytes.NewReader(itemToSave)

	r := httptest.NewRequest(http.MethodPost, "/items", itemToSaveReader)
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 201, res.StatusCode)
	assert.NotEmpty(t, res.Body)

	d = database.NewMockedDatabase(error_conflict)
	ih = NewItemsHandler(d)
	itemToSave = readTestData(t, "item-payload.json")
	itemToSaveReader = bytes.NewReader(itemToSave)

	r = httptest.NewRequest(http.MethodPost, "/items", itemToSaveReader)
	w = httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res = w.Result()
	defer res.Body.Close()

	assert.Equal(t, 500, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestServer_DeleteItem(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)
	r := httptest.NewRequest(http.MethodDelete, "/items/fe9dd883-7b95-4d7a-80d9-0c80423a8e16", nil)
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 204, res.StatusCode)
	assert.NotEmpty(t, res.Body)

	d = database.NewMockedDatabase(error_not_found)
	ih = NewItemsHandler(d)
	w = httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 404, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestServer_UpdateItem(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)
	itemToUpdate := readTestData(t, "item-payload.json")
	itemToUpdateReader := bytes.NewReader(itemToUpdate)

	r := httptest.NewRequest(http.MethodPut, "/items/fe9dd883-7b95-4d7a-80d9-0c80423a8e16", itemToUpdateReader)
	r.Header.Add("If-Match", "some-e-tag")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)
	assert.NotEmpty(t, res.Body)

	d = database.NewMockedDatabase(error_outdated)
	ih = NewItemsHandler(d)
	itemToUpdate = readTestData(t, "item-payload.json")
	itemToUpdateReader = bytes.NewReader(itemToUpdate)

	r = httptest.NewRequest(http.MethodPut, "/items/fe9dd883-7b95-4d7a-80d9-0c80423a8e16", itemToUpdateReader)
	r.Header.Add("If-Match", "some-e-tag")
	w = httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res = w.Result()
	defer res.Body.Close()

	assert.Equal(t, 412, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestServer_UpdateItem_Missing_IfMatch_Header(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)
	itemToUpdate := readTestData(t, "item-payload.json")
	itemToUpdateReader := bytes.NewReader(itemToUpdate)

	r := httptest.NewRequest(http.MethodPut, "/items/fe9dd883-7b95-4d7a-80d9-0c80423a8e16", itemToUpdateReader)
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 428, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestServer_MethodNotAllowed(t *testing.T) {
	d := database.NewMockedDatabase(nil)
	ih := NewItemsHandler(d)

	r := httptest.NewRequest(http.MethodOptions, "/items/fe9dd883-7b95-4d7a-80d9-0c80423a8e16", nil)
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 405, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}
