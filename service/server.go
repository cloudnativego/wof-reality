package service

import (
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer(appEnv *cfenv.App) *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	repo := initRepository(appEnv)

	initRoutes(mx, formatter, repo)

	n.UseHandler(mx)
	return n
}

func initRepository(appEnv *cfenv.App) (repo realityRepository) {
	return
}

func initRoutes(mx *mux.Router, formatter *render.Render, repo realityRepository) {
	mx.HandleFunc("/reality/{gameId}", updateRealityHandler(formatter, repo)).Methods("PUT")
	mx.HandleFunc("/reality/{gameId}", getRealityHandler(formatter, repo)).Methods("GET")
}
