package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.FromString(ps.ByName("id"))
	if err != nil {
		log.Warn().Err(err).Msg("invalid id")
		writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	ok := h.cache.Delete(id)
	if !ok {
		writeResponse(w, http.StatusNotFound, nil)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}
