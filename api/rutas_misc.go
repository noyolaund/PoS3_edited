package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func configurarRutasMisc(enrutador *mux.Router) {
	enrutador.HandleFunc("/usuario/logueado", func(w http.ResponseWriter, r *http.Request) {
		s := Sesion{
			request:        r,
			responseWriter: w,
		}
		usuario, err := s.obtenerUsuarioLogueado()
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		responderHttpExitoso(usuario, w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)
}
