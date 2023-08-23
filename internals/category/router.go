package category

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

	route.Post("/", r.CreateCategoryHandler)
	route.Patch("/{id}", r.UpdateCategoryHandler)
	route.Delete("/{id}", r.DeleteCategoryHandler)
	route.Get("/{id}", r.GetCategoryOneByIDHandler)
	route.Get("/", r.GetCategorysHandler)

	return route
}

func (r *Router) CreateCategoryHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	var input struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
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
	res, err := r.service.CreateCategory(ctx, input.Name, input.Desc)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "create Category success", http.StatusCreated, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) UpdateCategoryHandler(w http.ResponseWriter, req *http.Request) {
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
		Name string `json:"name"`
		Desc string `json:"desc"`
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
	res, err := r.service.UpdateCategory(ctx, id, input.Name, input.Desc)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "update Category name success", http.StatusOK, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) DeleteCategoryHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	err = r.service.DeleteCategory(ctx, id)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "delete Category name success", http.StatusOK, nil, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) GetCategoryOneByIDHandler(w http.ResponseWriter, req *http.Request) {
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	res, err := r.service.GetCategoryById(ctx, id)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "get one Category success", http.StatusOK, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) GetCategorysHandler(w http.ResponseWriter, req *http.Request) {
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
	res, length, next, err := r.service.GetCategoryByCursor(ctx, limitInt, cursor)
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

	if err = resp.WriteResponse(w, "get all Category success", http.StatusOK, res, metaResp); err != nil {
		log.Error().Err(err)
		return
	}
}
