package cargo

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}/", buildHandler).Methods("POST")
	r.HandleFunc("/{remote}/{user}/{repo}/{goos}/{goarch}/{version}/", downloadHandler).Methods("GET")
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

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	cargo := newCargo(mux.Vars(r))
	if !cargo.isExist() {
		http.NotFound(w, r)
		return
	}
	path, err := cargo.get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, path)
}
