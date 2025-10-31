package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func configurarRutasUtiles(enrutador *mux.Router) {
	enrutador.HandleFunc("/fechaYHora", func(w http.ResponseWriter, r *http.Request) {
		responderHttpExitoso(obtenerFechaActualFormateada(), w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)
}
