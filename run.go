package cargo

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}", buildHandler).Methods("POST")
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}", redirectHandler).Methods("GET")
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}/{filename}.tar.gz", downloadHandler).Methods("GET")
	http.Handle("/", r)

	return http.ListenAndServe(":8080", nil)
}

func buildHandler(w http.ResponseWriter, r *http.Request) {
	err := newCargo(mux.Vars(r)).build()
	if err != nil {
		switch err.(type) {
		case aleadyExistsError:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	cargo := newCargo(mux.Vars(r))
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
	cargo := newCargo(mux.Vars(r))
	if !cargo.isExist() {
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
