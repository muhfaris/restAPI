package router

import (
	"encoding/json"
	"net/http"

	"github.com/muhfaris/restAPI/api"
)

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *api.Error)
)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var errs []string

	// Ignore error from form parsing as it's insignificant.
	r.ParseForm()

	data, err := fn(w, r)
	if err != nil {
		logger.Err.WithError(err.Err).Println("Serve error.")
		errs = append(errs, err.Error())
		w.WriteHeader(err.StatusCode)
	}

	resp := api.Response{
		Data: data,
		BaseResponse: api.BaseResponse{
			Errors: errs,
		},
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		logger.Err.WithError(err).Println("Encode response error.")
		return
	}
}
