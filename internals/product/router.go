package product

import (
	"encoding/json"
	"errors"
	"flukis/product/utils/helper"
	"flukis/product/utils/resp"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Router struct {
	service Service
}

func NewRouter(
	service Service,
) *Router {
	return &Router{
		service: service,
	}
}

func (r *Router) Routes() *chi.Mux {
	route := chi.NewMux()

	route.Post("/", r.CreateProductHandler)
	route.Post("/upload-image/{id}", r.UploadProductImageHandler)
	route.Get("/{id}", r.GetProductOneByIDHandler)

	return route
}

func (r *Router) GetProductOneByIDHandler(w http.ResponseWriter, req *http.Request) {
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	res, err := r.service.GetProductByID(ctx, id)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "get one product success", http.StatusOK, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) CreateProductHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"desc"`
		Price       float64 `json:"price"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&input)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	res, err := r.service.CreateProduct(ctx, input.Name, input.Description, input.Price)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "create Product success", http.StatusCreated, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) UploadProductImageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		if err := resp.WriteError(w, http.StatusMethodNotAllowed, errors.New("method not allowed")); err != nil {
			return
		}
		return
	}
	productId := chi.URLParam(req, "id")
	id, err := ulid.Parse(productId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	imageData, err := helper.UploadImageHandler(req)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	res, err := r.service.UpdateImageProduct(ctx, id, imageData)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "create Product success", http.StatusCreated, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}
