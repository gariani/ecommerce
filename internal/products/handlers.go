package products

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gariani/ecommerce/internal/json"
)

type handler struct {
	service Service
}

func Newhandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindProductById(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		http.Error(w, "missing product id", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		log.Println("invalid ID content.")
		http.Error(w, "invalid ID content.", http.StatusNotFound)
		return
	}

	product, err := h.service.FindProductById(r.Context(), id)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	json.Write(w, http.StatusOK, product)

}
