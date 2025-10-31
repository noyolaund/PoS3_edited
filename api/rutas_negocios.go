package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func configurarRutasDeNegocios(enrutador *mux.Router) {

	enrutador.HandleFunc("/negocio", func(w http.ResponseWriter, r *http.Request) {
		var negocio Negocio
		err := json.NewDecoder(r.Body).Decode(&negocio)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		nc := NegociosController{}
		responderHttpExitoso(nc.nuevo(&negocio), w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodPost)

	enrutador.HandleFunc("/negocio/verificar/{token}", func(w http.ResponseWriter, r *http.Request) {
		variablesDePeticion := mux.Vars(r)
		if variablesDePeticion == nil {
			responderHttpConError(errors.New("no hay token al verificar negocio"), w, r)
			return
		}
		token := variablesDePeticion["token"]
		nc := NegociosController{}
		responderHttpExitoso(nc.verificarPorToken(token), w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)

	enrutador.HandleFunc("/negocio/eliminar/{token}", func(w http.ResponseWriter, r *http.Request) {
		variablesDePeticion := mux.Vars(r)
		if variablesDePeticion == nil {
			responderHttpConError(errors.New("no hay token al eliminar negocio"), w, r)
			return
		}
		token := variablesDePeticion["token"]
		nc := NegociosController{}
		responderHttpExitoso(nc.eliminarNegocioPorToken(token), w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)
}
