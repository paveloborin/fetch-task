package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/paveloborin/fetch-task/pkg/model"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
)

var httpMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodConnect: true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request body")
	}
	defer r.Body.Close()

	req := struct {
		Method string `json:"method"`
		URI    string `json:"uri"`
	}{}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Warn().Err(err).Msg("unmarshal body err")
		writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	method := strings.ToUpper(req.Method)
	if _, ok := httpMethods[method]; !ok {
		log.Warn().Err(fmt.Errorf("not allowed method %s", req.Method)).Msg("invalid method param")
		writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	u, err := url.ParseRequestURI(req.URI)
	if err != nil {
		log.Warn().Err(fmt.Errorf("invalid uri %s", req.URI)).Msg("invalid uri param")
		writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	task := &model.Task{ID: uuid.NewV4()}
	if err := getURI(u, method, task); err != nil {
		log.Warn().Err(err).Msg("get url error")
		writeResponse(w, http.StatusUnprocessableEntity, nil)
		return
	}

	go func() {
		h.cache.Add(task)
	}()

	writeResponse(w, http.StatusOK, task)
}

func getURI(uri *url.URL, method string, task *model.Task) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	task.Body = string(html)

	headers := map[string]string{}
	for name, value := range resp.Header {
		headers[name] = strings.Join(value, " ")

	}
	task.Headers = headers

	task.Status = resp.StatusCode
	task.ContentLength = len(html)
	return nil
}
