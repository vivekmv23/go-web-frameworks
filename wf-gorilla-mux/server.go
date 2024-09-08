package wfgorillamux

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/vivekmv23/go-web-frameworks/database"
	"github.com/vivekmv23/go-web-frameworks/lib"
	"github.com/vivekmv23/go-web-frameworks/web"
)

type GorillaMuxWebServer struct {
	d database.ItemDatabase
}

func NewGorillaMuxWebServer(d database.ItemDatabase) *GorillaMuxWebServer {
	return &GorillaMuxWebServer{d: d}
}

type ItemsHandler struct {
	d database.ItemDatabase
}

func (ws *GorillaMuxWebServer) Start(port int) {

	router := mux.NewRouter()

	itemsRouter := router.PathPrefix("/items").Subrouter()

	itemsRouter.Use(LogRequestMiddleware)
	itemsRouter.Use(AuthenticationMiddleware)
	itemsRouter.Use(LogResponseMiddleware)

	NewItemsHandler(ws.d, itemsRouter)

	log.Printf("Starting Server %d, Using Gorilla/Mux...\n", port)
	p := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(p, router)
	if err != nil {
		log.Fatalf("Failed to start server on port %d: %s", port, err)
	}
}

func NewItemsHandler(d database.ItemDatabase, itemsRouter *mux.Router) *ItemsHandler {
	ItemsHandler := &ItemsHandler{d: d}

	itemsRouter.HandleFunc("", ItemsHandler.GetAllItems).Methods(http.MethodGet)
	itemsRouter.HandleFunc("", ItemsHandler.CreateItem).Methods(http.MethodPost)
	itemsRouter.HandleFunc("/{id}", ItemsHandler.GetItemById).Methods(http.MethodGet)
	itemsRouter.HandleFunc("/{id}", ItemsHandler.DeleteItemById).Methods(http.MethodDelete)
	itemsRouter.HandleFunc("/{id}", ItemsHandler.UpdateItem).Methods(http.MethodPut)

	return ItemsHandler
}

func (i ItemsHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := i.d.GetAllItems()
	if err != nil {
		web.ErrorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		web.SuccessResponse(http.StatusOK, w, r, items)
	}
}

func (i ItemsHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var itemToCreate lib.Item

	if err := json.NewDecoder(r.Body).Decode(&itemToCreate); err != nil {
		web.ErrorResponse(http.StatusBadRequest, w, r, err)
		return
	}

	if err := i.d.SaveItem(&itemToCreate); err != nil {
		web.ErrorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		web.SuccessResponse(http.StatusCreated, w, r, itemToCreate)
	}
}

func (i ItemsHandler) GetItemById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	idToGet, err := uuid.Parse(id)
	if err != nil {
		web.ErrorResponse(http.StatusBadRequest, w, r, err)
		return
	}

	item, err := i.d.GetItemById(idToGet)

	if err != nil {
		// based on err, status code will be changed, e.g. 404 NOT_FOUND
		web.ErrorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		w.Header().Add("Etag", item.UpdatedOn.String())
		web.SuccessResponse(http.StatusOK, w, r, item)
	}
}

func (i ItemsHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ifMatch := r.Header.Get("If-Match")
	if ifMatch == "" {
		web.ErrorResponse(http.StatusPreconditionRequired, w, r, fmt.Errorf("If-Match is required header for update"))
		return
	}

	id := mux.Vars(r)["id"]
	idToUpdate, err := uuid.Parse(id)

	if err != nil {
		web.ErrorResponse(http.StatusBadRequest, w, r, err)
		return
	}

	var itemToUpdate lib.Item

	if err := json.NewDecoder(r.Body).Decode(&itemToUpdate); err != nil {
		web.ErrorResponse(http.StatusBadRequest, w, r, err)
		return
	}

	itemToUpdate.Id = idToUpdate

	updatedItem, err := i.d.UpdateItem(itemToUpdate, ifMatch)

	if err != nil {
		web.ErrorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		web.SuccessResponse(http.StatusOK, w, r, updatedItem)
	}
}

func (i ItemsHandler) DeleteItemById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idToDelete, err := uuid.Parse(id)

	if err != nil {
		web.ErrorResponse(http.StatusBadRequest, w, r, err)
		return
	}

	if err := i.d.DeleteItemById(idToDelete); err != nil {
		web.ErrorResponse(http.StatusInternalServerError, w, r, err)
	} else {
		web.SuccessResponse(http.StatusNoContent, w, r, nil)
	}
}

func AuthenticationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		if !web.IsAuthorized(r) {
			web.ErrorResponse(http.StatusUnauthorized, w, r, fmt.Errorf("unauthorized, remove header 'unauthorized'"))
			return
		}
		h.ServeHTTP(w, r)
	})
}

func LogRequestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: UserAgent: %s, Host: %s, Method: %s, URI: %s, Cookies: %q", r.UserAgent(), r.Host, r.Method, r.RequestURI, r.Cookies())
		h.ServeHTTP(w, r)
	})
}

func LogResponseMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(w http.ResponseWriter) {
			log.Printf("RESPONSE: Content-Type: %s", w.Header().Get("Content-Type"))
		}(w)
		h.ServeHTTP(w, r)
	})
}
