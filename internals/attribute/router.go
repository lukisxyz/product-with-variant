package attribute

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

	route.Post("/", r.CreateAttributeHandler)
	route.Patch("/{id}", r.UpdateAttributeHandler)
	route.Delete("/{id}", r.DeleteAttributeHandler)
	route.Get("/{id}", r.GetAttributeOneByIDHandler)
	route.Get("/", r.GetAttributesHandler)

	return route
}

func (r *Router) CreateAttributeHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	var input struct {
		Name string `json:"name"`
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
	res, err := r.service.CreateAttr(ctx, input.Name)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "create attribute success", http.StatusCreated, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) UpdateAttributeHandler(w http.ResponseWriter, req *http.Request) {
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
	res, err := r.service.UpdateNameAttr(ctx, id, input.Name)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "update attribute name success", http.StatusOK, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) DeleteAttributeHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	err = r.service.DeleteAttr(ctx, id)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "delete attribute name success", http.StatusOK, nil, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) GetAttributeOneByIDHandler(w http.ResponseWriter, req *http.Request) {
	categoryId := chi.URLParam(req, "id")
	id, err := ulid.Parse(categoryId)
	if err != nil {
		if err = resp.WriteError(w, http.StatusBadRequest, err); err != nil {
			return
		}
		return
	}
	ctx := req.Context()
	res, err := r.service.GetAttrById(ctx, id)
	if err != nil {
		if err = resp.WriteError(w, http.StatusInternalServerError, err); err != nil {
			log.Error().Err(err)
			return
		}
		return
	}
	if err = resp.WriteResponse(w, "get one attribute success", http.StatusOK, res, nil); err != nil {
		log.Error().Err(err)
		return
	}
}

func (r *Router) GetAttributesHandler(w http.ResponseWriter, req *http.Request) {
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
	res, length, next, err := r.service.GetAttrByCursor(ctx, limitInt, cursor)
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

	if err = resp.WriteResponse(w, "get all attribute success", http.StatusOK, res, metaResp); err != nil {
		log.Error().Err(err)
		return
	}
}
