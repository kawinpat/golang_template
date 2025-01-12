package handler

import (
	"golang_template/helper"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"message": "Hello World!"})
}
