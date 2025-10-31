package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func configurarRutasPermisos(enrutador *mux.Router) {
	enrutador.HandleFunc("/permisos", func(w http.ResponseWriter, r *http.Request) {
		pc := PermisosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.todos(), w, r)
	}).Name("ModificarYVerPermisos").Methods(http.MethodGet)

	enrutador.HandleFunc("/permisos/de/{idUsuario}", func(w http.ResponseWriter, r *http.Request) {
		variablesDePeticion := mux.Vars(r)
		idUsuario, err := strconv.Atoi(variablesDePeticion["idUsuario"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := PermisosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.deUnUsuario(idUsuario), w, r)
	}).Name("ModificarYVerPermisos").Methods(http.MethodGet)

	enrutador.HandleFunc("/permisos/para/{idUsuario}", func(w http.ResponseWriter, r *http.Request) {
		var permisos []int
		err := json.NewDecoder(r.Body).Decode(&permisos)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		variablesDePeticion := mux.Vars(r)
		idUsuario, err := strconv.Atoi(variablesDePeticion["idUsuario"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := PermisosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.asignarA(idUsuario, permisos), w, r)
	}).Name("ModificarYVerPermisos").Methods(http.MethodPut)
}
