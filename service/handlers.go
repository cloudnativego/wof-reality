package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func updateRealityHandler(formatter *render.Render, repo realityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		return
	}
}

func getRealityHandler(formatter *render.Render, repo realityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		return
	}
}
