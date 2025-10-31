package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func configurarRutasDeUsuarios(enrutador *mux.Router) {
	enrutador.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		usuariosController := UsuariosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(usuariosController.todos(), w, r)
	}).Name("VerUsuarios").Methods(http.MethodGet)

	enrutador.HandleFunc("/usuario/caja/{idUsuario}", func(w http.ResponseWriter, r *http.Request) {
		variablesDePeticion := mux.Vars(r)
		idUsuario, err := strconv.Atoi(variablesDePeticion["idUsuario"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		usuariosController := UsuariosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(usuariosController.uno(idUsuario), w, r)
	}).Name("VerReporteCaja").Methods(http.MethodGet)

	enrutador.HandleFunc("/usuario", func(w http.ResponseWriter, r *http.Request) {
		var usuario Usuario
		err := json.NewDecoder(r.Body).Decode(&usuario)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		usuariosController := UsuariosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(usuariosController.nuevo(&usuario), w, r)
	}).Name("RegistrarUsuario").Methods(http.MethodPost)
}
