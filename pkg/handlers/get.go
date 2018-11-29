package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
)

func (h *Handler) GetOneByIdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.FromString(ps.ByName("id"))
	if err != nil {
		log.Warn().Err(err).Msg("invalid id")
		writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	task, ok := h.cache.Get(id)
	if !ok {
		writeResponse(w, http.StatusNotFound, nil)
		return
	}

	writeResponse(w, http.StatusOK, task)
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	startPage := 0
	countTaskOnPage := 0

	page := r.URL.Query().Get("page")
	if len(page) > 0 {
		startPage, err = strconv.Atoi(page)
		if err != nil {
			writeResponse(w, http.StatusUnprocessableEntity, nil)
			return
		}
	}

	count := r.URL.Query().Get("count")
	if len(count) > 0 {
		countTaskOnPage, err = strconv.Atoi(count)
		if err != nil {
			writeResponse(w, http.StatusUnprocessableEntity, nil)
			return
		}
	}
	fmt.Println(startPage, countTaskOnPage)
	tasks := h.cache.GetAll(startPage, countTaskOnPage)
	if len(tasks) == 0 {
		writeResponse(w, http.StatusNotFound, nil)
		return
	}

	writeResponse(w, http.StatusOK, tasks)
}
