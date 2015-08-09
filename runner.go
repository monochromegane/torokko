package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func Run() error {

	queue := make(chan *params, 1024)

	r := mux.NewRouter()
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}", storeHandler(queue)).Methods("POST")
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}", redirectHandler).Methods("GET")
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}/{filename}.tar.gz", downloadHandler).Methods("GET")
	http.Handle("/", r)

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	// start build worker
	startWorker(queue, 10)
	return nil
}

func storeHandler(queue chan *params) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := newCargo(newParams(mux.Vars(r), r.Header.Get("Authorization"))).store(queue)
		if err != nil {
			switch err.(type) {
			case aleadyExistsError:
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	cargo := newCargo(newParams(mux.Vars(r), r.Header.Get("Authorization")))
	if !cargo.isExist() {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r,
		r.URL.Path+"/"+cargo.downloadFileName(),
		http.StatusSeeOther,
	)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	cargo := newCargo(newParams(mux.Vars(r), r.Header.Get("Authorization")))
	if !cargo.isExist() || !cargo.isAuthorized() {
		http.NotFound(w, r)
		return
	}
	filepath, err := cargo.get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(r.URL.Path))
	http.ServeFile(w, r, filepath)
}
