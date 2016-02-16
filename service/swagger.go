package service

import (
	"io/ioutil"
	"net/http"
)

func (s *TranscodingService) swaggerManifest(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(s.config.SwaggerManifest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// In order to support pure-JavaScript clients (like Swagger-UI), the
	// server must set CORS headers.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
