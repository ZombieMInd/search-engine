package server

import (
	"encoding/json"
	"net/http"

	"github.com/ZombieMInd/go-logger/internal/logger"
)

func (s *server) handleLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &logger.LogRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := req.Validate()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		req.IP = r.RemoteAddr

		err = s.services.Log.Save(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusCreated, nil)
	}
}
