package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vivekmv23/go-web-frameworks/database"
	"github.com/vivekmv23/go-web-frameworks/lib"
)

type WebServer interface {
	Start(port int)
}

func SuccessResponse(statusCode int, w http.ResponseWriter, r *http.Request, response any) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if response != nil {
		resJson, err := json.Marshal(response)
		if err != nil {
			ErrorResponse(http.StatusInternalServerError, w, r, err)
			return
		}
		w.Write(resJson)
	}
}

func ErrorResponse(statusCode int, w http.ResponseWriter, r *http.Request, err error) {

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

func IsAuthorized(r *http.Request) bool {
	isAuth := r.Header.Get("unauthorized") == ""
	log.Printf("Stubbed auth check for host: %s authorized: %v", r.Host, isAuth)
	return isAuth
}
