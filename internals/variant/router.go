package variant

import (
	"encoding/json"
	"flukis/product/utils/resp"
	"net/http"
	"strconv"

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

	route.Post("/", r.CreateVariantHandler)
	route.Get("/{id}", r.GetVariantOneByIDHandler)
	route.Get("/", r.GetVariantsHandler)
	route.Patch("/{id}", r.UpdateDataVariantHandler)
	route.Delete("/{id}", r.DeleteVariantHandler)

	return route
}

func (r *Router) DeleteVariantHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	err = r.service.DeleteVariant(ctx, id)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "delete variant name success", http.StatusOK, nil, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) GetVariantOneByIDHandler(w http.ResponseWriter, req *http.Request) {
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	res, err := r.service.GetVariantByID(ctx, id)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "get one variant success", http.StatusOK, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) GetVariantsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	limitStr := req.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	cursor := req.URL.Query().Get("cursor")
	res, length, next, err := r.service.GetVariantsByCursor(ctx, limitInt, cursor)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}

	var metaResp struct {
		Limit    int    `json:"limit"`
		ThisPage int    `json:"total_this_page"`
		Next     string `json:"next_cursor"`
	}

	metaResp.Limit = limitInt
	metaResp.Next = next
	metaResp.ThisPage = length

	if err = resp.WriteResponse(w, "get all variants success", http.StatusOK, res, metaResp); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) UpdateDataVariantHandler(w http.ResponseWriter, req *http.Request) {
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	if err := req.ParseForm(); err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	var input struct {
		Name          string    `json:"name"`
		Description   string    `json:"desc"`
		Price         float64   `json:"price"`
		MainProductId ulid.ULID `json:"main_id"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&input)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	res, err := r.service.UpdateDataVariant(ctx, id, input.Name, input.Description, input.Price, input.MainProductId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "update variant success", http.StatusOK, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) CreateVariantHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	var input struct {
		Name          string    `json:"name"`
		Description   string    `json:"desc"`
		Price         float64   `json:"price"`
		MainProductId ulid.ULID `json:"main_id"`
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
	res, err := r.service.CreateVariant(ctx, input.Name, input.Description, input.Price, input.MainProductId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "create variant success", http.StatusCreated, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}
