package wfstandardlib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/vivekmv23/go-web-frameworks/database"
	"github.com/vivekmv23/go-web-frameworks/lib"
)

var (
	ItemsEndpointRegex       = regexp.MustCompile(`^/items/*$`)
	ItemsWithIDEndpointRegex = regexp.MustCompile(`^/items/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)
)

type StandardLibWebServer struct {
}

func (ws *StandardLibWebServer) Start(port int) {
	mux := http.NewServeMux()

	mux.Handle("/items", &ItemsHandler{})
	mux.Handle("/items/", &ItemsHandler{})

	log.Printf("Starting Server %d, Using Standard Lib...\n", port)
	p := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(p, mux)
	if err != nil {
		log.Fatalf("Failed to start server on port %d: %s", port, err)
	}
}

type ItemsHandler struct {
}

// Satisfying the interface for handler
func (i *ItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Authorization checks
	if !isAuthorized(r) {
		errorResponse(http.StatusUnauthorized, w, r, fmt.Errorf("unauthorized, remove header 'unauthorized'"))
		return
	}

	switch {
	// Explicitly routing request based on method and url pattern :(
	case r.Method == http.MethodGet && ItemsEndpointRegex.MatchString(r.URL.Path):
		getAllItem(w, r)

	case r.Method == http.MethodGet && ItemsWithIDEndpointRegex.MatchString(r.URL.Path):
		getItem(w, r)

	case r.Method == http.MethodPost && ItemsEndpointRegex.MatchString(r.URL.Path):
		createItem(w, r)

	case r.Method == http.MethodDelete && ItemsWithIDEndpointRegex.MatchString(r.URL.Path):
		deleteItem(w, r)

	case r.Method == http.MethodPut && ItemsWithIDEndpointRegex.MatchString(r.URL.Path):
		updateItem(w, r)

	default:
		errorResponse(http.StatusMethodNotAllowed, w, r, fmt.Errorf("method %s and/or on url %s not allowed", r.Method, r.URL.Path))
	}

}

func createItem(w http.ResponseWriter, r *http.Request) {
	var itemToCreate lib.Item

	if err := json.NewDecoder(r.Body).Decode(&itemToCreate); err != nil {
		errorResponse(http.StatusBadRequest, w, r, err)
		return
	}

	if err := database.SaveItem(&itemToCreate); err != nil {
		errorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		successResponse(http.StatusCreated, w, r, itemToCreate)
	}

}

func getItem(w http.ResponseWriter, r *http.Request) {
	matches := ItemsWithIDEndpointRegex.FindStringSubmatch(r.RequestURI)
	idToGet, _ := uuid.Parse(matches[1]) // 0: full string, 1: sub string matched

	item, err := database.GetItemById(idToGet)

	if err != nil {
		// based on err, status code will be changed, e.g. 404 NOT_FOUND
		errorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		w.Header().Add("Etag", item.UpdatedOn.String())
		successResponse(http.StatusOK, w, r, item)
	}

}

func updateItem(w http.ResponseWriter, r *http.Request) {
	ifMatch := r.Header.Get("If-Match")
	if ifMatch == "" {
		errorResponse(http.StatusPreconditionRequired, w, r, fmt.Errorf("If-Match is required header for update"))
		return
	}

	matches := ItemsWithIDEndpointRegex.FindStringSubmatch(r.RequestURI)
	idToUpdate, _ := uuid.Parse(matches[1])

	var itemToUpdate lib.Item

	if err := json.NewDecoder(r.Body).Decode(&itemToUpdate); err != nil {
		errorResponse(http.StatusBadRequest, w, r, err)
		return
	}

	itemToUpdate.Id = idToUpdate

	updatedItem, err := database.UpdateItem(itemToUpdate, ifMatch)

	if err != nil {
		errorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		successResponse(http.StatusOK, w, r, updatedItem)
	}
}

func getAllItem(w http.ResponseWriter, r *http.Request) {
	items, err := database.GetAllItems()
	if err != nil {
		errorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		successResponse(http.StatusOK, w, r, items)
	}

}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	matches := ItemsWithIDEndpointRegex.FindStringSubmatch(r.RequestURI)
	idToDelete, _ := uuid.Parse(matches[1])
	if err := database.DeleteItemById(idToDelete); err != nil {
		errorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		successResponse(http.StatusNoContent, w, r, nil)
	}
}

// Call before each
func isAuthorized(r *http.Request) bool {
	log.Println("Stubbed auth check for host:", r.Host)
	return r.Header.Get("unauthorized") == ""
}

func successResponse(statusCode int, w http.ResponseWriter, r *http.Request, response any) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if response != nil {
		resJson, err := json.Marshal(response)
		if err != nil {
			errorResponse(http.StatusInternalServerError, w, r, err)
			return
		}
		w.Write(resJson)
	}
}

func errorResponse(statusCode int, w http.ResponseWriter, r *http.Request, err error) {

	e := lib.Error{
		Error: err.Error(),
		Path:  r.RequestURI,
	}

	statusCode = mapErrorToHttpStatus(err, statusCode)
	

	ejson, err := json.Marshal(e)

	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		w.Write([]byte("Failed to form the error struct"))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(ejson)
	}

}

func mapErrorToHttpStatus(err error, statusCode int) int {

	_, isNotFound := err.(*database.NotFound)
	if isNotFound {
		return http.StatusNotFound
	}

	_, isOutDated := err.(*database.Outdated)
	if isOutDated {
		return http.StatusPreconditionFailed
	}

	return statusCode
}
